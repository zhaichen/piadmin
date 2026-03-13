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
  voltage: VoltageInfo
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

export interface VoltageInfo {
  core: number
  sdram_c: number
  sdram_i: number
  sdram_p: number
  throttle: ThrottleInfo
}

export interface ThrottleInfo {
  raw: number
  under_voltage: boolean
  freq_capped: boolean
  throttled: boolean
  soft_temp_limit: boolean
  under_voltage_occurred: boolean
  freq_capped_occurred: boolean
  throttled_occurred: boolean
  soft_temp_limit_occurred: boolean
}

export interface HistoryEntry {
  timestamp: string
  temperature: TemperatureInfo[]
  voltage: VoltageInfo
}

export interface WSMessage {
  type: string
  data: SystemSnapshot
}
