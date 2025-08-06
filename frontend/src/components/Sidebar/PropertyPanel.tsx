/**
 * 属性面板组件
 * 
 * 显示和编辑选中元素的属性
 * 支持节点和路径的详细配置
 */

import React, { useState, useEffect } from 'react' // Import useState and useEffect
import { motion, AnimatePresence } from 'framer-motion'
import { useAppStore, useDataActions } from '@/stores/useAppStore' // Import useDataActions
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/Select'
import { Separator } from '@/components/ui/Separator'
import { Badge } from '@/components/ui/Badge'
import type { Node, Path, NodeType, PathType } from '@/types' // Import specific types

export const PropertyPanel: React.FC = () => {
  const { nodes, paths, selectedElements } = useAppStore()
  
  const selectedNodeIds = Array.from(selectedElements.selectedNodes)
  const selectedPathIds = Array.from(selectedElements.selectedPaths)
  
  const selectedNode = selectedNodeIds.length === 1 ? nodes[selectedNodeIds[0]] : null
  const selectedPath = selectedPathIds.length === 1 ? paths[selectedPathIds[0]] : null
  
  // If no elements are selected
  if (selectedNodeIds.length === 0 && selectedPathIds.length === 0) {
    return (
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="p-6 text-center"
      >
        <div className="text-gray-400 text-6xl mb-4">👈</div>
        <h3 className="text-lg font-medium text-gray-900 mb-2">选择元素</h3>
        <p className="text-gray-500 text-sm">
          点击画布上的节点或路径来查看和编辑属性
        </p>
      </motion.div>
    )
  }
  
  // Multiple selection state
  if (selectedNodeIds.length > 1 || selectedPathIds.length > 1) {
    return (
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="p-6"
      >
        <h3 className="text-lg font-medium text-gray-900 mb-4">批量编辑</h3>
        <div className="space-y-4">
          {selectedNodeIds.length > 0 && (
            <div>
              <Badge variant="secondary" className="mb-2">
                {selectedNodeIds.length} 个节点
              </Badge>
              <div className="space-y-2">
                <Button variant="outline" size="sm" className="w-full">
                  批量设置类型
                </Button>
                <Button variant="outline" size="sm" className="w-full">
                  批量设置状态
                </Button>
              </div>
            </div>
          )}
          
          {selectedPathIds.length > 0 && (
            <div>
              <Badge variant="secondary" className="mb-2">
                {selectedPathIds.length} 条路径
              </Badge>
              <div className="space-y-2">
                <Button variant="outline" size="sm" className="w-full">
                  批量设置权重
                </Button>
                <Button variant="outline" size="sm" className="w-full">
                  批量设置类型
                </Button>
              </div>
            </div>
          )}
          
          <Separator />
          
          <div className="space-y-2">
            <Button variant="destructive" size="sm" className="w-full">
              删除选中项
            </Button>
          </div>
        </div>
      </motion.div>
    )
  }
  
  return (
    <div className="h-full overflow-y-auto">
      <AnimatePresence mode="wait">
        {selectedNode && (
          <NodePropertyPanel key={selectedNode.id} node={selectedNode} />
        )}
        {selectedPath && (
          <PathPropertyPanel key={selectedPath.id} path={selectedPath} />
        )}
      </AnimatePresence>
    </div>
  )
}

// 节点属性面板
const NodePropertyPanel: React.FC<{ node: Node }> = ({ node }) => {
  const { setNode, setDirty } = useDataActions()

  // Local state for editable fields
  const [name, setName] = useState(node.name)
  const [type, setType] = useState<NodeType>(node.type)
  const [status, setStatus] = useState(node.status)
  const [posX, setPosX] = useState(node.position.x)
  const [posY, setPosY] = useState(node.position.y)
  const [posZ, setPosZ] = useState(node.position.z || 0)

  // Update local state when node prop changes (e.g., new node selected)
  useEffect(() => {
    setName(node.name)
    setType(node.type)
    setStatus(node.status)
    setPosX(node.position.x)
    setPosY(node.position.y)
    setPosZ(node.position.z || 0)
  }, [node])

  const handleSave = () => {
    setNode(node.id, {
      ...node,
      name,
      type,
      status,
      position: {
        x: parseFloat(posX.toFixed(2)), // Ensure two decimal places
        y: parseFloat(posY.toFixed(2)),
        z: parseFloat(posZ.toFixed(2)),
      },
      updated_at: new Date().toISOString(),
    })
    setDirty(true) // Mark as dirty after saving
    console.log('Node saved:', { id: node.id, name, type, status, position: { x: posX, y: posY, z: posZ } })
  }

  const handleDelete = () => {
    // TODO: Implement actual delete logic
    console.log(`Delete node: ${node.id}`)
  }

  return (
    <motion.div
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: -20 }}
      className="p-6 space-y-6"
    >
      <div>
        <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center">
          📍 节点属性
          <Badge variant="outline" className="ml-2">
            {type}
          </Badge>
        </h3>
      </div>
      
      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            节点ID
          </label>
          <Input value={node.id} readOnly /> {/* Added readOnly */}
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            节点名称
          </label>
          <Input 
            value={name} 
            onChange={(e) => setName(e.target.value)} 
            placeholder="输入节点名称" 
          />
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            节点类型
          </label>
          <Select value={type} onValueChange={(value: NodeType) => setType(value)}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="normal">普通节点</SelectItem>
              <SelectItem value="start">起始节点</SelectItem>
              <SelectItem value="end">结束节点</SelectItem>
              <SelectItem value="waypoint">路径点</SelectItem>
              <SelectItem value="charging">充电站</SelectItem>
              <SelectItem value="storage">存储点</SelectItem>
            </SelectContent>
          </Select>
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            状态
          </label>
          <Select value={status} onValueChange={(value) => setStatus(value as 'active' | 'inactive' | 'error')}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="active">激活</SelectItem>
              <SelectItem value="inactive">未激活</SelectItem>
              <SelectItem value="error">错误</SelectItem>
            </SelectContent>
          </Select>
        </div>
        
        <Separator />
        
        <div>
          <h4 className="text-sm font-medium text-gray-700 mb-3">坐标位置</h4>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-xs text-gray-500 mb-1">X坐标</label>
              <Input
                type="number"
                value={posX}
                onChange={(e) => setPosX(parseFloat(e.target.value))}
                step="0.1"
              />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">Y坐标</label>
              <Input
                type="number"
                value={posY}
                onChange={(e) => setPosY(parseFloat(e.target.value))}
                step="0.1"
              />
            </div>
          </div>
          <div className="mt-3">
            <label className="block text-xs text-gray-500 mb-1">Z坐标</label>
            <Input
              type="number"
              value={posZ}
              onChange={(e) => setPosZ(parseFloat(e.target.value))}
              step="0.1"
            />
          </div>
        </div>
        
        <Separator />
        
        <div>
          <h4 className="text-sm font-medium text-gray-700 mb-3">时间信息</h4>
          <div className="space-y-2 text-xs text-gray-500">
            <div>创建时间: {new Date(node.created_at).toLocaleString()}</div>
            <div>更新时间: {new Date(node.updated_at).toLocaleString()}</div>
          </div>
        </div>
        
        <div className="flex space-x-2 pt-4">
          <Button className="flex-1" onClick={handleSave}>
            保存更改
          </Button>
          <Button variant="destructive" size="icon" onClick={handleDelete}>
            🗑️
          </Button>
        </div>
      </div>
    </motion.div>
  )
}

// 路径属性面板
const PathPropertyPanel: React.FC<{ path: Path }> = ({ path }) => {
  const { nodes } = useAppStore()
  const { setPath, setDirty } = useDataActions()

  // Local state for editable fields
  const [name, setName] = useState(path.name)
  const [type, setType] = useState<PathType>(path.type)
  const [status, setStatus] = useState(path.status)
  const [weight, setWeight] = useState(path.weight)

  // Update local state when path prop changes
  useEffect(() => {
    setName(path.name)
    setType(path.type)
    setStatus(path.status)
    setWeight(path.weight)
  }, [path])

  const handleSave = () => {
    setPath(path.id, {
      ...path,
      name,
      type,
      status,
      weight: parseFloat(weight.toFixed(2)), // Ensure two decimal places
      updated_at: new Date().toISOString(),
    })
    setDirty(true) // Mark as dirty after saving
    console.log('Path saved:', { id: path.id, name, type, status, weight })
  }

  const handleDelete = () => {
    // TODO: Implement actual delete logic
    console.log(`Delete path: ${path.id}`)
  }
  
  const startNode = nodes.find(n => n.id === path.start_node_id)
  const endNode = nodes.find(n => n.id === path.end_node_id)
  
  return (
    <motion.div
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: -20 }}
      className="p-6 space-y-6"
    >
      <div>
        <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center">
          🔗 路径属性
          <Badge variant="outline" className="ml-2">
            {type}
          </Badge>
        </h3>
      </div>
      
      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            路径ID
          </label>
          <Input value={path.id} readOnly /> {/* Added readOnly */}
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            路径名称
          </label>
          <Input 
            value={name} 
            onChange={(e) => setName(e.target.value)} 
            placeholder="输入路径名称" 
          />
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            路径类型
          </label>
          <Select value={type} onValueChange={(value: PathType) => setType(value)}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="normal">普通路径</SelectItem>
              <SelectItem value="bidirectional">双向路径</SelectItem>
              <SelectItem value="one-way">单向路径</SelectItem>
              <SelectItem value="emergency">紧急路径</SelectItem>
            </SelectContent>
          </Select>
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            状态
          </label>
          <Select value={status} onValueChange={(value) => setStatus(value as 'active' | 'inactive' | 'blocked')}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="active">激活</SelectItem>
              <SelectItem value="inactive">未激活</SelectItem>
              <SelectItem value="blocked">阻塞</SelectItem>
            </SelectContent>
          </Select>
        </div>
        
        <Separator />
        
        <div>
          <h4 className="text-sm font-medium text-gray-700 mb-3">连接节点</h4>
          <div className="space-y-3">
            <div className="p-3 bg-gray-50 rounded-lg">
              <div className="text-xs text-gray-500 mb-1">起始节点</div>
              <div className="font-medium">{startNode?.name || '未知节点'}</div>
              <div className="text-xs text-gray-500">ID: {path.start_node_id}</div>
            </div>
            <div className="text-center text-gray-400">↓</div>
            <div className="p-3 bg-gray-50 rounded-lg">
              <div className="text-xs text-gray-500 mb-1">结束节点</div>
              <div className="font-medium">{endNode?.name || '未知节点'}</div>
              <div className="text-xs text-gray-500">ID: {path.end_node_id}</div>
            </div>
          </div>
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            路径权重
          </label>
          <Input
            type="number"
            value={weight}
            onChange={(e) => setWeight(parseFloat(e.target.value))}
            step="0.1"
            min="0"
          />
        </div>
        
        <Separator />
        
        <div>
          <h4 className="text-sm font-medium text-gray-700 mb-3">时间信息</h4>
          <div className="space-y-2 text-xs text-gray-500">
            <div>创建时间: {new Date(path.created_at).toLocaleString()}</div>
            <div>更新时间: {new Date(path.updated_at).toLocaleString()}</div>
          </div>
        </div>
        
        <div className="flex space-x-2 pt-4">
          <Button className="flex-1" onClick={handleSave}>
            保存更改
          </Button>
          <Button variant="destructive" size="icon" onClick={handleDelete}>
            🗑️
          </Button>
        </div>
      </div>
    </motion.div>
  )
}
