<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { PinInfo } from '@/types/gpio'
import { gpioAvailable, getGpioPins, setGpioPin, exportGpioPin, unexportGpioPin } from '@/api/client'

const available = ref(false)
const pins = ref<PinInfo[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await gpioAvailable()
    available.value = res.available
    if (available.value) {
      pins.value = await getGpioPins()
    }
  } catch {
    available.value = false
  } finally {
    loading.value = false
  }
}

async function handleExport(pin: number) {
  await exportGpioPin(pin)
  await load()
}

async function handleUnexport(pin: number) {
  await unexportGpioPin(pin)
  await load()
}

async function handleSetOutput(pin: number, value: number) {
  await setGpioPin(pin, 'out', value)
  await load()
}

async function handleSetInput(pin: number) {
  await setGpioPin(pin, 'in', 0)
  await load()
}

onMounted(load)
</script>

<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-200">GPIO Control</h2>
      <button @click="load" :disabled="loading" class="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-sm text-gray-300 rounded-lg transition">
        Refresh
      </button>
    </div>

    <div v-if="!available" class="text-center text-gray-500 py-20">
      GPIO is only available on Raspberry Pi (Linux with /sys/class/gpio)
    </div>

    <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
      <div
        v-for="pin in pins"
        :key="pin.number"
        class="bg-gray-900 rounded-xl border border-gray-800 p-4 text-center"
      >
        <div class="text-xs text-gray-500 mb-1">GPIO</div>
        <div class="text-2xl font-bold text-gray-200 mb-2">{{ pin.number }}</div>

        <div class="text-xs mb-3">
          <span class="px-1.5 py-0.5 rounded" :class="{
            'bg-emerald-900/50 text-emerald-400': pin.mode === 'out',
            'bg-blue-900/50 text-blue-400': pin.mode === 'in',
            'bg-gray-800 text-gray-500': pin.mode === 'unexported'
          }">{{ pin.mode }}</span>
        </div>

        <div v-if="pin.mode === 'out'" class="space-y-2">
          <div class="text-lg font-bold" :class="pin.value ? 'text-emerald-400' : 'text-gray-500'">
            {{ pin.value ? 'HIGH' : 'LOW' }}
          </div>
          <div class="flex gap-1 justify-center">
            <button @click="handleSetOutput(pin.number, 1)" class="text-xs px-2 py-1 bg-emerald-900/30 text-emerald-400 rounded hover:bg-emerald-900/50">HIGH</button>
            <button @click="handleSetOutput(pin.number, 0)" class="text-xs px-2 py-1 bg-gray-800 text-gray-400 rounded hover:bg-gray-700">LOW</button>
          </div>
        </div>

        <div v-else-if="pin.mode === 'in'" class="space-y-2">
          <div class="text-lg font-bold" :class="pin.value ? 'text-emerald-400' : 'text-gray-500'">
            {{ pin.value ? 'HIGH' : 'LOW' }}
          </div>
        </div>

        <div class="mt-2 flex gap-1 justify-center">
          <template v-if="pin.mode === 'unexported'">
            <button @click="handleExport(pin.number)" class="text-xs px-2 py-1 bg-blue-900/30 text-blue-400 rounded hover:bg-blue-900/50">Export</button>
          </template>
          <template v-else>
            <button @click="handleSetInput(pin.number)" v-if="pin.mode !== 'in'" class="text-xs px-1.5 py-1 text-blue-400 hover:text-blue-300">IN</button>
            <button @click="handleSetOutput(pin.number, 0)" v-if="pin.mode !== 'out'" class="text-xs px-1.5 py-1 text-emerald-400 hover:text-emerald-300">OUT</button>
            <button @click="handleUnexport(pin.number)" class="text-xs px-1.5 py-1 text-red-400 hover:text-red-300">Unexport</button>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>
