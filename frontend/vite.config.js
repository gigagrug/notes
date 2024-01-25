// vite.config.js
import { resolve } from "path"
import { defineConfig } from "vite"

const root = resolve(__dirname, "src")
export default defineConfig({
  root,
  build: {
    outDir: resolve(__dirname, "dist"),
    emptyOutDir: true,
    rollupOptions: {
      input: {
        index: resolve(root, "index.html"),
        blog: resolve(root, "getBlog.html")
      }
    }
  }
})
