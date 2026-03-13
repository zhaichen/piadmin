<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

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
</script>

<template>
  <div class="min-h-screen flex">
    <!-- Sidebar -->
    <aside class="w-56 bg-gray-900 border-r border-gray-800 flex flex-col shrink-0">
      <div class="p-4 border-b border-gray-800">
        <h1 class="text-xl font-bold text-emerald-400">PiAdmin</h1>
      </div>
      <nav class="flex-1 py-2">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="flex items-center gap-3 px-4 py-2.5 text-sm transition-colors"
          :class="isActive(item.path)
            ? 'bg-gray-800 text-emerald-400 border-r-2 border-emerald-400'
            : 'text-gray-400 hover:text-gray-200 hover:bg-gray-800/50'"
        >
          <span class="text-base w-5 text-center">{{ item.icon }}</span>
          {{ item.label }}
        </router-link>
      </nav>
      <div class="p-4 border-t border-gray-800">
        <button @click="logout" class="text-gray-500 hover:text-gray-300 text-sm transition w-full text-left">
          Logout
        </button>
      </div>
    </aside>

    <!-- Main content -->
    <main class="flex-1 overflow-auto">
      <router-view />
    </main>
  </div>
</template>
