# AGENTS.md - Development Guidelines for Dentika

## Build/Lint/Test Commands

### Frontend (Vue.js + Vite)
```bash
# Development server
npm run dev

# Production build
npm run build

# Preview production build
npm run preview
```

### Backend (Go + Fiber)
```bash
# Development server
go run main.go

# Build executable
go build -o dentika-server main.go

# Build for Linux
go build -o dentika-server main.go
```

### Testing
- No automated testing framework configured
- Manual testing via browser and API tools
- Notification testing via `NotificationTest.vue` component

### Linting
- No linting tools configured
- Code style enforced through manual review

## Code Style Guidelines

### Go Backend Conventions

#### File Structure
- `handlers/` - HTTP route handlers (one per feature)
- `models/` - Database models and schemas
- `services/` - Business logic layer
- `middleware/` - Authentication and utility middleware
- `database/` - Database connection and migrations

#### Naming Conventions
- **Functions/Methods**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase consistently
- **Structs**: PascalCase for exported fields
- **Constants**: PascalCase
- **Files**: snake_case matching package name

#### Import Organization
```go
import (
    "standard/library/packages"
    ""
    "third/party/packages"
    ""
    "local/project/packages"
)
```

#### Error Handling
- Always check and handle errors immediately
- Return appropriate HTTP status codes
- Use structured error responses: `c.Status(code).JSON(fiber.Map{"error": "message"})`
- Log errors with context using `log.Printf()`

#### Database Operations
- Use GORM for all database interactions
- Follow repository pattern through direct model usage
- Handle database errors consistently
- Use transactions for multi-step operations

### Vue.js Frontend Conventions

#### Component Structure
```vue
<template>
  <!-- Template content -->
</template>

<script>
export default {
  name: 'ComponentName',
  props: {
    // Prop definitions with validation
  },
  emits: ['event-name'],
  methods: {
    // Methods
  }
}
</script>
```

#### Naming Conventions
- **Components**: PascalCase (e.g., `AppointmentCard.vue`)
- **Files**: PascalCase for components, camelCase for utilities
- **Variables/Methods**: camelCase
- **Props**: camelCase with descriptive names
- **Events**: kebab-case for emitted events

#### Vue Composition API
- Use `<script setup>` syntax when possible
- Prefer Composition API over Options API for new components
- Use reactive refs and computed properties appropriately
- Leverage Vue 3 features (Teleport, Suspense, etc.)

#### Styling
- Use TailwindCSS for all styling
- Follow mobile-first responsive design
- Use consistent spacing and color schemes
- Prefer utility classes over custom CSS

### General Guidelines

#### Security
- Never log sensitive information (passwords, tokens, keys)
- Validate all user inputs
- Use proper authentication middleware
- Follow principle of least privilege for role-based access

#### Performance
- Optimize database queries with proper indexing
- Use pagination for large datasets
- Implement caching where appropriate
- Minimize bundle size in frontend builds

#### Code Quality
- Write self-documenting code with clear variable names
- Add comments for complex business logic only
- Follow DRY principle (Don't Repeat Yourself)
- Keep functions focused on single responsibility

#### Git Workflow
- Use descriptive commit messages
- Follow feature branch workflow
- Never commit sensitive data or credentials
- Test changes before committing

### Environment Variables

#### Backend (.env)
```
DB_HOST=localhost
DB_PORT=3306
DB_NAME=dentika
DB_USER=root
DB_PASSWORD=password
JWT_SECRET=your-secret-key
NATS_URL=nats://localhost:4222
TIMEZONE=UTC
PORT=3000
```

#### Frontend (.env)
- Configuration handled through Vite build process
- API proxying configured in `vite.config.js`

### Development Setup

1. **Database**: MySQL with GORM auto-migrations
2. **Real-time**: NATS.js for WebSocket communication
3. **Authentication**: JWT-based with role-based access control
4. **Frontend**: PWA with service worker and offline support
5. **Build**: Vite for fast development and optimized production builds

### Testing Strategy

- **Unit Tests**: Not currently implemented
- **Integration Tests**: Manual API testing
- **E2E Tests**: Manual browser testing
- **Performance**: Manual load testing

### Deployment

- Frontend builds to `/dist/` with code splitting
- Backend runs on configurable port (default: 3000)
- Static files served from `/uploads/` directory
- PWA features require HTTPS in production