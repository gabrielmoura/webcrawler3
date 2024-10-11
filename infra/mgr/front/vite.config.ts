import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react-swc'
import {TanStackRouterVite} from "@tanstack/router-vite-plugin";
import tsconfigPaths from "vite-tsconfig-paths";
// import { optimizeLodashImports } from "@optimize-lodash/rollup-plugin";


// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    TanStackRouterVite({
      routeFileIgnorePrefix: '-', // Ignora arquivos que começam com "-"
    }),
    tsconfigPaths(),
    // optimizeLodashImports()
  ],
  server: {
    watch: {
      usePolling: true,
    },
  },
})