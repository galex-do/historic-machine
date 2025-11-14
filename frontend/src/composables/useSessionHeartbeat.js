import { ref, watch, onBeforeUnmount } from 'vue'
import authService from '@/services/authService.js'

let heartbeatInterval = null

export function useSessionHeartbeat() {
  const is_running = ref(false)

  const start_heartbeat = () => {
    if (heartbeatInterval || !authService.isAuthenticated()) {
      return
    }

    is_running.value = true
    heartbeatInterval = setInterval(() => {
      if (authService.isAuthenticated()) {
        authService.sendHeartbeat()
      } else {
        stop_heartbeat()
      }
    }, 60000)
  }

  const stop_heartbeat = () => {
    if (heartbeatInterval) {
      clearInterval(heartbeatInterval)
      heartbeatInterval = null
      is_running.value = false
    }
  }

  onBeforeUnmount(() => {
    stop_heartbeat()
  })

  return {
    start_heartbeat,
    stop_heartbeat,
    is_running
  }
}
