<script setup lang="ts">
import type { DiskInfo } from '@/types/system'
import Card from '@/components/common/Card.vue'

const props = defineProps<{
  disks: DiskInfo[]
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
  <Card title="Disk">
    <div class="space-y-3">
      <div v-for="disk in disks" :key="disk.mount_point">
        <div class="flex justify-between text-sm mb-1">
          <span class="text-gray-300 truncate" :title="disk.device">{{ disk.mount_point }}</span>
          <span class="text-gray-400 shrink-0 ml-2">{{ formatBytes(disk.used) }} / {{ formatBytes(disk.total) }}</span>
        </div>
        <div class="h-2.5 bg-gray-800 rounded-full overflow-hidden">
          <div
            class="h-full rounded-full transition-all duration-500"
            :class="barColor(disk.used_percent)"
            :style="{ width: disk.used_percent + '%' }"
          />
        </div>
      </div>
      <div v-if="disks.length === 0" class="text-gray-500 text-sm">No disk data</div>
    </div>
  </Card>
</template>
