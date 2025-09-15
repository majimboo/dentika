# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Dentika is a comprehensive dental practice management application with a Go backend and Vue.js frontend. The system supports multi-clinic operations with user management, patient records, appointments, inventory, peer reviews, and a real-time notification system.

## Architecture

### Backend (Go + Fiber)
- **Framework**: Go with Fiber web framework
- **Database**: MySQL with GORM ORM
- **Real-time**: NATS.js for WebSocket communication
- **Authentication**: JWT-based authentication with role-based access control
- **File Structure**:
  - `main.go` - Application entry point with server setup
  - `handlers/` - HTTP route handlers organized by feature
  - `models/` - Database models and schemas
  - `services/` - Business logic layer
  - `middleware/` - Authentication and CORS middleware
  - `database/` - Database connection and migrations

### Frontend (Vue.js 3 PWA)
- **Framework**: Vue 3 with Composition API
- **App Type**: Progressive Web App (PWA) with mobile-first native design
- **State Management**: Pinia stores
- **UI**: TailwindCSS with custom components designed for mobile-first experience
- **Routing**: Vue Router with role-based guards
- **Build Tool**: Vite with PWA support and service worker
- **Real-time**: NATS.js WebSocket integration
- **File Structure**:
  - `src/views/` - Page components organized by feature
  - `src/components/` - Reusable UI components
  - `src/stores/` - Pinia state management
  - `src/composables/` - Vue composables (including NATS integration)
  - `src/layouts/` - Layout components

### Key Features
- **Multi-clinic Support**: Clinics can have multiple branches
- **Role-based Access**: Super Admin, Admin, Doctor, Secretary, Assistant roles
- **Real-time Notifications**: NATS-powered notification system with database persistence
- **Patient Management**: Comprehensive patient records with dental charts
- **Appointment System**: Calendar views with scheduling and management
- **Inventory Management**: Multi-level inventory with platform-clinic ordering
- **Peer Review System**: Case review and collaboration features

## Development Commands

### Frontend (from `/frontend/`)
- **Development**: `npm run dev` - Start Vite dev server on port 5173
- **Build**: `npm run build` - Build for production
- **Preview**: `npm run preview` - Preview production build

### Backend (from `/server/`)
- **Development**: `go run main.go` - Start Go server on port 9483
- **Build**: `go build -o server.exe main.go` - Build executable
- **Build (Linux)**: `go build -o dentika-server main.go` - Build for Linux

## Database Setup

The application uses MySQL. Connection details are in `.env`:
- Database models are in `server/models/`
- GORM handles migrations automatically on startup
- Seeded data includes 3 clinics with users and patients (see `SEEDED_DATA.md`)

## NATS Configuration

Real-time features use NATS with WebSocket:
- **Server**: `nats://localhost:4222` (configurable via `NATS_URL`)
- **Frontend WebSocket**:
  - Local: `wss://localhost:9222`
  - Production: `wss://nats.majidarif.com:9222`
- **Subject Structure**: `dentika.{scope}.{id}.{feature}` (see `NOTIFICATION_SYSTEM.md`)

## Key Development Patterns

### Authentication Flow
1. JWT tokens stored in localStorage
2. Role-based route guards in Vue Router
3. Middleware validates tokens on backend routes
4. User context available via Pinia stores

### Notification System
- Database-persisted notifications with real-time delivery
- NATS subjects for user, clinic, and system-wide notifications
- Vue composable `useNats.js` handles WebSocket connections
- Notification store manages state and API interactions

### Component Architecture
- Feature-based organization (appointments, patients, inventory, etc.)
- Shared components in `/components/`
- Role-specific sidebars and navigation
- **Mobile-first design**: Native mobile styling with touch-optimized interactions
- **Responsive navigation**: Desktop sidebars transform to mobile bottom navigation
- **PWA UI patterns**: Native app-like interface with proper touch targets and gestures

### API Patterns
- RESTful endpoints with consistent response format
- Middleware for authentication and CORS
- Error handling with proper HTTP status codes
- File uploads handled via `/uploads/` endpoint

## Environment Variables

### Backend (`.env`)
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` - Database config
- `JWT_SECRET` - JWT signing key
- `NATS_URL` - NATS server URL
- `TIMEZONE` - Application timezone (default: UTC)

### Frontend (`.env`)
- Build-time configuration handled by Vite
- API proxying configured in `vite.config.js`

## Production Considerations

- Frontend builds to `/dist/` with code splitting by feature
- Assets are chunked and hashed for caching
- **PWA features**: Service worker, offline support, installable app experience
- **Mobile optimization**: Touch-friendly interface, native mobile interactions
- **App manifest**: Configured for native app-like installation on mobile devices
- Server runs on port 9483, frontend on 5173 (development)
- Proxy configuration routes `/api` and `/uploads` to backend

## Testing

No automated testing framework is currently configured. Manual testing is done via:
- Browser-based UI testing
- API testing via development tools
- Notification testing via `NotificationTest.vue` component