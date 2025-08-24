import { fetchURL, removePrefix, createURL } from "./utils";
import { baseURL } from "@/utils/constants";

export async function fetch(url: string, password: string = "") {
  url = removePrefix(url);

  const res = await fetchURL(
    `/api/public/share${url}`,
    {
      headers: { "X-SHARE-PASSWORD": encodeURIComponent(password) },
    },
    false
  );

  const data = (await res.json()) as Resource;
  data.url = `/share${url}`;

  if (data.isDir) {
    if (!data.url.endsWith("/")) data.url += "/";
    data.items = data.items.map((item: any, index: any) => {
      item.index = index;
      item.url = `${data.url}${encodeURIComponent(item.name)}`;

      if (item.isDir) {
        item.url += "/";
      }

      return item;
    });
  }

  return data;
}

import { useDownloadStore } from "@/stores/download";

export async function download(
  format: DownloadFormat,
  hash: string,
  token: string,
  ...files: string[]
) {
  const downloadStore = useDownloadStore();
  let url = `${baseURL}/api/public/dl/${hash}`;
  let filename: string;

  if (files.length === 1) {
    url += encodeURIComponent(files[0]) + "?";
    filename = files[0].split("/").pop() || "download";
  } else {
    let arg = "";
    for (const file of files) {
      arg += encodeURIComponent(file) + ",";
    }
    arg = arg.substring(0, arg.length - 1);
    arg = encodeURIComponent(arg);
    url += `/?files=${arg}&`;
    filename = "archive";
  }

  if (format) {
    url += `algo=${format}&`;
    filename += `.${format}`;
  }

  if (token) {
    url += `token=${token}&`;
  }

  const response = await window.fetch(url);
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const contentLength = response.headers.get("Content-Length");
  const total = parseInt(contentLength || "0", 10);
  let loaded = 0;

  const downloadId = downloadStore.add(filename);

  const reader = response.body!.getReader();
  const stream = new ReadableStream({
    start(controller) {
      function push() {
        reader
          .read()
          .then(({ done, value }) => {
            if (done) {
              controller.close();
              downloadStore.remove(downloadId);
              return;
            }
            loaded += value.length;
            if (total > 0) {
              const progress = Math.round((loaded / total) * 100);
              downloadStore.update(downloadId, progress);
            }
            controller.enqueue(value);
            push();
          })
          .catch((err: any) => {
            console.error(err);
            controller.error(err);
            downloadStore.remove(downloadId);
          });
      }
      push();
    },
  });

  const blob = await new Response(stream).blob();
  const blobUrl = window.URL.createObjectURL(blob);

  const a = document.createElement("a");
  a.style.display = "none";
  a.href = blobUrl;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  window.URL.revokeObjectURL(blobUrl);
  document.body.removeChild(a);
}

export function getDownloadURL(res: Resource, inline = false) {
  const params = {
    ...(inline && { inline: "true" }),
    ...(res.token && { token: res.token }),
  };

  return createURL("api/public/dl/" + res.hash + res.path, params);
}
