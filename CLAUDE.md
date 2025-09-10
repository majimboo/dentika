# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Dentika is a comprehensive dental practice management application with:
- **Backend**: Go server using Fiber web framework with GORM ORM and MySQL database
- **Frontend**: Vue.js 3 SPA with Vite, TailwindCSS, Pinia state management, and Vue Router

## Development Commands

### Backend (Go server in /server)
```bash
# Development (from server directory)
go run main.go

# Build production binary
go build -o dentika-server.exe main.go

# Install dependencies
go mod tidy
```

### Frontend (Vue.js in /frontend)
```bash
# Development server (from frontend directory)
npm run dev

# Production build
npm run build

# Preview production build
npm run preview

# Install dependencies
npm install
```

## Architecture & Key Components

### Backend Structure (`/server`)
- **main.go**: Application entry point with Fiber app setup, middleware, and database migration
- **models/**: GORM models for all entities (User, Patient, Appointment, DentalRecord, Inventory, etc.)
- **handlers/**: HTTP route handlers organized by domain (auth, appointment, dental_record, consent, inventory, etc.)
- **database/**: Database connection setup using GORM with MySQL
- **middleware/**: Custom middleware (likely auth and CORS)

### Frontend Structure (`/frontend/src`)
- **main.js**: Vue app bootstrap with FontAwesome icons, Pinia stores, and router setup
- **router/**: Vue Router configuration with auth guards
- **views/**: Page components organized by feature areas
- **stores/**: Pinia stores for state management (auth, appointment, clinic, dental records, etc.)
- **components/**: Reusable Vue components
- **api/**: Axios configuration with authentication and error handling interceptors
- **layouts/**: Layout components including AppLayout for authenticated routes

### Database Models
Key entities include: User, AuthToken, Clinic, Branch, Patient, PatientDocument, Appointment, DentalRecord, ConsentTemplate, ConsentForm, Prescription, Inventory items, and Analytics.

## Configuration

### Environment Variables
- **Backend (.env)**: DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD, PORT, TIMEZONE
- **Frontend (.env)**: VITE_API_URL (defaults to http://localhost:3000)

### Development Setup
- Backend runs on port specified in .env (default likely 9483 based on Vite proxy config)
- Frontend dev server runs on port 5173 with proxy to backend API
- Database: MySQL with auto-migration on server startup

### API Communication
- Frontend uses Axios with base URL configuration
- API routes prefixed with `/api`
- Authentication via JWT tokens stored in localStorage
- Automatic token refresh and redirect to login on 401 responses

## Key Features & Domains

1. **Authentication & User Management**: Login, register, user roles, staff management
2. **Patient Management**: Patient records, documents, dental charts
3. **Appointment Scheduling**: Calendar view, appointment management, reminders
4. **Clinical Features**: Dental procedures, diagnosis, treatment plans, consent forms
5. **Inventory Management**: Stock tracking, restock alerts, item management
6. **Analytics**: Patient analytics, daily sales tracking
7. **Clinic Management**: Multi-clinic/branch support

## Development Notes

- Uses FontAwesome icons extensively (imported in main.js)
- TailwindCSS for styling with Vite plugin
- Pinia for centralized state management
- Vue Router with authentication guards
- GORM auto-migration runs on server startup
- Default admin user and sample data seeding on first run
- File upload handling for patient documents
- Real-time connection status monitoring via connection store