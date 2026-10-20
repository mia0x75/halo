import tailwindcss from "@tailwindcss/vite"
import react from "@vitejs/plugin-react"
import { fileURLToPath, URL } from "node:url"
import { defineConfig } from "vite"

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: { "@": fileURLToPath(new URL("./src", import.meta.url)) },
  },
  server: {
    // Vite 默认只监听 localhost 的 IPv6 回环，浏览器用 127.0.0.1 访问会连不上；
    // 监听所有网卡后 IPv4、IPv6 与容器外访问都可用。
    host: true,
    proxy: {
      // The Go backend serves GraphQL under /api; it listens on :443 with the
      // self-signed certificate from certs/, hence secure:false for dev only.
      "/api": {
        target: "https://127.0.0.1:443",
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
