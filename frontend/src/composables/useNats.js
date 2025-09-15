import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useNotificationStore } from '../stores/notification'
import {
    wsconnect
} from '@nats-io/nats-core'

// Singleton NATS instance
let natsInstance = null
let isInitializing = false

class NATSService {
  constructor() {
    this.connection = null
    this.isConnected = false
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = -1
    this.subscriptions = new Map()
    this.heartbeatInterval = null
    this.authStore = null
    this.notificationStore = null
    this.isInitialized = false
  }

  async initialize(authStore, notificationStore) {
    // Prevent multiple simultaneous initializations
    if (isInitializing) {
      // Wait for current initialization to complete
      while (isInitializing) {
        await new Promise(resolve => setTimeout(resolve, 100))
      }
      return
    }

    if (this.isInitialized && this.isConnected && this.connection) {
      return
    }

    isInitializing = true

    try {
      // Set stores
      this.authStore = authStore
      this.notificationStore = notificationStore

      if (!authStore.token) {
        console.warn('Cannot connect to NATS: No authentication token')
        return
      }

      await this.connect()
      this.isInitialized = true
    } finally {
      isInitializing = false
    }
  }

  async connect() {
    // Don't cleanup if we're already connected
    if (this.isConnected && this.connection) {
      return
    }

    // Only cleanup if we have an existing broken connection
    if (this.connection && !this.isConnected) {
      await this.disconnect()
    }

    let wsUrl
    if (window.location.hostname === 'localhost') {
      wsUrl = `wss://${window.location.hostname}:9222`
    } else {
      wsUrl = 'wss://nats.majidarif.com:9222'
    }

    try {
      this.connection = await wsconnect({
        servers: [wsUrl],
        maxReconnectAttempts: this.maxReconnectAttempts,
        reconnectTimeWait: 1000,
        maxReconnectTimeWait: 30000,
        reconnectJitter: 100,
        reconnectJitterTLS: 100,
        token: this.authStore.token,
        debug: false,
        verbose: false,
      })

      console.log('NATS connected successfully')
      this.isConnected = true
      this.reconnectAttempts = 0

      this.setupConnectionEvents()
      await this.setupSubscriptions()
      this.startHeartbeat()

    } catch (error) {
      console.error('Failed to connect to NATS:', error)
      this.isConnected = false
      this.scheduleReconnect()
    }
  }

  setupConnectionEvents() {
    if (!this.connection) return

    this.connection.closed().then(() => {
      console.log('NATS connection closed')
      this.isConnected = false
      this.stopHeartbeat()
      this.clearSubscriptions()

      if (this.reconnectAttempts < this.maxReconnectAttempts || this.maxReconnectAttempts === -1) {
        this.scheduleReconnect()
      }
    }).catch((error) => {
      console.error('NATS connection error:', error)
      this.isConnected = false
      this.scheduleReconnect()
    })
  }

  async setupSubscriptions() {
    if (!this.connection || !this.authStore.user) return

    const userId = this.authStore.user.id
    const clinicId = this.authStore.user.clinic_id

    try {
      await this.subscribe(`dentika.user.${userId}.notifications`, this.handleUserNotification.bind(this))
      await this.subscribe(`dentika.user.${userId}.appointments`, this.handleAppointmentUpdate.bind(this))
      await this.subscribe(`dentika.user.${userId}.patients`, this.handlePatientUpdate.bind(this))

      if (clinicId) {
        await this.subscribe(`dentika.clinic.${clinicId}.notifications`, this.handleClinicNotification.bind(this))
        await this.subscribe(`dentika.clinic.${clinicId}.appointments`, this.handleAppointmentUpdate.bind(this))
        await this.subscribe(`dentika.clinic.${clinicId}.patients`, this.handlePatientUpdate.bind(this))
      }

      await this.subscribe('dentika.system.notifications', this.handleSystemNotification.bind(this))
      await this.subscribe('dentika.system.alerts', this.handleSystemAlert.bind(this))

      console.log('NATS subscriptions setup complete')
    } catch (error) {
      console.error('Failed to setup NATS subscriptions:', error)
    }
  }

  async subscribe(subject, handler) {
    if (!this.connection) return

    try {
      const subscription = this.connection.subscribe(subject)
      this.subscriptions.set(subject, subscription)

      ;(async () => {
        for await (const message of subscription) {
          try {
            const data = JSON.parse(new TextDecoder().decode(message.data))
            handler(data)
          } catch (error) {
            console.error(`Failed to parse message from ${subject}:`, error)
          }
        }
      })()

      console.log(`Subscribed to: ${subject}`)
    } catch (error) {
      console.error(`Failed to subscribe to ${subject}:`, error)
    }
  }

  clearSubscriptions() {
    this.subscriptions.forEach((subscription, subject) => {
      try {
        subscription.unsubscribe()
      } catch (error) {
        console.error(`Error unsubscribing from ${subject}:`, error)
      }
    })
    this.subscriptions.clear()
  }

  async disconnect() {
    this.stopHeartbeat()
    this.clearSubscriptions()

    if (this.connection) {
      try {
        await this.connection.close()
        console.log('NATS connection closed gracefully')
      } catch (error) {
        console.error('Error closing NATS connection:', error)
      }
    }

    this.connection = null
    this.isConnected = false
    this.reconnectAttempts = 0
  }

  scheduleReconnect() {
    if (this.maxReconnectAttempts !== -1 && this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log('Max reconnection attempts reached')
      return
    }

    this.reconnectAttempts++
    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts - 1), 30000)

    console.log(`Scheduling NATS reconnection attempt ${this.reconnectAttempts} in ${delay}ms`)

    setTimeout(() => {
      this.connect()
    }, delay)
  }

  startHeartbeat() {
    this.heartbeatInterval = setInterval(async () => {
      if (this.connection && this.isConnected) {
        try {
          await this.publish('dentika.heartbeat', {
            user_id: this.authStore.user?.id,
            timestamp: Date.now()
          })
        } catch (error) {
          console.error('Heartbeat failed:', error)
        }
      }
    }, 30000)
  }

  stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }
  }

  async publish(subject, data) {
    if (!this.connection) {
      console.warn('Cannot publish message: NATS not connected')
      return false
    }

    try {
      const payload = new TextEncoder().encode(JSON.stringify(data))
      this.connection.publish(subject, payload)
      return true
    } catch (error) {
      console.error(`Failed to publish to ${subject}:`, error)
      return false
    }
  }

  // Message handlers
  handleUserNotification(data) {
    this.notificationStore.addNotificationFromNATS(data)
    window.dispatchEvent(new CustomEvent('app-notification', { detail: data }))
  }

  handleAppointmentUpdate(data) {
    switch (data.type) {
      case 'appointment_update':
        window.dispatchEvent(new CustomEvent('appointment-updated', { detail: data.payload }))
        break
      case 'appointment_reminder':
        this.handleAppointmentReminder(data.payload)
        break
      default:
        console.log('Unknown appointment message type:', data.type)
    }
  }

  handleAppointmentReminder(reminder) {
    if ('Notification' in window && Notification.permission === 'granted') {
      new Notification(`Appointment Reminder`, {
        body: `${reminder.patient_name} has an appointment in ${reminder.minutes_until} minutes`,
        icon: '/favicon.ico',
        tag: `reminder-${reminder.appointment_id}`
      })
    }

    window.dispatchEvent(new CustomEvent('appointment-reminder', { detail: reminder }))
  }

  handlePatientUpdate(data) {
    window.dispatchEvent(new CustomEvent('patient-updated', { detail: data.payload }))
  }

  handleClinicNotification(data) {
    // Add to notification store (same as user notifications)
    this.notificationStore.addNotificationFromNATS(data)

    window.dispatchEvent(new CustomEvent('clinic-notification', { detail: data }))

    if ('Notification' in window && Notification.permission === 'granted') {
      new Notification(data.title || 'Clinic Notification', {
        body: data.message,
        icon: '/favicon.svg',
        tag: `clinic-${data.type}-${Date.now()}`
      })
    }
  }

  handleSystemNotification(data) {
    window.dispatchEvent(new CustomEvent('system-notification', { detail: data }))

    if ('Notification' in window && Notification.permission === 'granted') {
      new Notification(data.title || 'System Notification', {
        body: data.message,
        icon: '/favicon.svg',
        tag: `system-${data.type}-${Date.now()}`
      })
    }
  }

  handleSystemAlert(data) {
    window.dispatchEvent(new CustomEvent('system-alert', { detail: data }))
  }

  async requestNotificationPermission() {
    if ('Notification' in window && Notification.permission === 'default') {
      const permission = await Notification.requestPermission()
      return permission === 'granted'
    }
    return Notification.permission === 'granted'
  }
}

export function useNats() {
  if (!natsInstance) {
    natsInstance = new NATSService()
  }

  const authStore = useAuthStore()
  const notificationStore = useNotificationStore()

  const connect = async () => {
    await natsInstance.initialize(authStore, notificationStore)
  }

  const disconnect = async () => {
    await natsInstance.disconnect()
  }

  const publish = async (subject, data) => {
    return await natsInstance.publish(subject, data)
  }

  const requestNotificationPermission = async () => {
    return await natsInstance.requestNotificationPermission()
  }

  // Utility functions for publishing common events
  const publishAppointmentUpdate = async (appointment) => {
    const subject = `dentika.clinic.${appointment.clinic_id}.appointments`
    return await publish(subject, {
      type: 'appointment_update',
      payload: appointment
    })
  }

  const publishPatientUpdate = async (patient) => {
    const subject = `dentika.clinic.${patient.clinic_id}.patients`
    return await publish(subject, {
      type: 'patient_update',
      payload: patient
    })
  }

  const publishUserNotification = async (userId, notification) => {
    const subject = `dentika.user.${userId}.notifications`
    return await publish(subject, notification)
  }

  const publishClinicNotification = async (clinicId, notification) => {
    const subject = `dentika.clinic.${clinicId}.notifications`
    return await publish(subject, notification)
  }

  const publishSystemNotification = async (notification) => {
    const subject = 'dentika.system.notifications'
    return await publish(subject, notification)
  }

  // Lifecycle management
  onMounted(() => {
    if (authStore.token && authStore.user) {
      notificationStore.initialize()
      connect()
      requestNotificationPermission()
    }
  })

  onUnmounted(() => {
    // Don't disconnect on unmount since this is a singleton
    // The connection should stay alive across component lifecycles
  })

  return {
    connection: ref(natsInstance.connection),
    isConnected: ref(natsInstance.isConnected),
    connect,
    disconnect,
    publish,
    requestNotificationPermission,
    // Utility publishing functions
    publishAppointmentUpdate,
    publishPatientUpdate,
    publishUserNotification,
    publishClinicNotification,
    publishSystemNotification
  }
}