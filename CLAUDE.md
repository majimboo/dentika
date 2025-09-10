# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Architecture

This is a comprehensive dental clinic management system with a Vue.js frontend and Go backend:

- **Frontend**: Vue 3 + Vite + TailwindCSS application in `frontend/`
- **Backend**: Go Fiber REST API server in `server/`
- **Database**: MySQL with GORM ORM
- **Authentication**: Role-based JWT authentication with multi-tenant support

### Multi-Tenant Architecture Overview

The system is built around a multi-tenant architecture where:
- Each clinic operates as an isolated tenant with their own data
- Users are assigned to specific clinics via `clinic_id` foreign keys
- All API operations are automatically scoped to the user's clinic (except SuperAdmin)
- Middleware enforces data isolation at the database query level

### Backend Structure

**Handler Organization** (9 main handlers):
- `auth.go`: Authentication, registration, user sessions
- `clinic.go`: Clinic and branch management (SuperAdmin features)
- `patient.go`: Patient CRUD, search, medical records, auto-generates dental records
- `appointment.go`: Scheduling, status updates, conflict detection, countdown features
- `dental_record.go`: Digital dental charts, tooth-by-tooth updates, history tracking
- `procedure.go`: Procedure/diagnosis templates and appointment associations
- `users.go`: User management and avatar uploads
- `upload.go`: File upload handling
- `analytics.go`: Lead tracking, sales metrics, and patient behavior analytics

**Model Structure** (8 core model files):
- `user.go`: Extended user model with 5 role types and permission methods
- `clinic.go`: Multi-branch clinic management
- `patient.go`: Comprehensive patient profiles with emergency contacts and medical history
- `appointment.go`: Advanced scheduling with arrival tracking and countdown logic
- `dental_record.go`: JSON-based tooth condition storage with full dental chart support
- `procedure.go`: Template-based procedure and diagnosis system
- `document.go`: Digital consent forms and prescriptions with signature support
- `analytics.go`: Lead tracking, sales metrics, and patient behavior analytics

### Frontend Architecture

**Component Specialization**:
- `DentalChart.vue`: Interactive tooth visualization supporting permanent (32) and primary (20) teeth
- `ToothComponent.vue`: Individual tooth representation with condition coloring and surface indicators
- `AppointmentCountdown.vue`: Real-time countdown sidebar for upcoming appointments (30-minute window)
- `ToothDetailsPanel.vue` & `ToothEditModal.vue`: Dental record editing with condition history

**State Management** (7 specialized stores):
- `auth.js`: Extended with role-based getters and clinic context
- `clinic.js`: Multi-clinic and branch management for SuperAdmin
- `patient.js`: Patient search, pagination, dental record integration
- `appointment.js`: Calendar filtering, status management, real-time countdown data
- `dentalRecord.js`: Tooth-level CRUD operations, condition tracking, bulk updates
- `notification.js`: Toast notifications and user feedback system
- `connection.js`: API connection status and error handling

## Development Commands

### Frontend (from `frontend/` directory)
```bash
npm install          # Install dependencies
npm run dev          # Start development server on port 5173
npm run build        # Build for production
npm run preview      # Preview production build
```

### Backend (from `server/` directory)  
```bash
go mod tidy          # Download and organize dependencies
go run main.go       # Start development server on port 3000
go build -o dentika-server.exe main.go  # Build production executable
```

## API Architecture

The REST API follows dental clinic workflow patterns:

### Core Endpoint Groups
- `/api/auth/*`: Authentication and user session management
- `/api/clinics/*`: Multi-tenant clinic and branch management (SuperAdmin only)
- `/api/patients/*`: Patient management with automatic dental record creation
- `/api/appointments/*`: Scheduling with conflict detection and countdown features
- `/api/dental-records/*`: Tooth-level condition tracking with history
- `/api/procedure-templates/*` & `/api/diagnosis-templates/*`: Clinical template management
- `/api/appointments/{id}/procedures` & `/api/appointments/{id}/diagnoses`: Per-appointment clinical data

### Permission-Based Route Protection
- Role middleware enforces: `SuperAdmin` → all access, `Doctor` → clinical modifications, `Secretary` → patient/appointment management
- Clinic scoping automatically applied to all data operations (except SuperAdmin)
- Special routes like dental record editing restricted to `Doctor` role only

## Database Design

**Auto-Migration Strategy**: All 20+ models auto-migrate on server startup
**Multi-Tenant Data Isolation**: Every model (except User/Templates) includes `clinic_id` for data scoping
**Dental Record Storage**: Uses JSON fields for flexible tooth condition tracking (supports both permanent/primary teeth)
**Appointment Conflict Prevention**: Database-level constraints prevent double-booking doctors

## User Roles & Dental Workflow

- **SuperAdmin**: System management, create clinics, cross-clinic access
- **ClinicOwner**: Clinic/branch management, staff oversight, template creation  
- **Doctor**: Full clinical access - dental records, procedures, diagnoses, prescriptions
- **Secretary**: Front desk operations - scheduling, patient check-in, basic patient info
- **Assistant**: Limited clinical support access

## Key Implementation Details

### Dental Chart System
- Supports both permanent teeth (1-32) and primary teeth (A-T) numbering systems
- Each tooth stores: condition, affected surfaces, notes, update history
- Real-time visual updates with color-coded conditions and surface indicators
- Bulk update capabilities for treatment planning

### Appointment Countdown Feature
- Auto-refreshes every 30 seconds to show appointments within next 30 minutes
- Real-time countdown timers with color-coded urgency (red <5min, yellow <15min)
- Patient arrival tracking with late arrival detection (>10min past scheduled time)
- Integration with appointment status workflow

### Multi-Clinic Data Patterns
- SuperAdmin sees global clinic list, can create new clinics with auto-generated main branches
- Regular users automatically filtered to their assigned clinic data
- Branch-level appointment scheduling with doctor availability conflict checking
- Patient dental records auto-created for both permanent and primary teeth on patient creation