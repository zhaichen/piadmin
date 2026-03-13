<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { ServiceInfo } from '@/types/service'
import { getServices, serviceAction } from '@/api/client'

const services = ref<ServiceInfo[]>([])
const loading = ref(false)
const filter = ref('')

async function load() {
  loading.value = true
  try {
    services.value = await getServices() || []
  } finally {
    loading.value = false
  }
}

async function handleAction(name: string, action: string) {
  await serviceAction(name, action)
  await load()
}

const filtered = () => {
  if (!filter.value) return services.value
  const q = filter.value.toLowerCase()
  return services.value.filter(s => s.name.toLowerCase().includes(q) || s.description.toLowerCase().includes(q))
}

onMounted(load)
</script>

<template>
  <div class="p-4 md:p-6">
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 mb-6">
      <h2 class="text-lg font-semibold text-gray-200">Services (systemd)</h2>
      <div class="flex gap-3 w-full sm:w-auto">
        <input
          v-model="filter"
          placeholder="Filter..."
          class="flex-1 sm:flex-initial px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-lg text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-emerald-500"
        />
        <button @click="load" :disabled="loading" class="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-sm text-gray-300 rounded-lg transition shrink-0">
          Refresh
        </button>
      </div>
    </div>

    <div v-if="services.length === 0 && !loading" class="text-center text-gray-500 py-20">
      systemd services are only available on Linux
    </div>

    <div v-else class="bg-gray-900 rounded-xl border border-gray-800 overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="text-gray-500 text-left border-b border-gray-800">
            <th class="px-4 py-3 font-medium">Service</th>
            <th class="px-4 py-3 font-medium">Description</th>
            <th class="px-4 py-3 font-medium">State</th>
            <th class="px-4 py-3 font-medium">Sub</th>
            <th class="px-4 py-3 font-medium w-40">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in filtered()" :key="s.name" class="border-t border-gray-800/50 hover:bg-gray-800/30">
            <td class="px-4 py-2 text-gray-200 font-mono text-xs">{{ s.name }}</td>
            <td class="px-4 py-2 text-gray-400 truncate max-w-xs">{{ s.description }}</td>
            <td class="px-4 py-2">
              <span class="text-xs px-1.5 py-0.5 rounded" :class="{
                'bg-emerald-900/50 text-emerald-400': s.active_state === 'active',
                'bg-red-900/50 text-red-400': s.active_state === 'failed',
                'bg-gray-800 text-gray-500': !['active','failed'].includes(s.active_state)
              }">{{ s.active_state }}</span>
            </td>
            <td class="px-4 py-2 text-gray-400 text-xs">{{ s.sub_state }}</td>
            <td class="px-4 py-2 flex gap-1">
              <button @click="handleAction(s.name, 'start')" class="text-xs px-2 py-1 bg-emerald-900/30 text-emerald-400 rounded hover:bg-emerald-900/50 transition">Start</button>
              <button @click="handleAction(s.name, 'stop')" class="text-xs px-2 py-1 bg-red-900/30 text-red-400 rounded hover:bg-red-900/50 transition">Stop</button>
              <button @click="handleAction(s.name, 'restart')" class="text-xs px-2 py-1 bg-blue-900/30 text-blue-400 rounded hover:bg-blue-900/50 transition">Restart</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
