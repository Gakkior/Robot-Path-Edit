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

export type NodeID = string;
export type NodeType = 'normal' | 'start' | 'end' | 'waypoint' | 'charging' | 'storage' | 'robot' | 'charging_station';

export interface Node {
id: NodeID;
name: string;
type: NodeType;
position: Position;
status: 'active' | 'inactive' | 'error';
properties?: Record<string, any>;
metadata: Metadata;
created_at: string;
updated_at: string;
}

export interface CreateNodeRequest {
name: string;
type?: NodeType;
position: Position;
properties?: Record<string, any>;
metadata?: Metadata;
}

export interface UpdateNodeRequest {
id: NodeID;
name?: string;
type?: NodeType;
position?: Position;
status?: 'active' | 'inactive' | 'error';
properties?: Record<string, any>;
metadata?: Metadata;
}

// ============= 路径相关类型 =============

export type PathID = string;
export type PathType = 'normal' | 'bidirectional' | 'one-way' | 'emergency' | 'direct' | 'curved';

export interface Path {
id: PathID;
from: NodeID;
to: NodeID;
type: PathType;
weight: number;
status: 'active' | 'inactive' | 'blocked';
properties?: Record<string, any>;
metadata: Metadata;
created_at: string;
updated_at: string;
}

export interface CreatePathRequest {
from: NodeID;
to: NodeID;
name?: string;
type?: PathType;
weight?: number;
properties?: Record<string, any>;
metadata?: Metadata;
}

export interface UpdatePathRequest {
id: PathID;
from?: NodeID;
to?: NodeID;
name?: string;
type?: PathType;
weight?: number;
status?: 'active' | 'inactive' | 'blocked';
properties?: Record<string, any>;
metadata?: Metadata;
}

// ============= 画布相关类型 =============

export interface CanvasData {
nodes: Record<NodeID, Node>
paths: Record<PathID, Path>
}

export interface Viewport {
x: number;
y: number;
scale: number;
}

export interface SelectionState {
selectedNodes: Set<NodeID>
selectedPaths: Set<PathID>
hoveredNode: NodeID | null
hoveredPath: PathID | null
}

// ============= 布局算法类型 =============

export type LayoutAlgorithm = 
| 'force-directed'
| 'hierarchical' 
| 'circular'
| 'grid'
| 'tree'
| 'pipeline'
| 'radial'
| 'custom'

export interface LayoutConfig {
algorithm: LayoutAlgorithm
parameters: Record<string, any>
}

// ============= 路径生成算法类型 =============

export type PathGenerationAlgorithm = 
| 'nearest-neighbor'
| 'full-connectivity'
| 'grid'
| 'mst' // 最小生成树
| 'shortest-path'

export interface PathGenerationConfig {
algorithm: PathGenerationAlgorithm
parameters: Record<string, any>
}

// ============= 命令模式类型 (Undo/Redo) =============

// 定义可撤销/重做的操作类型
export type CommandType = 
| 'addNode' 
| 'deleteNode' 
| 'updateNodePosition'
| 'addPath'
| 'deletePath'
| 'updateElementProperties' // For generic property updates

// 历史记录条目
export interface Command {
type: CommandType;
payload: any; // Specific data for the command (e.g., node, path, old/new position)
}

// 命令历史状态
export interface CommandHistoryState {
history: Command[];
historyPointer: number; // Points to the current state in history
}

// ============= API响应类型 =============

export interface ApiResponse<T = any> {
success: boolean
data?: T
error?: string
message?: string
}

export interface ListResponse<T> {
items: T[];
total: number;
page: number;
pageSize: number;
}

export interface HealthStatus {
status: string; // e.g., "ok", "error"
message: string;
timestamp: string;
}

export interface SystemStats {
uptime: string;
memoryUsage: string;
cpuUsage: string;
goroutines: number;
}

// ============= 应用状态类型 =============

export interface AppState {
// 数据状态
nodes: Node[];
paths: Path[];
selectedNodeId: NodeID | null;
selectedPathId: PathID | null;
editMode: EditMode;
viewport: Viewport;
isDirty: boolean; // Indicates if there are unsaved changes
history: Command[];
historyPointer: number;
isConnecting: boolean; // State for path creation
connectingNodeId: NodeID | null; // The first node selected for path creation
}

// ============= 事件类型 =============

export interface CanvasEvent {
type: 'node-click' | 'node-drag' | 'path-click' | 'canvas-click'
target?: Node | Path
position?: Position
modifiers?: {
  shift: boolean
  ctrl: boolean
  alt: boolean
}
}

// ============= 插件系统类型 =============

export interface Plugin {
id: string
name: string
version: string
description: string
author: string

// 插件生命周期
install(app: any): void
uninstall(app: any): void

// 插件能力
layouts?: LayoutAlgorithm[]
pathGenerators?: PathGenerationAlgorithm[]
components?: React.ComponentType[]
}

// ============= 配置类型 =============

export interface AppConfig {
// 画布配置
canvas: {
  width: number
  height: number
  backgroundColor: string
  gridSize: number
  showGrid: boolean
}

// 节点配置
node: {
  defaultRadius: number
  defaultColor: string
  selectedColor: string
  hoveredColor: string
}

// 路径配置
path: {
  defaultWidth: number
  defaultColor: string
  selectedColor: string
  hoveredColor: string
}

// 性能配置
performance: {
  maxNodes: number
  maxPaths: number
  renderThrottle: number
}

// 快捷键配置
shortcuts: Record<string, string>
}

// ============= 数据库连接类型 =============

export type DatabaseType = 'mysql' | 'sqlite' | 'postgresql' | 'sqlserver' | 'oracle' | 'mongodb' | 'cassandra' | 'redis' | 'other';

export interface DatabaseConnection {
id: string;
name: string;
type: DatabaseType;
host: string;
port: number;
username?: string;
password?: string;
database?: string;
url?: string; // For JDBC-style URLs or custom connection strings
properties?: Record<string, string>; // Additional connection properties (e.g., SSL, SSH)
createdAt: string;
updatedAt: string;
}

export interface CreateConnectionRequest {
name: string;
type: DatabaseType;
host: string;
port: number;
database?: string;
username?: string;
password?: string;
url?: string;
properties?: Record<string, string>;
}

export interface UpdateConnectionRequest {
id: string;
name?: string;
type?: DatabaseType;
host?: string;
port?: number;
database?: string;
username?: string;
password?: string;
url?: string;
properties?: Record<string, string>;
}

// ============= 表映射类型 =============

export interface TableMapping {
id: string;
connectionId: string;
nodeTableName: string;
nodeIdField: string;
nodeXField: string;
nodeYField: string;
nodeNameField: string;
pathTableName: string;
pathIdField: string;
pathFromField: string;
pathToField: string;
createdAt: string;
updatedAt: string;
}

export interface CreateTableMappingRequest {
connectionId: string;
nodeTableName: string;
nodeIdField: string;
nodeXField: string;
nodeYField: string;
nodeNameField: string;
pathTableName: string;
pathIdField: string;
pathFromField: string;
pathToField: string;
}

export interface UpdateTableMappingRequest {
id: string;
connectionId?: string;
nodeTableName?: string;
nodeIdField?: string;
nodeXField?: string;
nodeYField?: string;
nodeNameField?: string;
pathTableName?: string;
pathIdField?: string;
pathFromField?: string;
pathToField?: string;
}

export type TableMappingType = 'node' | 'path';

// ============= 数据同步和验证类型 =============

export interface ValidateTableResponse {
is_valid: boolean
message: string
columns?: { name: string; type: string }[]
missing_fields?: string[]
}

export interface SyncResult {
nodes_synced: number
paths_synced: number
message: string
}

// ============= 模板类型 =============

export type LayoutType = 'tree' | 'grid' | 'circular' | 'force-directed' | 'pipeline' | 'hierarchical' | 'radial' | 'custom';

export interface Template {
id: string;
name: string;
description?: string;
category?: string; // e.g., "factory", "warehouse", "lab"
layoutType: LayoutType;
nodes: Node[];
paths: Path[];
isPublic: boolean;
createdAt: string;
updatedAt: string;
}

export interface CreateTemplateRequest {
name: string;
description?: string;
category?: string;
layoutType: LayoutType;
nodes: Node[];
paths: Path[];
isPublic?: boolean;
}

export interface UpdateTemplateRequest {
id: string;
name?: string;
description?: string;
category?: string;
layoutType?: LayoutType;
nodes?: Node[];
paths?: Path[];
isPublic?: boolean;
}

export interface ApplyTemplateRequest {
width?: number;
height?: number;
}

export interface ApplyTemplateResponse {
nodes: Node[];
paths: Path[];
}

export interface ExportDataRequest {
type: 'nodes' | 'paths' | 'all';
format: 'csv' | 'xlsx';
}

export interface ExportTemplateResponse {
template_data: string; // Base64 encoded or JSON string
}

export interface TemplateStats {
total_templates: number;
public_templates: number;
categories: string[];
popular_templates: Array<{ id: string; name: string; usage_count: number }>;
}
