<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { FileEntry } from '@/types/file'
import { listFiles, downloadFile, deleteFile, mkdirFile, uploadFile } from '@/api/client'

const currentPath = ref('/')
const entries = ref<FileEntry[]>([])
const loading = ref(false)
const fileInput = ref<HTMLInputElement>()

async function load(path?: string) {
  if (path !== undefined) currentPath.value = path
  loading.value = true
  try {
    entries.value = await listFiles(currentPath.value)
  } catch (e: any) {
    entries.value = []
  } finally {
    loading.value = false
  }
}

function navigateTo(entry: FileEntry) {
  if (entry.is_dir) {
    load(entry.path)
  }
}

function goUp() {
  const parts = currentPath.value.split('/').filter(Boolean)
  parts.pop()
  load('/' + parts.join('/'))
}

async function handleDownload(entry: FileEntry) {
  const blob = await downloadFile(entry.path)
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = entry.name
  a.click()
  URL.revokeObjectURL(url)
}

async function handleDelete(entry: FileEntry) {
  if (!confirm(`Delete "${entry.name}"?`)) return
  await deleteFile(entry.path)
  await load()
}

async function handleMkdir() {
  const name = prompt('Directory name:')
  if (!name) return
  const path = currentPath.value === '/' ? `/${name}` : `${currentPath.value}/${name}`
  await mkdirFile(path)
  await load()
}

async function handleUpload() {
  fileInput.value?.click()
}

async function onFileSelected(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return
  for (const file of input.files) {
    await uploadFile(currentPath.value, file)
  }
  input.value = ''
  await load()
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return '-'
  const units = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}

function formatTime(ts: number): string {
  return new Date(ts * 1000).toLocaleString()
}

onMounted(() => load())
</script>

<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-200">File Manager</h2>
      <div class="flex gap-2">
        <button @click="handleMkdir" class="px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-sm text-gray-300 rounded-lg transition">New Folder</button>
        <button @click="handleUpload" class="px-3 py-1.5 bg-emerald-700 hover:bg-emerald-600 text-sm text-white rounded-lg transition">Upload</button>
        <input ref="fileInput" type="file" multiple class="hidden" @change="onFileSelected" />
      </div>
    </div>

    <!-- Breadcrumb -->
    <div class="flex items-center gap-1 mb-4 text-sm">
      <button @click="load('/')" class="text-emerald-400 hover:text-emerald-300">/</button>
      <template v-for="(part, i) in currentPath.split('/').filter(Boolean)" :key="i">
        <span class="text-gray-600">/</span>
        <button
          @click="load('/' + currentPath.split('/').filter(Boolean).slice(0, i + 1).join('/'))"
          class="text-gray-300 hover:text-emerald-400 transition"
        >{{ part }}</button>
      </template>
    </div>

    <div class="bg-gray-900 rounded-xl border border-gray-800 overflow-hidden">
      <table class="w-full text-sm">
        <thead>
          <tr class="text-gray-500 text-left border-b border-gray-800">
            <th class="px-4 py-3 font-medium">Name</th>
            <th class="px-4 py-3 font-medium text-right">Size</th>
            <th class="px-4 py-3 font-medium">Mode</th>
            <th class="px-4 py-3 font-medium">Modified</th>
            <th class="px-4 py-3 font-medium w-28"></th>
          </tr>
        </thead>
        <tbody>
          <!-- Go up -->
          <tr v-if="currentPath !== '/'" class="border-t border-gray-800/50 hover:bg-gray-800/30 cursor-pointer" @click="goUp">
            <td class="px-4 py-2 text-blue-400" colspan="5">..</td>
          </tr>
          <tr
            v-for="entry in entries"
            :key="entry.path"
            class="border-t border-gray-800/50 hover:bg-gray-800/30"
            :class="{ 'cursor-pointer': entry.is_dir }"
            @click="navigateTo(entry)"
          >
            <td class="px-4 py-2" :class="entry.is_dir ? 'text-blue-400' : 'text-gray-200'">
              <span class="mr-2">{{ entry.is_dir ? '📁' : '📄' }}</span>{{ entry.name }}
            </td>
            <td class="px-4 py-2 text-right text-gray-400">{{ entry.is_dir ? '-' : formatBytes(entry.size) }}</td>
            <td class="px-4 py-2 text-gray-500 font-mono text-xs">{{ entry.mode }}</td>
            <td class="px-4 py-2 text-gray-500 text-xs">{{ formatTime(entry.mod_time) }}</td>
            <td class="px-4 py-2 flex gap-2" @click.stop>
              <button v-if="!entry.is_dir" @click="handleDownload(entry)" class="text-xs text-emerald-400 hover:text-emerald-300">Download</button>
              <button @click="handleDelete(entry)" class="text-xs text-red-400 hover:text-red-300">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
