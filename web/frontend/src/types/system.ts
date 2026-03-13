export interface SystemSnapshot {
  timestamp: string
  hostname: string
  os: string
  platform: string
  arch: string
  kernel_version: string
  uptime: number
  cpu: CPUInfo
  memory: MemoryInfo
  disks: DiskInfo[]
  network: NetworkInfo[]
  temperature: TemperatureInfo[]
}

export interface CPUInfo {
  model_name: string
  cores: number
  threads: number
  usage_total: number
  usage_per: number[]
}

export interface MemoryInfo {
  total: number
  used: number
  available: number
  used_percent: number
  swap_total: number
  swap_used: number
}

export interface DiskInfo {
  device: string
  mount_point: string
  fs_type: string
  total: number
  used: number
  free: number
  used_percent: number
}

export interface NetworkInfo {
  name: string
  bytes_sent: number
  bytes_recv: number
  packs_sent: number
  packs_recv: number
}

export interface TemperatureInfo {
  sensor_key: string
  temperature: number
}

export interface WSMessage {
  type: string
  data: SystemSnapshot
}
