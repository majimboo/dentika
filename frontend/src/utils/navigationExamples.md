# Contextual Navigation Usage Guide

## Overview
For pages that can be accessed from multiple places, use the contextual navigation system instead of regular router.push. This ensures the back button takes users to where they actually came from.

## How It Works

### Before (Static Back Navigation):
- Patient details always goes back to Patient List
- Even if you came from Notifications or Appointments

### After (Contextual Back Navigation):
- Patient details goes back to where you actually came from
- Notifications → Patient Details → Back to Notifications
- Patient List → Patient Details → Back to Patient List

## Usage

### Import the composable:
```javascript
import { useNavigation } from '../composables/useNavigation'
const { navigateWithContext } = useNavigation()
```

### Use navigateWithContext instead of router.push:
```javascript
// OLD - Static navigation
router.push(`/patients/${patientId}`)

// NEW - Contextual navigation
navigateWithContext(`/patients/${patientId}`)
```

## Examples of Updated Components

### PatientList.vue
```javascript
const viewPatient = (patient) => {
  navigateWithContext(`/patients/${patient.id}`)
}
```

### Notifications.vue
```javascript
if (notification.patientId) {
  navigateWithContext(`/patients/${notification.patientId}`)
}
```

### AppointmentList.vue
```javascript
const viewAppointment = (appointment) => {
  navigateWithContext(`/appointments/${appointment.id}`)
}
```

## How Navigation Context Works

1. **Setting Context**: When you use `navigateWithContext()`, it automatically captures:
   - Current route name and path
   - Route parameters and query
   - Page title

2. **Using Context**: When the back button is pressed:
   - First checks if navigation context exists
   - If yes, navigates to the stored context
   - If no, falls back to parent route logic
   - Last resort: browser back button

3. **Context Lifecycle**:
   - Context is set when navigating with `navigateWithContext()`
   - Context is cleared after being used for back navigation
   - Context is automatically cleared on direct navigation

## When to Use

Use `navigateWithContext()` for pages that can be accessed from multiple places:

- ✅ Patient details (from list, appointments, notifications)
- ✅ Appointment details (from calendar, list, notifications)
- ✅ Inventory details (from list, shop, orders)
- ❌ New patient form (always goes back to patient list)
- ❌ Settings pages (always go back to dashboard)

## Fallback Behavior

The system gracefully falls back in this order:
1. Navigation context (where you came from)
2. Parent route (defined in route meta)
3. Browser back button (router.go(-1))