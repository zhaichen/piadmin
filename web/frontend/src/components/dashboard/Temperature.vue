<script setup lang="ts">
import type { TemperatureInfo } from '@/types/system'
import Card from '@/components/common/Card.vue'

const props = defineProps<{
  temperatures: TemperatureInfo[]
}>()

function tempColor(temp: number): string {
  if (temp > 80) return 'text-red-400'
  if (temp > 60) return 'text-amber-400'
  return 'text-emerald-400'
}
</script>

<template>
  <Card title="Temperature">
    <div class="space-y-3">
      <div v-for="t in temperatures" :key="t.sensor_key" class="flex justify-between items-center">
        <span class="text-gray-300 text-sm">{{ t.sensor_key }}</span>
        <span class="text-2xl font-bold" :class="tempColor(t.temperature)">
          {{ t.temperature.toFixed(1) }}&deg;C
        </span>
      </div>
      <div v-if="temperatures.length === 0" class="text-gray-500 text-sm">No sensor data</div>
    </div>
  </Card>
</template>
