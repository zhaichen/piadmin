import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { SystemSnapshot } from '@/types/system'

export const useMonitorStore = defineStore('monitor', () => {
  const snapshot = ref<SystemSnapshot | null>(null)
  const connected = ref(false)

  function update(data: SystemSnapshot) {
    snapshot.value = data
  }

  return { snapshot, connected, update }
})
