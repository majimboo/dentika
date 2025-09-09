import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '../stores/auth'

export function useWebSocket() {
  const socket = ref(null)
  const isConnected = ref(false)
  const reconnectAttempts = ref(0)
  const maxReconnectAttempts = 5
  const reconnectDelay = ref(1000) // Start with 1 second
  let reconnectTimeout = null

  const authStore = useAuthStore()

  const connect = () => {
    if (!authStore.token) {
      console.warn('Cannot connect to WebSocket: No authentication token')
      return
    }

    // Determine WebSocket URL based on environment
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const wsUrl = `${protocol}//${host}/ws`

    try {
      socket.value = new WebSocket(wsUrl)

      socket.value.onopen = () => {
        console.log('WebSocket connected')
        isConnected.value = true
        reconnectAttempts.value = 0
        reconnectDelay.value = 1000

        // Send authentication
        send({
          type: 'auth',
          token: authStore.token
        })
      }

      socket.value.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          handleMessage(data)
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      socket.value.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason)
        isConnected.value = false
        
        // Attempt reconnection if not a deliberate close
        if (event.code !== 1000 && reconnectAttempts.value < maxReconnectAttempts) {
          scheduleReconnect()
        }
      }

      socket.value.onerror = (error) => {
        console.error('WebSocket error:', error)
        isConnected.value = false
      }

    } catch (error) {
      console.error('Failed to create WebSocket connection:', error)
    }
  }

  const disconnect = () => {
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout)
      reconnectTimeout = null
    }

    if (socket.value && socket.value.readyState === WebSocket.OPEN) {
      socket.value.close(1000, 'Manual disconnect')
    }

    socket.value = null
    isConnected.value = false
    reconnectAttempts.value = 0
  }

  const scheduleReconnect = () => {
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout)
    }

    reconnectTimeout = setTimeout(() => {
      reconnectAttempts.value++
      reconnectDelay.value = Math.min(reconnectDelay.value * 2, 30000) // Max 30 seconds
      console.log(`Attempting to reconnect... (${reconnectAttempts.value}/${maxReconnectAttempts})`)
      connect()
    }, reconnectDelay.value)
  }

  const send = (data) => {
    if (socket.value && socket.value.readyState === WebSocket.OPEN) {
      socket.value.send(JSON.stringify(data))
      return true
    } else {
      console.warn('Cannot send message: WebSocket not connected')
      return false
    }
  }

  const handleMessage = (data) => {
    switch (data.type) {
      case 'appointment_update':
        handleAppointmentUpdate(data.payload)
        break
      case 'appointment_reminder':
        handleAppointmentReminder(data.payload)
        break
      case 'patient_update':
        handlePatientUpdate(data.payload)
        break
      case 'notification':
        handleNotification(data.payload)
        break
      case 'system_alert':
        handleSystemAlert(data.payload)
        break
      default:
        console.log('Unknown message type:', data.type, data)
    }
  }

  const handleAppointmentUpdate = (appointment) => {
    // Emit custom event that components can listen to
    window.dispatchEvent(new CustomEvent('appointment-updated', {
      detail: appointment
    }))
  }

  const handleAppointmentReminder = (reminder) => {
    // Show notification for upcoming appointments
    if ('Notification' in window && Notification.permission === 'granted') {
      new Notification(`Appointment Reminder`, {
        body: `${reminder.patient_name} has an appointment in ${reminder.minutes_until} minutes`,
        icon: '/favicon.ico',
        tag: `reminder-${reminder.appointment_id}`
      })
    }

    // Emit event for in-app handling
    window.dispatchEvent(new CustomEvent('appointment-reminder', {
      detail: reminder
    }))
  }

  const handlePatientUpdate = (patient) => {
    window.dispatchEvent(new CustomEvent('patient-updated', {
      detail: patient
    }))
  }

  const handleNotification = (notification) => {
    window.dispatchEvent(new CustomEvent('app-notification', {
      detail: notification
    }))
  }

  const handleSystemAlert = (alert) => {
    window.dispatchEvent(new CustomEvent('system-alert', {
      detail: alert
    }))
  }

  // Request notification permission
  const requestNotificationPermission = async () => {
    if ('Notification' in window && Notification.permission === 'default') {
      const permission = await Notification.requestPermission()
      return permission === 'granted'
    }
    return Notification.permission === 'granted'
  }

  // Lifecycle management
  onMounted(() => {
    if (authStore.token) {
      connect()
      requestNotificationPermission()
    }
  })

  onUnmounted(() => {
    disconnect()
  })

  return {
    socket,
    isConnected,
    connect,
    disconnect,
    send,
    requestNotificationPermission
  }
}