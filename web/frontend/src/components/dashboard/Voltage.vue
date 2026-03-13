<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { VoltageInfo, HistoryEntry } from '@/types/system'
import Card from '@/components/common/Card.vue'

use([LineChart, GridComponent, TooltipComponent, CanvasRenderer])

const props = defineProps<{
  voltage: VoltageInfo
  history?: HistoryEntry[]
}>()

const emit = defineEmits<{
  clearHistory: []
}>()

interface VoltPoint {
  time: string
  core: number
}

const historyData = ref<VoltPoint[]>([])

watch(() => props.history, (entries) => {
  if (!entries) return
  historyData.value = entries
    .filter(e => e.voltage && e.voltage.core > 0)
    .map(e => ({ time: e.timestamp, core: e.voltage.core }))
}, { immediate: true })

watch(() => props.voltage, (v) => {
  if (!v || v.core <= 0) return
  historyData.value.push({ time: new Date().toISOString(), core: v.core })
  if (historyData.value.length > 300) {
    historyData.value = historyData.value.slice(-300)
  }
})

const chartOption = computed(() => ({
  animation: false,
  grid: { top: 10, right: 10, bottom: 20, left: 40 },
  tooltip: {
    trigger: 'axis',
    formatter: (params: any) => {
      const time = new Date(params[0].axisValue).toLocaleTimeString()
      return `${time}<br/>Core: ${params[0].value[1].toFixed(4)}V`
    },
  },
  xAxis: { type: 'time', show: false },
  yAxis: {
    type: 'value',
    min: (v: any) => Math.floor(v.min * 100) / 100,
    max: (v: any) => Math.ceil(v.max * 100) / 100,
    axisLabel: { color: '#9ca3af', fontSize: 10, formatter: '{value}V' },
  },
  series: [{
    name: 'Core',
    type: 'line',
    showSymbol: false,
    lineStyle: { width: 1.5 },
    color: '#8b5cf6',
    areaStyle: { color: 'rgba(139, 92, 246, 0.1)' },
    data: historyData.value.map(p => [p.time, p.core]),
  }],
}))

function handleClear() {
  historyData.value = []
  emit('clearHistory')
}

interface ThrottleFlag {
  label: string
  current: boolean
  occurred: boolean
}

const throttleFlags = computed<ThrottleFlag[]>(() => {
  const t = props.voltage.throttle
  if (!t) return []
  return [
    { label: 'Under-voltage', current: t.under_voltage, occurred: t.under_voltage_occurred },
    { label: 'Freq Capped', current: t.freq_capped, occurred: t.freq_capped_occurred },
    { label: 'Throttled', current: t.throttled, occurred: t.throttled_occurred },
    { label: 'Soft Temp Limit', current: t.soft_temp_limit, occurred: t.soft_temp_limit_occurred },
  ]
})

function flagClass(f: ThrottleFlag): string {
  if (f.current) return 'bg-red-900/50 text-red-400 border-red-800'
  if (f.occurred) return 'bg-amber-900/50 text-amber-400 border-amber-800'
  return 'bg-emerald-900/50 text-emerald-400 border-emerald-800'
}
</script>

<template>
  <Card title="Voltage">
    <div class="grid grid-cols-2 gap-3">
      <div class="text-center">
        <div class="text-xs text-gray-500">Core</div>
        <div class="text-lg font-bold text-purple-400">{{ voltage.core.toFixed(4) }}V</div>
      </div>
      <div class="text-center">
        <div class="text-xs text-gray-500">SDRAM_C</div>
        <div class="text-lg font-bold text-purple-400">{{ voltage.sdram_c.toFixed(4) }}V</div>
      </div>
      <div class="text-center">
        <div class="text-xs text-gray-500">SDRAM_I</div>
        <div class="text-lg font-bold text-purple-400">{{ voltage.sdram_i.toFixed(4) }}V</div>
      </div>
      <div class="text-center">
        <div class="text-xs text-gray-500">SDRAM_P</div>
        <div class="text-lg font-bold text-purple-400">{{ voltage.sdram_p.toFixed(4) }}V</div>
      </div>
    </div>

    <div v-if="throttleFlags.length > 0" class="mt-3 flex flex-wrap gap-1.5">
      <span
        v-for="f in throttleFlags"
        :key="f.label"
        class="text-xs px-2 py-0.5 rounded border"
        :class="flagClass(f)"
      >
        {{ f.label }}
      </span>
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
