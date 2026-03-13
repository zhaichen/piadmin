<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { ProcessInfo } from '@/types/process'
import { getProcesses, killProcess } from '@/api/client'

const processes = ref<ProcessInfo[]>([])
const loading = ref(false)
const filter = ref('')

async function load() {
  loading.value = true
  try {
    processes.value = await getProcesses()
  } finally {
    loading.value = false
  }
}

async function handleKill(pid: number, force: boolean) {
  if (!confirm(`${force ? 'Force kill' : 'Terminate'} process ${pid}?`)) return
  await killProcess(pid, force)
  await load()
}

const filtered = () => {
  if (!filter.value) return processes.value
  const q = filter.value.toLowerCase()
  return processes.value.filter(
    p => p.name.toLowerCase().includes(q) || p.command.toLowerCase().includes(q) || String(p.pid).includes(q)
  )
}

onMounted(load)
</script>

<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-200">Processes</h2>
      <div class="flex gap-3">
        <input
          v-model="filter"
          placeholder="Filter..."
          class="px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-lg text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-emerald-500"
        />
        <button @click="load" :disabled="loading" class="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-sm text-gray-300 rounded-lg transition">
          Refresh
        </button>
      </div>
    </div>

    <div class="bg-gray-900 rounded-xl border border-gray-800 overflow-hidden">
      <table class="w-full text-sm">
        <thead>
          <tr class="text-gray-500 text-left border-b border-gray-800">
            <th class="px-4 py-3 font-medium">PID</th>
            <th class="px-4 py-3 font-medium">Name</th>
            <th class="px-4 py-3 font-medium">User</th>
            <th class="px-4 py-3 font-medium text-right">CPU%</th>
            <th class="px-4 py-3 font-medium text-right">Mem%</th>
            <th class="px-4 py-3 font-medium">Status</th>
            <th class="px-4 py-3 font-medium w-20"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in filtered()" :key="p.pid" class="border-t border-gray-800/50 hover:bg-gray-800/30">
            <td class="px-4 py-2 text-gray-400 font-mono">{{ p.pid }}</td>
            <td class="px-4 py-2 text-gray-200">{{ p.name }}</td>
            <td class="px-4 py-2 text-gray-400">{{ p.user }}</td>
            <td class="px-4 py-2 text-right" :class="p.cpu > 50 ? 'text-amber-400' : 'text-gray-300'">
              {{ p.cpu.toFixed(1) }}
            </td>
            <td class="px-4 py-2 text-right" :class="p.memory > 10 ? 'text-amber-400' : 'text-gray-300'">
              {{ p.memory.toFixed(1) }}
            </td>
            <td class="px-4 py-2">
              <span class="text-xs px-1.5 py-0.5 rounded" :class="{
                'bg-emerald-900/50 text-emerald-400': p.status === 'R',
                'bg-blue-900/50 text-blue-400': p.status === 'S',
                'bg-gray-800 text-gray-500': !['R','S'].includes(p.status)
              }">{{ p.status }}</span>
            </td>
            <td class="px-4 py-2">
              <button @click="handleKill(p.pid, false)" class="text-red-500 hover:text-red-400 text-xs mr-2">Kill</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
