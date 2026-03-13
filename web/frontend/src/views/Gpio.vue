<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { PinInfo } from '@/types/gpio'
import { gpioAvailable, getGpioPins, setGpioPin } from '@/api/client'

const available = ref(false)
const pins = ref<PinInfo[]>([])
const loading = ref(false)
const error = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    const res = await gpioAvailable()
    available.value = res.available
    if (available.value) {
      pins.value = await getGpioPins()
    }
  } catch (e: any) {
    error.value = e.message
    available.value = false
  } finally {
    loading.value = false
  }
}

async function handleSetHigh(pin: number) {
  try {
    await setGpioPin(pin, 'out', 1)
    await load()
  } catch (e: any) {
    error.value = e.message
  }
}

async function handleSetLow(pin: number) {
  try {
    await setGpioPin(pin, 'out', 0)
    await load()
  } catch (e: any) {
    error.value = e.message
  }
}

async function handleSetInput(pin: number) {
  try {
    await setGpioPin(pin, 'in', 0)
    await load()
  } catch (e: any) {
    error.value = e.message
  }
}

onMounted(load)
</script>

<template>
  <div class="p-4 md:p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-200">GPIO Control</h2>
      <button @click="load" :disabled="loading" class="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-sm text-gray-300 rounded-lg transition">
        Refresh
      </button>
    </div>

    <div v-if="error" class="mb-4 p-3 bg-red-900/30 border border-red-800 rounded-lg text-red-400 text-sm">
      {{ error }}
    </div>

    <div v-if="!available && !loading" class="text-center text-gray-500 py-20">
      GPIO is only available on Raspberry Pi with libgpiod installed
    </div>

    <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
      <div
        v-for="pin in pins"
        :key="pin.number"
        class="bg-gray-900 rounded-xl border border-gray-800 p-4 text-center"
      >
        <div class="text-xs text-gray-500 mb-1">GPIO</div>
        <div class="text-2xl font-bold text-gray-200 mb-2">{{ pin.number }}</div>

        <!-- Direction badge -->
        <div class="text-xs mb-3">
          <span class="px-1.5 py-0.5 rounded" :class="{
            'bg-emerald-900/50 text-emerald-400': pin.mode === 'output',
            'bg-blue-900/50 text-blue-400': pin.mode === 'input',
          }">{{ pin.mode }}</span>
        </div>

        <!-- Value display -->
        <div class="text-lg font-bold mb-3" :class="pin.value ? 'text-emerald-400' : 'text-gray-500'">
          {{ pin.value ? 'HIGH' : 'LOW' }}
        </div>

        <!-- Controls -->
        <div class="flex flex-col gap-1.5">
          <div class="flex gap-1 justify-center">
            <button @click="handleSetHigh(pin.number)" class="text-xs px-2 py-1 bg-emerald-900/30 text-emerald-400 rounded hover:bg-emerald-900/50 transition">HIGH</button>
            <button @click="handleSetLow(pin.number)" class="text-xs px-2 py-1 bg-gray-800 text-gray-400 rounded hover:bg-gray-700 transition">LOW</button>
          </div>
          <button
            v-if="pin.mode !== 'input'"
            @click="handleSetInput(pin.number)"
            class="text-xs px-2 py-1 text-blue-400 hover:text-blue-300 transition"
          >Set Input</button>
        </div>
      </div>
    </div>
  </div>
</template>
