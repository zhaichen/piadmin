import { ref, onMounted, onUnmounted } from 'vue'
import type { SystemSnapshot, WSMessage } from '@/types/system'

export function useWebSocket() {
  const data = ref<SystemSnapshot | null>(null)
  const connected = ref(false)
  let ws: WebSocket | null = null
  let reconnectTimer: number | undefined

  function getWsUrl(): string {
    const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
    const token = localStorage.getItem('piadmin_token') || ''
    return `${proto}//${location.host}/api/ws/monitor?token=${token}`
  }

  function connect() {
    if (ws) {
      ws.close()
      ws = null
    }

    ws = new WebSocket(getWsUrl())

    ws.onopen = () => {
      connected.value = true
    }

    ws.onmessage = (e: MessageEvent) => {
      const msg: WSMessage = JSON.parse(e.data)
      if (msg.type === 'snapshot') {
        data.value = msg.data
      }
    }

    ws.onclose = () => {
      connected.value = false
      ws = null
      reconnectTimer = window.setTimeout(connect, 3000)
    }

    ws.onerror = () => {
      ws?.close()
    }
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = undefined
    }
    if (ws) {
      ws.close()
      ws = null
    }
  }

  onMounted(connect)
  onUnmounted(disconnect)

  return { data, connected }
}
