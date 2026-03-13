<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { TemperatureInfo, HistoryEntry } from '@/types/system'
import Card from '@/components/common/Card.vue'

use([LineChart, GridComponent, TooltipComponent, CanvasRenderer])

const props = defineProps<{
  temperatures: TemperatureInfo[]
  history?: HistoryEntry[]
}>()

const emit = defineEmits<{
  clearHistory: []
}>()

interface TempPoint {
  time: string
  values: Record<string, number>
}

const historyData = ref<TempPoint[]>([])

watch(() => props.history, (entries) => {
  if (!entries) return
  historyData.value = entries.map(e => ({
    time: e.timestamp,
    values: Object.fromEntries((e.temperature || []).map(t => [t.sensor_key, t.temperature])),
  }))
}, { immediate: true })

watch(() => props.temperatures, (temps) => {
  if (!temps || temps.length === 0) return
  historyData.value.push({
    time: new Date().toISOString(),
    values: Object.fromEntries(temps.map(t => [t.sensor_key, t.temperature])),
  })
  if (historyData.value.length > 300) {
    historyData.value = historyData.value.slice(-300)
  }
})

const sensorColors = ['#10b981', '#f59e0b', '#3b82f6', '#ef4444']

const chartOption = computed(() => {
  const sensors = new Set<string>()
  for (const p of historyData.value) {
    for (const k of Object.keys(p.values)) sensors.add(k)
  }
  const sensorList = Array.from(sensors)

  return {
    animation: false,
    grid: { top: 10, right: 10, bottom: 20, left: 35 },
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const time = new Date(params[0].axisValue).toLocaleTimeString()
        const lines = params.map((p: any) => `${p.seriesName}: ${p.value[1].toFixed(1)}°C`)
        return `${time}<br/>${lines.join('<br/>')}`
      },
    },
    xAxis: { type: 'time', show: false },
    yAxis: { type: 'value', axisLabel: { color: '#9ca3af', fontSize: 10, formatter: '{value}°' } },
    series: sensorList.map((sensor, i) => ({
      name: sensor,
      type: 'line',
      showSymbol: false,
      lineStyle: { width: 1.5 },
      color: sensorColors[i % sensorColors.length],
      data: historyData.value
        .filter(p => p.values[sensor] !== undefined)
        .map(p => [p.time, p.values[sensor]]),
    })),
  }
})

function handleClear() {
  historyData.value = []
  emit('clearHistory')
}

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
    <div v-if="historyData.length > 1" class="mt-3">
      <VChart :option="chartOption" autoresize style="height: 120px" />
      <div class="flex justify-end mt-1">
        <button
          class="text-xs text-gray-500 hover:text-gray-300 transition-colors"
          @click="handleClear"
        >
          Clear History
        </button>
      </div>
    </div>
  </Card>
</template>
