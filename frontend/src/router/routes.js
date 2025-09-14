// Critical components loaded immediately
import Login from '../views/Login.vue';
import Dashboard from '../views/Dashboard.vue';
import AppLayout from '../layouts/AppLayout.vue';

// Lazy-loaded components with chunk names for better organization
const UserManagementWrapper = () => import('../views/UserManagementWrapper.vue');
const UserList = () => import('../views/UserList.vue');
const UserEdit = () => import('../views/UserEdit.vue');
const StaffEdit = () => import('../views/StaffEdit.vue');
const NewStaff = () => import('../views/NewStaff.vue');
const NewUser = () => import('../views/NewUser.vue');

const Agenda = () => import('../views/Agenda.vue');
const AppointmentCalendar = () => import('../views/AppointmentCalendar.vue');
const AppointmentList = () => import('../views/AppointmentList.vue');
const NewAppointment = () => import('../views/NewAppointment.vue');
const EditAppointment = () => import('../views/EditAppointment.vue');

const PatientList = () => import('../views/PatientList.vue');
const PatientForm = () => import('../views/PatientForm.vue');
const PatientDentalChart = () => import('../views/PatientDentalChart.vue');

const ProcedureManagement = () => import('../views/ProcedureManagement.vue');
const ProcedureForm = () => import('../views/ProcedureForm.vue');
const DiagnosisForm = () => import('../views/DiagnosisForm.vue');
const TreatmentPlanForm = () => import('../views/TreatmentPlanForm.vue');
const ConsentForm = () => import('../views/ConsentForm.vue');
const PrescriptionForm = () => import('../views/PrescriptionForm.vue');

const ClinicManagement = () => import('../views/ClinicManagement.vue');
const ClinicForm = () => import('../views/ClinicForm.vue');

const PeerReviewList = () => import('../views/PeerReviewList.vue');
const PeerReviewCase = () => import('../views/PeerReviewCase.vue');
const PeerReviewCreate = () => import('../views/PeerReviewCreate.vue');

const InventoryList = () => import('../views/InventoryList.vue');
const InventoryForm = () => import('../views/InventoryForm.vue');
const InventoryView = () => import('../views/InventoryView.vue');

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: { hideLayout: true },
    },
    {
        path: '/',
        component: AppLayout,
        meta: { requiresAuth: true },
        children: [
            // Home/Dashboard - Level 0 (no back button)
            {
                path: '',
                name: 'Dashboard',
                component: Dashboard,
                meta: { level: 0, title: 'Dashboard' },
            },

            // Level 1 pages (show back button to home)
            {
                path: '/agenda',
                name: 'Agenda',
                component: Agenda,
                meta: { level: 1, title: 'Daily Agenda', parent: 'Dashboard', requiresNonSuperAdmin: true },
            },
            {
                path: '/patients',
                name: 'PatientList',
                component: PatientList,
                meta: { level: 1, title: 'Patients', parent: 'Dashboard' },
            },
            {
                 path: '/appointments',
                 name: 'AppointmentCalendar',
                 component: AppointmentCalendar,
                 meta: { level: 1, title: 'Appointments', parent: 'Dashboard' },
             },
             {
                 path: '/appointments/list',
                 name: 'AppointmentList',
                 component: AppointmentList,
                 meta: { level: 1, title: 'Appointment List', parent: 'Dashboard' },
             },
            {
                path: '/procedures',
                name: 'ProcedureManagement',
                component: ProcedureManagement,
                meta: {
                    level: 1,
                    title: 'Procedures & Diagnoses',
                    parent: 'Dashboard',
                },
            },
            {
                path: '/users',
                name: 'UserManagement',
                component: UserManagementWrapper,
                meta: { level: 1, title: 'User Management', parent: 'Dashboard' },
            },
            {
                path: '/staff',
                name: 'StaffManagement',
                component: UserManagementWrapper,
                meta: { level: 1, title: 'Staff Management', parent: 'Dashboard' },
            },
            {
                path: '/staff/new',
                name: 'StaffCreate',
                component: NewStaff,
                meta: { level: 2, title: 'Add Staff Member', parent: 'StaffManagement' },
            },
            {
                path: '/clinics',
                name: 'ClinicManagement',
                component: ClinicManagement,
                meta: { level: 1, title: 'Clinics', parent: 'Dashboard' },
            },
            {
                path: '/peer-review',
                name: 'PeerReviewList',
                component: PeerReviewList,
                meta: { level: 1, title: 'Peer Review', parent: 'Dashboard' },
            },
            {
                path: '/peer-review/create',
                name: 'PeerReviewCreate',
                component: PeerReviewCreate,
                meta: {
                    level: 2,
                    title: 'Create Peer Review Case',
                    parent: 'PeerReviewList',
                },
            },

            // Level 2 pages (show back button to level 1)
            {
                path: '/clinics/new',
                name: 'ClinicCreate',
                component: ClinicForm,
                meta: {
                    level: 2,
                    title: 'New Clinic',
                    parent: 'ClinicManagement',
                },
            },
            {
                path: '/clinics/:id/edit',
                name: 'ClinicEdit',
                component: ClinicForm,
                meta: {
                    level: 2,
                    title: 'Edit Clinic',
                    parent: 'ClinicManagement',
                },
            },
            {
                path: '/patients/new',
                name: 'PatientCreate',
                component: PatientForm,
                meta: { level: 2, title: 'New Patient', parent: 'PatientList' },
            },
            {
                path: '/patients/:id/edit',
                name: 'PatientEdit',
                component: PatientForm,
                meta: {
                    level: 2,
                    title: 'Edit Patient',
                    parent: 'PatientList',
                },
            },
            {
                path: '/patients/:id',
                name: 'PatientView',
                component: PatientForm,
                meta: {
                    level: 2,
                    title: 'Patient Details',
                    parent: 'PatientList',
                },
            },
            {
                path: '/patients/:patientId/diagnosis/new',
                name: 'PatientDiagnosisCreate',
                component: () => import('../views/PatientDiagnosisForm.vue'),
                meta: {
                    level: 3,
                    title: 'Add Diagnosis',
                    parent: 'PatientView',
                },
            },
            {
                path: '/patients/:patientId/diagnosis/:diagnosisId/edit',
                name: 'PatientDiagnosisEdit',
                component: () => import('../views/PatientDiagnosisForm.vue'),
                meta: {
                    level: 3,
                    title: 'Edit Diagnosis',
                    parent: 'PatientView',
                },
            },
            {
                path: '/patients/:patientId/treatment-plan/new',
                name: 'PatientTreatmentPlanCreate',
                component: () => import('../views/PatientTreatmentPlanForm.vue'),
                meta: {
                    level: 3,
                    title: 'Add Treatment Plan',
                    parent: 'PatientView',
                },
            },
            {
                path: '/patients/:patientId/treatment-plan/:treatmentPlanId/edit',
                name: 'PatientTreatmentPlanEdit',
                component: () => import('../views/PatientTreatmentPlanForm.vue'),
                meta: {
                    level: 3,
                    title: 'Edit Treatment Plan',
                    parent: 'PatientView',
                },
            },
            {
                path: '/patients/:patientId/dental-chart',
                name: 'PatientDentalChartView',
                component: PatientDentalChart,
                meta: {
                    level: 2,
                    title: 'Patient Dental Chart',
                    parent: 'PatientList',
                },
            },

             {
                 path: '/appointments/new',
                 name: 'AppointmentCreate',
                 component: NewAppointment,
                 meta: {
                     level: 2,
                     title: 'New Appointment',
                     parent: 'AppointmentCalendar',
                 },
             },
             {
                 path: '/appointments/:id',
                 name: 'AppointmentView',
                 component: EditAppointment,
                 meta: {
                     level: 2,
                     title: 'Appointment Details',
                     parent: 'AppointmentCalendar',
                     readOnly: true,
                 },
             },
             {
                 path: '/appointments/:id/edit',
                 name: 'AppointmentEdit',
                 component: EditAppointment,
                 meta: {
                     level: 2,
                     title: 'Edit Appointment',
                     parent: 'AppointmentCalendar',
                 },
             },

            {
                path: '/procedures/new',
                name: 'ProcedureCreate',
                component: ProcedureForm,
                meta: {
                    level: 2,
                    title: 'New Procedure',
                    parent: 'ProcedureManagement',
                },
            },
            {
                path: '/procedures/:id/edit',
                name: 'ProcedureEdit',
                component: ProcedureForm,
                meta: {
                    level: 2,
                    title: 'Edit Procedure',
                    parent: 'ProcedureManagement',
                },
            },

            {
                path: '/diagnoses/new',
                name: 'DiagnosisCreate',
                component: DiagnosisForm,
                meta: {
                    level: 2,
                    title: 'New Diagnosis',
                    parent: 'ProcedureManagement',
                },
            },
            {
                path: '/diagnoses/:id/edit',
                name: 'DiagnosisEdit',
                component: DiagnosisForm,
                meta: {
                    level: 2,
                    title: 'Edit Diagnosis',
                    parent: 'ProcedureManagement',
                },
            },

            {
                path: '/treatments/new',
                name: 'TreatmentCreate',
                component: TreatmentPlanForm,
                meta: {
                    level: 2,
                    title: 'New Treatment Plan',
                    parent: 'ProcedureManagement',
                },
            },
            {
                path: '/treatments/:id',
                name: 'TreatmentView',
                component: TreatmentPlanForm,
                meta: {
                    level: 2,
                    title: 'Treatment Plan',
                    parent: 'ProcedureManagement',
                },
            },
            {
                path: '/treatments/:id/edit',
                name: 'TreatmentEdit',
                component: TreatmentPlanForm,
                meta: {
                    level: 2,
                    title: 'Edit Treatment Plan',
                    parent: 'ProcedureManagement',
                },
            },

            {
                path: '/users/:id/edit',
                name: 'UserEdit',
                component: UserEdit,
                meta: { level: 2, title: 'Edit User', parent: 'UserManagement' },
            },
            {
                path: '/users/new',
                name: 'UserCreate',
                component: NewUser,
                meta: { level: 2, title: 'New User', parent: 'UserManagement' },
            },
            {
                path: '/staff/:id/edit',
                name: 'StaffEdit',
                component: StaffEdit,
                meta: { level: 2, title: 'Edit Staff Member', parent: 'StaffManagement' },
            },


            {
                path: '/patients/:patientId/consent/new',
                name: 'ConsentCreate',
                component: ConsentForm,
                meta: {
                    level: 3,
                    title: 'New Consent Form',
                    parent: 'PatientView',
                },
            },

            {
                path: '/prescriptions/new',
                name: 'PrescriptionCreate',
                component: PrescriptionForm,
                meta: {
                    level: 2,
                    title: 'New Prescription',
                    parent: 'ProcedureManagement',
                },
            },

            // Peer Review Routes
             {
                 path: '/peer-review/:id',
                 name: 'PeerReviewCase',
                 component: PeerReviewCase,
                 meta: {
                     level: 2,
                     title: 'Peer Review Case',
                     parent: 'PeerReviewList',
                 },
             },

             // Inventory Management Routes
             {
                 path: '/inventory',
                 name: 'InventoryList',
                 component: InventoryList,
                 meta: { level: 1, title: 'Inventory', parent: 'Dashboard' },
             },
             {
                 path: '/inventory/new',
                 name: 'InventoryCreate',
                 component: InventoryForm,
                 meta: {
                     level: 2,
                     title: 'Add Inventory Item',
                     parent: 'InventoryList',
                 },
             },
              {
                  path: '/inventory/:id/edit',
                  name: 'InventoryEdit',
                  component: InventoryForm,
                  meta: {
                      level: 2,
                      title: 'Edit Inventory Item',
                      parent: 'InventoryList',
                  },
              },
              {
                  path: '/inventory/:id',
                  name: 'InventoryView',
                  component: InventoryView,
                  meta: {
                      level: 2,
                      title: 'Inventory Item Details',
                      parent: 'InventoryList',
                  },
              },
        ],
    },
];

export default routes;
