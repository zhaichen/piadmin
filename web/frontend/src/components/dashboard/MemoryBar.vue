<script setup lang="ts">
import type { MemoryInfo } from '@/types/system'
import Card from '@/components/common/Card.vue'

const props = defineProps<{
  memory: MemoryInfo
}>()

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}

function barColor(percent: number): string {
  if (percent > 90) return 'bg-red-500'
  if (percent > 70) return 'bg-amber-500'
  return 'bg-emerald-500'
}
</script>

<template>
  <Card title="Memory">
    <div class="space-y-4">
      <!-- RAM -->
      <div>
        <div class="flex justify-between text-sm mb-1">
          <span class="text-gray-300">RAM</span>
          <span class="text-gray-400">{{ formatBytes(memory.used) }} / {{ formatBytes(memory.total) }}</span>
        </div>
        <div class="h-3 bg-gray-800 rounded-full overflow-hidden">
          <div
            class="h-full rounded-full transition-all duration-500"
            :class="barColor(memory.used_percent)"
            :style="{ width: memory.used_percent + '%' }"
          />
        </div>
        <div class="text-right text-xs text-gray-500 mt-1">{{ memory.used_percent.toFixed(1) }}%</div>
      </div>

      <!-- Swap -->
      <div v-if="memory.swap_total > 0">
        <div class="flex justify-between text-sm mb-1">
          <span class="text-gray-300">Swap</span>
          <span class="text-gray-400">{{ formatBytes(memory.swap_used) }} / {{ formatBytes(memory.swap_total) }}</span>
        </div>
        <div class="h-3 bg-gray-800 rounded-full overflow-hidden">
          <div
            class="h-full rounded-full transition-all duration-500 bg-blue-500"
            :style="{ width: (memory.swap_total > 0 ? (memory.swap_used / memory.swap_total * 100) : 0) + '%' }"
          />
        </div>
      </div>
    </div>
  </Card>
</template>
