import Login from '../views/Login.vue';
import Register from '../views/Register.vue';
import Dashboard from '../views/Dashboard.vue';
import UserManagementWrapper from '../views/UserManagementWrapper.vue';
import UserList from '../views/UserList.vue'; // Keep for backward compatibility
import UserEdit from '../views/UserEdit.vue';
import StaffEdit from '../views/StaffEdit.vue';
import Agenda from '../views/Agenda.vue';
import NewStaff from '../views/NewStaff.vue';
import NewUser from '../views/NewUser.vue';
import PatientList from '../views/PatientList.vue';
import PatientForm from '../views/PatientForm.vue';
import AppointmentCalendar from '../views/AppointmentCalendar.vue';
import AppointmentList from '../views/AppointmentList.vue';
import NewAppointment from '../views/NewAppointment.vue';
import EditAppointment from '../views/EditAppointment.vue';
import ProcedureManagement from '../views/ProcedureManagement.vue';
import ProcedureForm from '../views/ProcedureForm.vue';
import DiagnosisForm from '../views/DiagnosisForm.vue';
import TreatmentPlanForm from '../views/TreatmentPlanForm.vue';
import ConsentForm from '../views/ConsentForm.vue';
import PrescriptionForm from '../views/PrescriptionForm.vue';
import PatientDentalChart from '../views/PatientDentalChart.vue';
import ClinicManagement from '../views/ClinicManagement.vue';
import ClinicForm from '../views/ClinicForm.vue';
import PeerReviewList from '../views/PeerReviewList.vue';
import PeerReviewCase from '../views/PeerReviewCase.vue';
import PeerReviewCreate from '../views/PeerReviewCreate.vue';
import InventoryList from '../views/InventoryList.vue';
import InventoryForm from '../views/InventoryForm.vue';
import AppLayout from '../layouts/AppLayout.vue';

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: { hideLayout: true },
    },
    {
        path: '/register',
        name: 'Register',
        component: Register,
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
        ],
    },
];

export default routes;
