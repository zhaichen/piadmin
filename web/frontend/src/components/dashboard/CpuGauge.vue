<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { GaugeChart } from 'echarts/charts'
import { CanvasRenderer } from 'echarts/renderers'
import type { CPUInfo } from '@/types/system'
import Card from '@/components/common/Card.vue'

use([GaugeChart, CanvasRenderer])

const props = defineProps<{
  cpu: CPUInfo
}>()

const option = computed(() => ({
  series: [
    {
      type: 'gauge',
      startAngle: 220,
      endAngle: -40,
      min: 0,
      max: 100,
      progress: { show: true, width: 14 },
      axisLine: { lineStyle: { width: 14, color: [[1, '#334155']] } },
      axisTick: { show: false },
      splitLine: { show: false },
      axisLabel: { show: false },
      pointer: { show: false },
      title: { show: false },
      detail: {
        valueAnimation: true,
        fontSize: 24,
        fontWeight: 'bold',
        color: '#e2e8f0',
        formatter: '{value}%',
        offsetCenter: [0, '0%'],
      },
      data: [{ value: Math.round(props.cpu.usage_total * 10) / 10 }],
      itemStyle: {
        color: props.cpu.usage_total > 80 ? '#ef4444' : props.cpu.usage_total > 50 ? '#f59e0b' : '#10b981',
      },
    },
  ],
}))
</script>

<template>
  <Card title="CPU">
    <VChart :option="option" autoresize style="height: 180px" />
    <div class="text-xs text-gray-500 mt-2 text-center">
      {{ cpu.model_name }} &middot; {{ cpu.cores }}C/{{ cpu.threads }}T
    </div>
  </Card>
</template>
