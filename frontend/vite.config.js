import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import tailwindcss from '@tailwindcss/vite';
import path from 'path';

export default defineConfig({
    plugins: [vue(), tailwindcss()],
    resolve: {
        alias: {
            '@': path.resolve(__dirname, 'src'),
        },
    },
    base: '/',
    build: {
        outDir: 'dist',
        assetsDir: 'assets',
        rollupOptions: {
            output: {
                entryFileNames: `assets/[name]-[hash].js`,
                chunkFileNames: `assets/[name]-[hash].js`,
                assetFileNames: `assets/[name]-[hash].[ext]`,
                manualChunks: {
                    // Vendor libraries
                    'vue-vendor': ['vue', 'vue-router', 'pinia'],
                    'ui-vendor': ['@fortawesome/fontawesome-svg-core', '@fortawesome/vue-fontawesome'],
                    'icons': [
                        '@fortawesome/free-solid-svg-icons',
                        '@fortawesome/free-regular-svg-icons',
                        '@fortawesome/free-brands-svg-icons'
                    ],
                    'utils': ['axios', 'lodash'],

                    // Feature-based chunks (matching route chunk names)
                    'user-management': [
                        'src/views/UserManagementWrapper.vue',
                        'src/views/UserList.vue',
                        'src/views/UserEdit.vue',
                        'src/views/StaffEdit.vue',
                        'src/views/NewStaff.vue',
                        'src/views/NewUser.vue'
                    ],
                    'appointments': [
                        'src/views/Agenda.vue',
                        'src/views/AppointmentCalendar.vue',
                        'src/views/AppointmentList.vue',
                        'src/views/NewAppointment.vue',
                        'src/views/EditAppointment.vue'
                    ],
                    'patients': [
                        'src/views/PatientList.vue',
                        'src/views/PatientForm.vue',
                        'src/views/PatientDentalChart.vue',
                        'src/views/PatientDiagnosisForm.vue',
                        'src/views/PatientTreatmentPlanForm.vue'
                    ],
                    'clinical': [
                        'src/views/ProcedureManagement.vue',
                        'src/views/ProcedureForm.vue',
                        'src/views/DiagnosisForm.vue',
                        'src/views/TreatmentPlanForm.vue',
                        'src/views/ConsentForm.vue',
                        'src/views/PrescriptionForm.vue'
                    ],
                    'admin': [
                        'src/views/ClinicManagement.vue',
                        'src/views/ClinicForm.vue'
                    ],
                    'peer-review': [
                        'src/views/PeerReviewList.vue',
                        'src/views/PeerReviewCase.vue',
                        'src/views/PeerReviewCreate.vue'
                    ],
                    'inventory': [
                        'src/views/InventoryList.vue',
                        'src/views/InventoryForm.vue',
                        'src/views/InventoryView.vue'
                    ]
                }
            },
        },
        // Additional optimization settings
        target: 'esnext',
        minify: 'terser',
        terserOptions: {
            compress: {
                drop_console: true,
                drop_debugger: true
            }
        },
        cssCodeSplit: true,
        chunkSizeWarningLimit: 500
    },
    server: {
        port: 5173,
        proxy: {
            '/api': {
                target: 'http://localhost:9483',
                changeOrigin: true,
            },
            '/upload': {
                target: 'http://localhost:9483',
                changeOrigin: true,
            },
        },
        allowedHosts: ['localhost', 'app.majidarif.com'],
    },
});
