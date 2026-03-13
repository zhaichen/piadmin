<script setup lang="ts">
import type { NetworkInfo } from '@/types/system'
import Card from '@/components/common/Card.vue'

const props = defineProps<{
  network: NetworkInfo[]
}>()

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}
</script>

<template>
  <Card title="Network">
    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="text-gray-500 text-left">
            <th class="pb-2 font-medium">Interface</th>
            <th class="pb-2 font-medium text-right">Sent</th>
            <th class="pb-2 font-medium text-right">Received</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="n in network" :key="n.name" class="border-t border-gray-800">
            <td class="py-2 text-gray-300">{{ n.name }}</td>
            <td class="py-2 text-right text-blue-400">{{ formatBytes(n.bytes_sent) }}</td>
            <td class="py-2 text-right text-emerald-400">{{ formatBytes(n.bytes_recv) }}</td>
          </tr>
        </tbody>
      </table>
      <div v-if="network.length === 0" class="text-gray-500 text-sm">No network data</div>
    </div>
  </Card>
</template>
