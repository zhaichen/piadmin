<script setup lang="ts">
import { useWebSocket } from '@/composables/useWebSocket'
import SystemInfo from '@/components/dashboard/SystemInfo.vue'
import CpuGauge from '@/components/dashboard/CpuGauge.vue'
import MemoryBar from '@/components/dashboard/MemoryBar.vue'
import DiskUsage from '@/components/dashboard/DiskUsage.vue'
import NetworkTraffic from '@/components/dashboard/NetworkTraffic.vue'
import Temperature from '@/components/dashboard/Temperature.vue'
import Voltage from '@/components/dashboard/Voltage.vue'

const { data, connected } = useWebSocket()
</script>

<template>
  <div class="p-4 md:p-6">
    <div class="flex items-center gap-3 mb-6">
      <h2 class="text-lg font-semibold text-gray-200">Dashboard</h2>
      <span
        class="text-xs px-2 py-0.5 rounded-full"
        :class="connected ? 'bg-emerald-900 text-emerald-400' : 'bg-red-900 text-red-400'"
      >
        {{ connected ? 'Live' : 'Disconnected' }}
      </span>
    </div>

    <div v-if="!data" class="text-center text-gray-500 py-20">Connecting...</div>

    <template v-else>
      <SystemInfo :snapshot="data" />
      <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6 mt-6">
        <CpuGauge :cpu="data.cpu" />
        <MemoryBar :memory="data.memory" />
        <Temperature :temperatures="data.temperature" />
        <Voltage v-if="data.voltage && data.voltage.core > 0" :voltage="data.voltage" />
        <DiskUsage :disks="data.disks" />
        <NetworkTraffic :network="data.network" class="md:col-span-2 xl:col-span-2" />
      </div>
    </template>
  </div>
</template>
