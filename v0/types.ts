/**
 * 核心类型定义
 * 
 * 设计理念：
 * 1. 强类型约束，减少运行时错误
 * 2. 领域驱动设计，类型即文档
 * 3. 可扩展性，便于后续功能添加
 */

// ============= 基础类型 =============

export interface Position {
  x: number
  y: number
  z?: number
}

export interface Metadata {
  [key: string]: any
}

// ============= 节点相关类型 =============

export type NodeType = 'normal' | 'start' | 'end' | 'waypoint' | 'charging' | 'storage'

export interface Node {
  id: string
  name: string
  type: NodeType
  position: Position
  status: 'active' | 'inactive' | 'error'
  metadata: Metadata
  created_at: string
  updated_at: string
}

export interface CreateNodeRequest {
  name: string
  type?: NodeType
  position: Position
  metadata?: Metadata
}

export interface UpdateNodeRequest {
  id: string
  name?: string
  type?: NodeType
  position?: Position
  status?: 'active' | 'inactive' | 'error'
  metadata?: Metadata
}

// ============= 路径相关类型 =============

export type PathType = 'normal' | 'bidirectional' | 'one-way' | 'emergency'

export interface Path {
  id: string
  name: string
  type: PathType
  start_node_id: string
  end_node_id: string
  weight: number
  status: 'active' | 'inactive' | 'blocked'
  metadata: Metadata
  created_at: string
  updated_at: string
}

export interface CreatePathRequest {
  name: string
  type?: PathType
  start_node_id: string
  end_node_id: string
  weight?: number
  metadata?: Metadata
}

export interface UpdatePathRequest {
  id: string
  name?: string
  type?: PathType
  start_node_id?: string
  end_node_id?: string
  weight?: number
  status?: 'active' | 'inactive' | 'blocked'
  metadata?: Metadata
}

// ============= 画布相关类型 =============

export interface CanvasData {
  nodes: Record<string, Node>
  paths: Record<string, Path>
}

export interface ViewportState {
  x: number
  y: number
  scale: number
}

export interface SelectionState {
  selectedNodes: Set<string>
  selectedPaths: Set<string>
  hoveredNode: string | null
  hoveredPath: string | null
}

// ============= 布局算法类型 =============

export type LayoutAlgorithm = 
  | 'force-directed'
  | 'hierarchical' 
  | 'circular'
  | 'grid'
  | 'tree'

// ============= 路径生成算法类型 =============

export type PathGenerationAlgorithm = 
  | 'nearest-neighbor'
  | 'full-connectivity'
  | 'grid'
  | 'mst' // 最小生成树
  | 'shortest-path'

// ============= API响应类型 =============

export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  error?: string
  message?: string
}

// ============= 应用状态类型 =============

export interface AppState {
  // 数据状态
  nodes: Record<string, Node>
  paths: Record<string, Path>
  
  // UI状态
  currentView: 'canvas' | 'table'
  sidebarOpen: boolean
  selectedElements: SelectionState
  viewport: ViewportState
  
  // 编辑状态
  editMode: 'select' | 'add-node' | 'add-path' | 'delete'
  isDirty: boolean
  
  // 系统状态
  loading: boolean
  error: string | null
  connectionStatus: 'connected' | 'disconnected' | 'connecting'
}
