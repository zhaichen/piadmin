<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { systemPower } from '@/api/client'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const sidebarOpen = ref(false)
const powerLoading = ref(false)

const navItems = [
  { path: '/', label: 'Dashboard', icon: '⊞' },
  { path: '/processes', label: 'Processes', icon: '⚙' },
  { path: '/services', label: 'Services', icon: '☰' },
  { path: '/network', label: 'Network', icon: '⇄' },
  { path: '/terminal', label: 'Terminal', icon: '▸' },
  { path: '/files', label: 'Files', icon: '☷' },
  { path: '/gpio', label: 'GPIO', icon: '⊡' },
]

function logout() {
  authStore.logout()
  router.push('/login')
}

function isActive(path: string): boolean {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}

function closeSidebar() {
  sidebarOpen.value = false
}

function navigate(path: string) {
  router.push(path)
  closeSidebar()
}

async function handlePower(action: 'shutdown' | 'reboot') {
  const label = action === 'shutdown' ? '关机' : '重启'
  if (!confirm(`确定要${label}吗？`)) return
  powerLoading.value = true
  try {
    await systemPower(action)
    alert(`${label}指令已发送`)
  } catch (e: any) {
    alert(`${label}失败: ${e.message}`)
  } finally {
    powerLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex">
    <!-- Mobile overlay -->
    <div
      v-if="sidebarOpen"
      class="fixed inset-0 bg-black/60 z-40 md:hidden"
      @click="closeSidebar"
    ></div>

    <!-- Sidebar -->
    <aside
      class="fixed md:static inset-y-0 left-0 z-50 w-56 bg-gray-900 border-r border-gray-800 flex flex-col shrink-0 transition-transform duration-200 ease-in-out"
      :class="sidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'"
    >
      <div class="p-4 border-b border-gray-800 flex items-center justify-between">
        <h1 class="text-xl font-bold text-emerald-400">PiAdmin</h1>
        <button @click="closeSidebar" class="md:hidden text-gray-400 hover:text-gray-200 text-lg">✕</button>
      </div>
      <nav class="flex-1 py-2">
        <a
          v-for="item in navItems"
          :key="item.path"
          href="#"
          @click.prevent="navigate(item.path)"
          class="flex items-center gap-3 px-4 py-2.5 text-sm transition-colors"
          :class="isActive(item.path)
            ? 'bg-gray-800 text-emerald-400 border-r-2 border-emerald-400'
            : 'text-gray-400 hover:text-gray-200 hover:bg-gray-800/50'"
        >
          <span class="text-base w-5 text-center">{{ item.icon }}</span>
          {{ item.label }}
        </a>
      </nav>
      <div class="p-4 border-t border-gray-800 space-y-2">
        <div class="flex gap-2">
          <button
            @click="handlePower('reboot')"
            :disabled="powerLoading"
            class="flex-1 px-2 py-1.5 bg-amber-900/30 hover:bg-amber-900/50 text-amber-400 text-xs rounded-lg transition disabled:opacity-50"
          >
            ↻ 重启
          </button>
          <button
            @click="handlePower('shutdown')"
            :disabled="powerLoading"
            class="flex-1 px-2 py-1.5 bg-red-900/30 hover:bg-red-900/50 text-red-400 text-xs rounded-lg transition disabled:opacity-50"
          >
            ⏻ 关机
          </button>
        </div>
        <button @click="logout" class="text-gray-500 hover:text-gray-300 text-sm transition w-full text-left">
          Logout
        </button>
      </div>
    </aside>

    <!-- Main content -->
    <main class="flex-1 overflow-auto w-0 min-w-0">
      <!-- Mobile top bar -->
      <div class="md:hidden flex items-center gap-3 p-3 bg-gray-900 border-b border-gray-800 sticky top-0 z-30">
        <button @click="sidebarOpen = true" class="text-gray-300 hover:text-white text-xl">☰</button>
        <span class="text-emerald-400 font-bold">PiAdmin</span>
      </div>
      <router-view />
    </main>
  </div>
</template>
