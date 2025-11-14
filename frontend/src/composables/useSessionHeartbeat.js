import { ref, onBeforeUnmount } from 'vue'
import authService from '@/services/authService.js'

let heartbeatInterval = null

function generate_session_id() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

function get_anonymous_session_id() {
  let sessionId = localStorage.getItem('historia_anonymous_session_id')
  if (!sessionId) {
    sessionId = generate_session_id()
    localStorage.setItem('historia_anonymous_session_id', sessionId)
  }
  return sessionId
}

export function useSessionHeartbeat() {
  const is_running = ref(false)

  const send_heartbeat = async () => {
    if (authService.isAuthenticated()) {
      await authService.sendHeartbeat()
    } else {
      await authService.sendAnonymousHeartbeat(get_anonymous_session_id())
    }
  }

  const start_heartbeat = () => {
    if (heartbeatInterval) {
      return
    }

    is_running.value = true
    
    send_heartbeat()
    
    heartbeatInterval = setInterval(() => {
      send_heartbeat()
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
    is_running,
    get_anonymous_session_id
  }
}
