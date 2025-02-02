import { defineConfig, splitVendorChunkPlugin } from 'vite'
import preact from '@preact/preset-vite'
import { resolve } from 'path'
import { visualizer } from "rollup-plugin-visualizer";

// https://vitejs.dev/config/
export default defineConfig(({ }) => {

    return {
        server: {
            hmr: {
                port: 5000
            },
            watch: {
                ignored: ["**/.env"]
            }
        },
        build: {
            outDir: 'static/build/',

            rollupOptions: {
                output: {
                    manualChunks(id: string) {
                        // creating a chunk to @open-ish deps. Reducing the vendor chunk size
                        if (id.includes('preact')) {
                            return '@base';
                        }
                        // creating a chunk to react routes deps. Reducing the vendor chunk size
                        if (
                            id.includes('react-router-dom') ||
                            id.includes('@remix-run') ||
                            id.includes('react-router')
                        ) {
                            return '@router';
                        }
                    },
                },
            },
        },
        resolve: {
            alias: {
                "i18n": resolve(__dirname, 'i18n'),
                "~": resolve(__dirname, 'static'),
                "@": resolve(__dirname, 'pkg'),
                "ui": resolve(__dirname, 'ui'),
                "ws": resolve(__dirname, 'internal/websocket'),
            }
        },
        plugins: [
            preact(),
            splitVendorChunkPlugin(),
            visualizer()
        ],
    }
})
