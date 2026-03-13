<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import '@xterm/xterm/css/xterm.css'
import { getAuthToken } from '@/api/client'

const terminalRef = ref<HTMLDivElement>()
let term: Terminal | null = null
let ws: WebSocket | null = null
let fitAddon: FitAddon | null = null

function connect() {
  if (!terminalRef.value) return

  term = new Terminal({
    cursorBlink: true,
    fontSize: 14,
    fontFamily: 'Menlo, Monaco, "Courier New", monospace',
    theme: {
      background: '#0a0a0a',
      foreground: '#e2e8f0',
      cursor: '#10b981',
      selectionBackground: '#334155',
    },
  })

  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.loadAddon(new WebLinksAddon())
  term.open(terminalRef.value)
  fitAddon.fit()

  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  const token = getAuthToken()
  ws = new WebSocket(`${proto}//${location.host}/api/ws/terminal?token=${token}`)
  ws.binaryType = 'arraybuffer'

  ws.onopen = () => {
    // send initial size
    ws!.send(JSON.stringify({ type: 'resize', cols: term!.cols, rows: term!.rows }))
  }

  ws.onmessage = (e) => {
    if (e.data instanceof ArrayBuffer) {
      term!.write(new Uint8Array(e.data))
    } else {
      term!.write(e.data)
    }
  }

  ws.onclose = () => {
    term!.write('\r\n\x1b[31mConnection closed.\x1b[0m\r\n')
  }

  term.onData((data) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(data)
    }
  })

  term.onResize(({ cols, rows }) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'resize', cols, rows }))
    }
  })
}

function handleResize() {
  fitAddon?.fit()
}

onMounted(() => {
  connect()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  ws?.close()
  term?.dispose()
})
</script>

<template>
  <div class="h-full flex flex-col">
    <div class="px-6 py-3 border-b border-gray-800 flex items-center justify-between">
      <h2 class="text-lg font-semibold text-gray-200">Terminal</h2>
    </div>
    <div ref="terminalRef" class="flex-1 p-2 bg-[#0a0a0a]"></div>
  </div>
</template>
