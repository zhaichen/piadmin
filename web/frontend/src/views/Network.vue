<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { InterfaceInfo } from '@/types/network'
import { getNetworkInterfaces } from '@/api/client'

const interfaces = ref<InterfaceInfo[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    interfaces.value = await getNetworkInterfaces()
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="p-4 md:p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-200">Network Interfaces</h2>
      <button @click="load" :disabled="loading" class="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-sm text-gray-300 rounded-lg transition">
        Refresh
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div v-for="iface in interfaces" :key="iface.name" class="bg-gray-900 rounded-xl border border-gray-800 p-5">
        <div class="flex items-center justify-between mb-3">
          <h3 class="font-semibold text-gray-200">{{ iface.name }}</h3>
          <span class="text-xs px-2 py-0.5 rounded-full" :class="iface.flags.includes('up') ? 'bg-emerald-900 text-emerald-400' : 'bg-gray-800 text-gray-500'">
            {{ iface.flags.includes('up') ? 'UP' : 'DOWN' }}
          </span>
        </div>
        <div class="space-y-1.5 text-sm">
          <div><span class="text-gray-500">MAC:</span> <span class="text-gray-300 font-mono">{{ iface.mac || 'N/A' }}</span></div>
          <div><span class="text-gray-500">MTU:</span> <span class="text-gray-300">{{ iface.mtu }}</span></div>
          <div v-if="iface.addresses && iface.addresses.length > 0">
            <span class="text-gray-500">Addresses:</span>
            <div v-for="addr in iface.addresses" :key="addr" class="ml-4 text-gray-300 font-mono text-xs">{{ addr }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
