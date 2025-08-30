import { baseURL } from "@/utils/constants";
import { removePrefix } from "./utils";

const ssl = window.location.protocol === "https:";
const protocol = ssl ? "wss:" : "ws:";

export default function shell(
  url: string,
  onmessage: WebSocket["onmessage"]
): [WebSocket, Function] {
  url = removePrefix(url);
  url = `${protocol}//${window.location.host}${baseURL}/api/shell${url}`;

  const conn = new window.WebSocket(url);
  conn.onmessage = onmessage;

  const resize = (cols, rows) => {
    conn.send(JSON.stringify({ cols, rows }));
  };

  return [conn, resize];
}
