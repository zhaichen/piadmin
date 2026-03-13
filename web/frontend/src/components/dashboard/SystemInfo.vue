<script setup lang="ts">
import type { SystemSnapshot } from '@/types/system'

const props = defineProps<{
  snapshot: SystemSnapshot
}>()

function formatUptime(seconds: number): string {
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const mins = Math.floor((seconds % 3600) / 60)
  if (days > 0) return `${days}d ${hours}h ${mins}m`
  if (hours > 0) return `${hours}h ${mins}m`
  return `${mins}m`
}
</script>

<template>
  <div class="bg-gray-900 rounded-xl border border-gray-800 p-5 flex flex-wrap gap-x-8 gap-y-2 text-sm">
    <div><span class="text-gray-500">Host:</span> <span class="text-gray-200">{{ snapshot.hostname }}</span></div>
    <div><span class="text-gray-500">OS:</span> <span class="text-gray-200">{{ snapshot.platform }}</span></div>
    <div><span class="text-gray-500">Kernel:</span> <span class="text-gray-200">{{ snapshot.kernel_version }}</span></div>
    <div><span class="text-gray-500">Arch:</span> <span class="text-gray-200">{{ snapshot.arch }}</span></div>
    <div><span class="text-gray-500">Uptime:</span> <span class="text-gray-200">{{ formatUptime(snapshot.uptime) }}</span></div>
  </div>
</template>
