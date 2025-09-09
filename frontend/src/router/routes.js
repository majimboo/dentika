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
import AppointmentForm from '../views/AppointmentForm.vue';
import ProcedureManagement from '../views/ProcedureManagement.vue';
import ProcedureForm from '../views/ProcedureForm.vue';
import DiagnosisForm from '../views/DiagnosisForm.vue';
import TreatmentPlanForm from '../views/TreatmentPlanForm.vue';
import ConsentForm from '../views/ConsentForm.vue';
import PrescriptionForm from '../views/PrescriptionForm.vue';
import ClinicManagement from '../views/ClinicManagement.vue';
import ClinicForm from '../views/ClinicForm.vue';
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
                path: '/appointments/new',
                name: 'AppointmentCreate',
                component: AppointmentForm,
                meta: {
                    level: 2,
                    title: 'New Appointment',
                    parent: 'AppointmentCalendar',
                },
            },
            {
                path: '/appointments/:id/edit',
                name: 'AppointmentEdit',
                component: AppointmentForm,
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
                path: '/consent/new',
                name: 'ConsentCreate',
                component: ConsentForm,
                meta: {
                    level: 2,
                    title: 'New Consent Form',
                    parent: 'ProcedureManagement',
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
        ],
    },
];

export default routes;
