/**
* 主画布组件 - Konva.js版本
* 
* 使用Konva.js来实现画布交互功能
* 提供节点和路径的可视化渲染、拖拽、缩放、平移等交互
* 选择和编辑功能
* 响应式设计
* 
* 设计参考：
* - Figma的画布交互模式
* - Draw.io的编辑体验
* - Miro的协作功能
*/

'use client'

import React, { useRef, useEffect, useState, useCallback } from 'react'
import { Stage, Layer, Circle, Line, Text, Group } from 'react-konva'
import { v4 as uuidv4 } from 'uuid'
import { useAppStore } from '@/stores/useAppStore'
import {
useCreateNode,
useUpdateNode,
useDeleteNode,
useCreatePath,
useDeletePath,
useNodes,
usePaths,
} from '@/services' // Import from the new barrel file
import type { Node, Path, Position, NodeID } from '@/types'
import { distance, getLineCenter } from '@/utils/canvas' // Corrected import

interface CanvasProps {
width?: number;
height?: number;
}

export const Canvas: React.FC<CanvasProps> = ({ width = 800, height = 600 }) => {
const stageRef = useRef<any>(null)
const {
  nodes,
  paths,
  selectedNodeId,
  selectedPathId,
  editMode,
  viewport,
  isConnecting,
  connectingNodeId,
  setNodes,
  setPaths,
  setSelectedNodeId,
  setSelectedPathId,
  setViewport,
  setIsConnecting,
  setConnectingNodeId,
  deleteNode, // Use the action from store
  deletePath, // Use the action from store
  addNode, // Use the action from store
  addPath, // Use the action from store
} = useAppStore()

const createNodeMutation = useCreateNode()
const updateNodeMutation = useUpdateNode()
const createPathMutation = useCreatePath()
const deleteNodeMutation = useDeleteNode()
const deletePathMutation = useDeletePath()

const { data: fetchedNodes } = useNodes()
const { data: fetchedPaths } = usePaths()

useEffect(() => {
  if (fetchedNodes) {
    setNodes(fetchedNodes)
  }
}, [fetchedNodes, setNodes])

useEffect(() => {
  if (fetchedPaths) {
    setPaths(fetchedPaths)
  }
}, [fetchedPaths, setPaths])

const [dragStart, setDragStart] = useState<{ x: number; y: number } | null>(null)
const DRAG_THRESHOLD = 5 // Pixels to differentiate click from drag

const handleStageMouseDown = useCallback((e: any) => {
  const stage = stageRef.current
  if (!stage) return

  const pointerPos = stage.getPointerPosition()
  if (!pointerPos) return

  setDragStart({ x: pointerPos.x, y: pointerPos.y })

  // If click on empty space, clear selection
  if (e.target === stage) {
    setSelectedNodeId(null)
    setSelectedPathId(null)
    if (editMode === 'add-node') {
      const newId = uuidv4()
      const newNode: Node = {
        id: newId,
        name: `Node ${nodes.length + 1}`,
        type: 'normal',
        position: {
          x: (pointerPos.x - viewport.x) / viewport.scale,
          y: (pointerPos.y - viewport.y) / viewport.scale,
        },
        status: 'active',
        metadata: {},
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
      addNode(newNode) // Add to Zustand store
      createNodeMutation.mutate(newNode) // Persist to backend
    } else if (editMode === 'add-path') {
      setIsConnecting(false)
      setConnectingNodeId(null)
    }
  }
}, [editMode, nodes.length, viewport, setSelectedNodeId, setIsConnecting, setConnectingNodeId, addNode, createNodeMutation])

const handleStageMouseMove = useCallback((e: any) => {
  const stage = stageRef.current
  if (!stage || !dragStart) return

  const pointerPos = stage.getPointerPosition()
  if (!pointerPos) return

  const dx = pointerPos.x - dragStart.x
  const dy = pointerPos.y - dragStart.y

  // Only pan if not in add-node/add-path/delete mode and drag threshold is met
  if (editMode === 'select' && (Math.abs(dx) > DRAG_THRESHOLD || Math.abs(dy) > DRAG_THRESHOLD)) {
    setViewport({
      x: viewport.x + dx,
      y: viewport.y + dy,
    })
    setDragStart({ x: pointerPos.x, y: pointerPos.y }) // Update drag start for continuous panning
  }
}, [dragStart, editMode, viewport, setViewport])

const handleStageMouseUp = useCallback(() => {
  setDragStart(null)
}, [])

const handleWheel = useCallback((e: any) => {
  e.evt.preventDefault()
  const stage = stageRef.current
  if (!stage) return

  const scaleBy = 1.1
  const oldScale = viewport.scale
  const pointer = stage.getPointerPosition()

  if (!pointer) return

  const mousePointTo = {
    x: (pointer.x - viewport.x) / oldScale,
    y: (pointer.y - viewport.y) / oldScale,
  }

  const newScale = e.evt.deltaY < 0 ? oldScale * scaleBy : oldScale / scaleBy

  setViewport({
    scale: newScale,
    x: pointer.x - mousePointTo.x * newScale,
    y: pointer.y - mousePointTo.y * newScale,
  })
}, [viewport, setViewport])

const handleNodeDragEnd = useCallback((e: any) => {
  const nodeId = e.target.id()
  const newPosition = { x: e.target.x(), y: e.target.y() }
  updateNodeMutation.mutate({ id: nodeId, position: newPosition })
  // Update local state immediately for responsiveness
  useAppStore.getState().updateNode(nodeId, { position: newPosition })
}, [updateNodeMutation])

const handleNodeClick = useCallback((nodeId: NodeID) => {
  setSelectedNodeId(nodeId)
  setSelectedPathId(null) // Clear path selection

  if (editMode === 'delete') {
    deleteNode(nodeId) // Delete from Zustand store
    deleteNodeMutation.mutate(nodeId) // Persist to backend
  } else if (editMode === 'add-path') {
    if (!isConnecting) {
      setIsConnecting(true)
      setConnectingNodeId(nodeId)
    } else if (connectingNodeId && connectingNodeId !== nodeId) {
      // Create path between connectingNodeId and current nodeId
      const newPathId = uuidv4()
      const newPath: Path = {
        id: newPathId,
        name: `Path ${paths.length + 1}`,
        type: 'normal',
        from: connectingNodeId,
        to: nodeId,
        weight: 1,
        status: 'active',
        metadata: {},
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
      addPath(newPath) // Add to Zustand store
      createPathMutation.mutate(newPath) // Persist to backend
      setIsConnecting(false)
      setConnectingNodeId(null)
    } else if (connectingNodeId === nodeId) {
      // Clicked the same node again, cancel connection
      setIsConnecting(false)
      setConnectingNodeId(null)
    }
  }
}, [editMode, isConnecting, connectingNodeId, paths.length, setSelectedNodeId, setSelectedPathId, setIsConnecting, setConnectingNodeId, deleteNode, deleteNodeMutation, addPath, createPathMutation])

const handlePathClick = useCallback((pathId: any) => {
  setSelectedPathId(pathId)
  setSelectedNodeId(null) // Clear node selection

  if (editMode === 'delete') {
    deletePath(pathId) // Delete from Zustand store
    deletePathMutation.mutate(pathId) // Persist to backend
  }
}, [editMode, setSelectedPathId, setSelectedNodeId, deletePath, deletePathMutation])

useEffect(() => {
  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === 'Delete' || e.key === 'Backspace') {
      if (selectedNodeId) {
        deleteNode(selectedNodeId)
        deleteNodeMutation.mutate(selectedNodeId)
        setSelectedNodeId(null)
      } else if (selectedPathId) {
        deletePath(selectedPathId)
        deletePathMutation.mutate(selectedPathId)
        setSelectedPathId(null)
      }
    }
  }

  window.addEventListener('keydown', handleKeyDown)
  return () => window.removeEventListener('keydown', handleKeyDown)
}, [selectedNodeId, selectedPathId, deleteNode, deleteNodeMutation, deletePath, deletePathMutation, setSelectedNodeId, setSelectedPathId])

return (
  <Stage
    width={width}
    height={height}
    onMouseDown={handleStageMouseDown}
    onMouseMove={handleStageMouseMove}
    onMouseUp={handleStageMouseUp}
    onWheel={handleWheel}
    scaleX={viewport.scale}
    scaleY={viewport.scale}
    x={viewport.x}
    y={viewport.y}
    ref={stageRef}
    className="bg-gray-100"
  >
    <Layer>
      {/* Render Paths */}
      {paths.map((path) => {
        const fromNode = nodes.find((n) => n.id === path.from)
        const toNode = nodes.find((n) => n.id === path.to)

        if (!fromNode || !toNode) return null

        return (
          <Line
            key={path.id}
            points={[fromNode.position.x, fromNode.position.y, toNode.position.x, toNode.position.y]}
            stroke={selectedPathId === path.id ? 'blue' : 'gray'}
            strokeWidth={path.weight || 2}
            tension={0.5}
            lineCap="round"
            lineJoin="round"
            onClick={() => handlePathClick(path.id)}
            onTap={() => handlePathClick(path.id)}
            hitStrokeWidth={10} // Increase hit area for easier selection
          />
        )
      })}

      {/* Render Nodes */}
      {nodes.map((node) => (
        <Group
          key={node.id}
          id={node.id}
          x={node.position.x}
          y={node.position.y}
          draggable={editMode === 'select'}
          onDragEnd={handleNodeDragEnd}
          onClick={() => handleNodeClick(node.id)}
          onTap={() => handleNodeClick(node.id)}
        >
          <Circle
            radius={20}
            fill={selectedNodeId === node.id ? 'red' : 'green'}
            stroke="black"
            strokeWidth={1}
          />
          <Text
            text={node.name}
            fontSize={12}
            fill="white"
            align="center"
            verticalAlign="middle"
            offsetX={node.name.length * 3} // Adjust text position
            offsetY={6}
          />
        </Group>
      ))}

      {/* Render connecting line if in add-path mode and a node is selected */}
      {editMode === 'add-path' && isConnecting && connectingNodeId && (
        <Line
          points={[
            nodes.find(n => n.id === connectingNodeId)?.position.x || 0,
            nodes.find(n => n.id === connectingNodeId)?.position.y || 0,
            (stageRef.current?.getPointerPosition()?.x - viewport.x) / viewport.scale || 0,
            (stageRef.current?.getPointerPosition()?.y - viewport.y) / viewport.scale || 0,
          ]}
          stroke="purple"
          strokeWidth={2}
          dash={[10, 5]}
        />
      )}
    </Layer>
  </Stage>
)
}
