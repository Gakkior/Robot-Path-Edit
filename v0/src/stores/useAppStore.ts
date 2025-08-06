/**
 * 应用状态管理
 * 
 * 使用Zustand - 现代化的状态管理库
 * 
 * 设计特点：
 * 1. 类型安全的状态管理
 * 2. 基于Immer的不可变更新 (直接使用 produce)
 * 3. 中间件支持（持久化、开发工具等）
 * 4. 模块化的状态切片
 * 
 * 参考项目：
 * - Zustand官方示例
 * - Linear应用的状态管理模式
 * - Figma的编辑器状态设计
 */

import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import { produce } from 'immer' // Directly import produce from immer
import type { 
  Node, 
  Path, 
  Viewport, 
  EditMode, 
  Command, 
  CommandType, 
  NodeID, 
  PathID 
} from '@/types'

// ============= 完整的Store类型 =============

interface AppState {
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

interface AppActions {
  setNodes: (nodes: Node[]) => void;
  addNode: (node: Node) => void;
  updateNode: (nodeId: NodeID, updates: Partial<Node>) => void;
  deleteNode: (nodeId: NodeID) => void;
  setPaths: (paths: Path[]) => void;
  addPath: (path: Path) => void;
  updatePath: (pathId: PathID, updates: Partial<Path>) => void;
  deletePath: (pathId: PathID) => void;
  setSelectedNodeId: (id: NodeID | null) => void;
  setSelectedPathId: (id: PathID | null) => void;
  setEditMode: (mode: EditMode) => void;
  setViewport: (viewport: Partial<Viewport>) => void;
  setIsDirty: (dirty: boolean) => void;
  clearCanvas: () => void;
  // Undo/Redo actions
  pushCommand: (command: Command) => void;
  undo: () => void;
  redo: () => void;
  // Path connection actions
  setIsConnecting: (isConnecting: boolean) => void;
  setConnectingNodeId: (nodeId: NodeID | null) => void;
}

// Command history limit
const HISTORY_LIMIT = 100;

export const useAppStore = create<AppState & AppActions>()(
  persist(
    (set, get) => ({ // Removed immer middleware from here
      nodes: [],
      paths: [],
      selectedNodeId: null,
      selectedPathId: null,
      editMode: 'select',
      viewport: { x: 0, y: 0, scale: 1 },
      isDirty: false,
      history: [],
      historyPointer: -1,
      isConnecting: false,
      connectingNodeId: null,

      setNodes: (nodes) => set(produce((state: AppState) => { state.nodes = nodes; state.isDirty = true; })),
      addNode: (node) => set(produce((state: AppState) => {
        state.nodes.push(node);
        state.isDirty = true;
        get().pushCommand({ type: 'addNode', payload: node });
      })),
      updateNode: (nodeId, updates) => set(produce((state: AppState) => {
        const nodeIndex = state.nodes.findIndex(n => n.id === nodeId);
        if (nodeIndex !== -1) {
          const oldNode = { ...state.nodes[nodeIndex] };
          state.nodes[nodeIndex] = { ...state.nodes[nodeIndex], ...updates };
          state.isDirty = true;
          // Only push command if position actually changed for undo/redo
          if (updates.position?.x !== undefined || updates.position?.y !== undefined) {
            get().pushCommand({ type: 'updateNodePosition', payload: { id: nodeId, oldX: oldNode.position.x, oldY: oldNode.position.y, newX: updates.position.x, newY: updates.position.y } });
          }
        }
      })),
      deleteNode: (nodeId) => set(produce((state: AppState) => {
        const nodeToDelete = state.nodes.find(n => n.id === nodeId);
        if (nodeToDelete) {
          state.nodes = state.nodes.filter(n => n.id !== nodeId);
          state.paths = state.paths.filter(p => p.from !== nodeId && p.to !== nodeId);
          state.isDirty = true;
          get().pushCommand({ type: 'deleteNode', payload: nodeToDelete });
        }
      })),
      setPaths: (paths) => set(produce((state: AppState) => { state.paths = paths; state.isDirty = true; })),
      addPath: (path) => set(produce((state: AppState) => {
        state.paths.push(path);
        state.isDirty = true;
        get().pushCommand({ type: 'addPath', payload: path });
      })),
      updatePath: (pathId, updates) => set(produce((state: AppState) => {
        const pathIndex = state.paths.findIndex(p => p.id === pathId);
        if (pathIndex !== -1) {
          state.paths[pathIndex] = { ...state.paths[pathIndex], ...updates };
          state.isDirty = true;
        }
      })),
      deletePath: (pathId) => set(produce((state: AppState) => {
        const pathToDelete = state.paths.find(p => p.id === pathId);
        if (pathToDelete) {
          state.paths = state.paths.filter(p => p.id !== pathId);
          state.isDirty = true;
          get().pushCommand({ type: 'deletePath', payload: pathToDelete });
        }
      })),
      setSelectedNodeId: (id) => set({ selectedNodeId: id, selectedPathId: null }),
      setSelectedPathId: (id) => set({ selectedPathId: id, selectedNodeId: null }),
      setEditMode: (mode) => set({ editMode: mode }),
      setViewport: (updates) => set(produce((state: AppState) => { state.viewport = { ...state.viewport, ...updates }; })),
      setIsDirty: (dirty) => set({ isDirty: dirty }),
      clearCanvas: () => set(produce((state: AppState) => {
        state.nodes = [];
        state.paths = [];
        state.selectedNodeId = null;
        state.selectedPathId = null;
        state.isDirty = true;
        state.history = [];
        state.historyPointer = -1;
      })),

      pushCommand: (command) => set(produce((state: AppState) => {
        // Clear redo history
        state.history = state.history.slice(0, state.historyPointer + 1);
        // Add new command
        state.history.push(command);
        // Trim history if it exceeds limit
        if (state.history.length > HISTORY_LIMIT) {
          state.history.shift();
        }
        state.historyPointer = state.history.length - 1;
      })),
      undo: () => set(produce((state: AppState) => {
        if (state.historyPointer >= 0) {
          const command = state.history[state.historyPointer];
          switch (command.type) {
            case 'addNode':
              state.nodes = state.nodes.filter(n => n.id !== command.payload.id);
              break;
            case 'deleteNode':
              state.nodes.push(command.payload);
              // Re-add paths connected to this node if they were deleted with it
              // This requires more complex payload to store deleted paths
              break;
            case 'updateNodePosition':
              const nodeToUpdatePos = state.nodes.find(n => n.id === command.payload.id);
              if (nodeToUpdatePos) {
                nodeToUpdatePos.position.x = command.payload.oldX;
                nodeToUpdatePos.position.y = command.payload.oldY;
              }
              break;
            case 'addPath':
              state.paths = state.paths.filter(p => p.id !== command.payload.id);
              break;
            case 'deletePath':
              state.paths.push(command.payload);
              break;
          }
          state.historyPointer--;
          state.isDirty = true; // Undo/Redo also makes it dirty
        }
      })),
      redo: () => set(produce((state: AppState) => {
        if (state.historyPointer < state.history.length - 1) {
          state.historyPointer++;
          const command = state.history[state.historyPointer];
          switch (command.type) {
            case 'addNode':
              state.nodes.push(command.payload);
              break;
            case 'deleteNode':
              state.nodes = state.nodes.filter(n => n.id !== command.payload.id);
              state.paths = state.paths.filter(p => p.from !== command.payload.id && p.to !== command.payload.id);
              break;
            case 'updateNodePosition':
              const nodeToUpdatePos = state.nodes.find(n => n.id === command.payload.id);
              if (nodeToUpdatePos) {
                nodeToUpdatePos.position.x = command.payload.newX;
                nodeToUpdatePos.position.y = command.payload.newY;
              }
              break;
            case 'addPath':
              state.paths.push(command.payload);
              break;
            case 'deletePath':
              state.paths = state.paths.filter(p => p.id !== command.payload.id);
              break;
          }
          state.isDirty = true; // Undo/Redo also makes it dirty
        }
      })),
      setIsConnecting: (isConnecting) => set({ isConnecting }),
      setConnectingNodeId: (nodeId) => set({ connectingNodeId: nodeId }),
    })
  ),
  {
    name: 'robot-path-editor-storage', // unique name
    storage: createJSONStorage(() => localStorage), // (optional) by default, 'localStorage' is used
    partialize: (state) => ({
      nodes: state.nodes,
      paths: state.paths,
      viewport: state.viewport,
    }),
    version: 1,
  }
);

// Custom hook to expose undo/redo capabilities and states
export const useCommandActions = () => {
  const { undo, redo, history, historyPointer } = useAppStore();
  const canUndo = historyPointer >= 0;
  const canRedo = historyPointer < history.length - 1;
  return { undo, redo, canUndo, canRedo };
};
