<template>
  <div
    class="shell"
    :class="{ ['shell--hidden']: !showShell }"
    :style="{ height: `${this.shellHeight}em`, direction: 'ltr' }"
  >
    <div
      @pointerdown="startDrag()"
      @pointerup="stopDrag()"
      class="shell__divider"
      :style="this.shellDrag ? { background: `${checkTheme()}` } : ''"
    ></div>
    <div ref="terminal"></div>
    <div
      @pointerup="stopDrag()"
      class="shell__overlay"
      v-show="this.shellDrag"
    ></div>
  </div>
</template>

<script>
import { mapState, mapActions } from "pinia";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

import shell from "@/api/shell";
import { throttle } from "lodash-es";
import { theme } from "@/utils/constants";

import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";

export default {
  name: "shell",
  computed: {
    ...mapState(useLayoutStore, ["showShell"]),
    ...mapState(useFileStore, ["isFiles"]),
    path: function () {
      if (this.isFiles) {
        return this.$route.path;
      }
      return "";
    },
  },
  data: () => ({
    shellHeight: 25,
    shellDrag: false,
    fontsize: parseFloat(getComputedStyle(document.documentElement).fontSize),
    term: null,
    fitAddon: null,
    conn: null,
    resize: null,
  }),
  mounted() {
    this.term = new Terminal();
    this.fitAddon = new FitAddon();
    this.term.loadAddon(this.fitAddon);
    this.term.open(this.$refs.terminal);
    this.fitAddon.fit();
    this.term.onData((data) => {
      this.conn.send(data);
    });
    [this.conn, this.resize] = shell(this.path, (ev) => {
      this.term.write(ev.data);
    });
    window.addEventListener("resize", this.resizeTerm);
    this.resizeTerm();
  },
  beforeUnmount() {
    window.removeEventListener("resize", this.resizeTerm);
    this.conn.close();
  },
  methods: {
    ...mapActions(useLayoutStore, ["toggleShell"]),
    checkTheme() {
      if (theme == "dark") {
        return "rgba(255, 255, 255, 0.4)";
      }
      return "rgba(127, 127, 127, 0.4)";
    },
    startDrag() {
      document.addEventListener("pointermove", this.handleDrag);
      this.shellDrag = true;
    },
    stopDrag() {
      document.removeEventListener("pointermove", this.handleDrag);
      this.shellDrag = false;
    },
    handleDrag: throttle(function (event) {
      const top = window.innerHeight / this.fontsize - 4;
      const userPos = (window.innerHeight - event.clientY) / this.fontsize;
      const bottom =
        2.25 +
        document.querySelector(".shell__divider").offsetHeight / this.fontsize;
      if (userPos <= top && userPos >= bottom) {
        this.shellHeight = userPos.toFixed(2);
      }
    }, 32),
    resizeTerm: throttle(function () {
      const top = window.innerHeight / this.fontsize - 4;
      const bottom =
        2.25 +
        document.querySelector(".shell__divider").offsetHeight / this.fontsize;
      if (this.shellHeight > top) {
        this.shellHeight = top;
      } else if (this.shellHeight < bottom) {
        this.shellHeight = bottom;
      }
      this.fitAddon.fit();
      this.resize(this.term.cols, this.term.rows);
    }, 32),
  },
};
</script>

<style>
@import "xterm/css/xterm.css";
</style>
