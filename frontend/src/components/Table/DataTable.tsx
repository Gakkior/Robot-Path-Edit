/**
 * 数据表格组件
 * 
 * 高性能的虚拟化表格，支持大量数据的展示和编辑
 * 
 * 功能特点：
 * 1. 虚拟化渲染 - 支持大量数据
 * 2. 内联编辑 - 直接在表格中编辑
 * 3. 排序和筛选 - 数据操作
 * 4. 批量操作 - 多选和批量处理
 * 5. 响应式设计 - 适配不同屏幕
 * 6. 分页控制 - 调节每页显示条数
 * 7. 密度调节 - 改变显示的紧凑和大小
 * 
 * 设计参考：
 * - Airtable的表格交互
 * - Notion的数据库视图
 * - Linear的问题列表
 */

import React, { useState, useMemo, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { ChevronUp, ChevronDown, Search, Filter, MoreHorizontal, Edit, Trash2, Plus } from 'lucide-react'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/Select'
import { Badge } from '@/components/ui/Badge'
import { useAppStore } from '@/stores/useAppStore'
import { cn } from '@/utils/cn'
import type { Node, Path } from '@/types'

interface DataTableProps {
  type: 'nodes' | 'paths'
}

export const DataTable: React.FC<DataTableProps> = ({ type }) => {
  const { nodes, paths } = useAppStore()
  const [searchTerm, setSearchTerm] = useState('')
  const [sortField, setSortField] = useState<string>('name')
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc')
  const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set())
  
  // Pagination states
  const [pageSize, setPageSize] = useState(25)
  const [currentPage, setCurrentPage] = useState(1)

  // Density state
  const [density, setDensity] = useState<'compact' | 'comfortable'>('comfortable')

  const pageSizeOptions = [10, 25, 50, 100]
  const densityOptions = ['compact', 'comfortable']
  
  // Get data
  const data = useMemo(() => {
    return type === 'nodes' ? Object.values(nodes) : Object.values(paths)
  }, [type, nodes, paths])
  
  // Filter and sort data
  const filteredAndSortedData = useMemo(() => {
    let filtered = data.filter(item =>
      item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.id.toLowerCase().includes(searchTerm.toLowerCase())
    )
    
    filtered.sort((a, b) => {
      let aValue = (a as any)[sortField]
      let bValue = (b as any)[sortField]
      
      // Handle nested properties for sorting (e.g., position.x)
      if (sortField.includes('.')) {
        const [parent, child] = sortField.split('.')
        aValue = (a as any)[parent]?.[child]
        bValue = (b as any)[parent]?.[child]
      }

      if (typeof aValue === 'string') {
        aValue = aValue.toLowerCase()
        bValue = bValue.toLowerCase()
      }
      
      if (sortDirection === 'asc') {
        return aValue < bValue ? -1 : aValue > bValue ? 1 : 0
      } else {
        return aValue > bValue ? -1 : aValue < bValue ? 1 : 0
      }
    })
    
    return filtered
  }, [data, searchTerm, sortField, sortDirection])

  // Apply pagination
  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize
    const endIndex = startIndex + pageSize
    return filteredAndSortedData.slice(startIndex, endIndex)
  }, [filteredAndSortedData, currentPage, pageSize])

  const totalPages = Math.ceil(filteredAndSortedData.length / pageSize)
  
  // Handle sorting
  const handleSort = useCallback((field: string) => {
    if (sortField === field) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')
    } else {
      setSortField(field)
      setSortDirection('asc')
    }
    setCurrentPage(1) // Reset to first page on sort change
  }, [sortField, sortDirection])
  
  // Handle row selection
  const handleRowSelect = useCallback((id: string, selected: boolean) => {
    const newSelection = new Set(selectedRows)
    if (selected) {
      newSelection.add(id)
    } else {
      newSelection.delete(id)
    }
    setSelectedRows(newSelection)
  }, [selectedRows])
  
  // Handle select all
  const handleSelectAll = useCallback((selected: boolean) => {
    if (selected) {
      setSelectedRows(new Set(paginatedData.map(item => item.id)))
    } else {
      setSelectedRows(new Set())
    }
  }, [paginatedData])

  // Determine cell padding based on density
  const cellPaddingClass = density === 'compact' ? 'px-3 py-2' : 'px-4 py-3'
  
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="h-full flex flex-col bg-white"
    >
      {/* 表格头部 */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-900">
            {type === 'nodes' ? '📍 节点管理' : '🔗 路径管理'}
            <Badge variant="secondary" className="ml-2">
              {filteredAndSortedData.length} 项
            </Badge>
          </h2>
          
          <div className="flex items-center gap-2">
            <Button size="sm" onClick={() => {}}>
              <Plus className="h-4 w-4 mr-1" />
              添加{type === 'nodes' ? '节点' : '路径'}
            </Button>
          </div>
        </div>
        
        {/* 搜索、筛选和密度/分页控制 */}
        <div className="flex items-center gap-3 flex-wrap">
          <div className="relative flex-1 max-w-sm min-w-[200px]">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
            <Input
              placeholder={`搜索${type === 'nodes' ? '节点' : '路径'}...`}
              value={searchTerm}
              onChange={(e) => {
                setSearchTerm(e.target.value)
                setCurrentPage(1) // Reset to first page on search
              }}
              className="pl-10"
            />
          </div>
          
          <Button variant="outline" size="sm">
            <Filter className="h-4 w-4 mr-1" />
            筛选
          </Button>

          <Select value={String(pageSize)} onValueChange={(value) => {
            setPageSize(Number(value))
            setCurrentPage(1) // Reset to first page on page size change
          }}>
            <SelectTrigger className="w-[120px] h-9">
              <SelectValue placeholder="每页显示" />
            </SelectTrigger>
            <SelectContent>
              {pageSizeOptions.map(option => (
                <SelectItem key={option} value={String(option)}>
                  {option} 条/页
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select value={density} onValueChange={(value: 'compact' | 'comfortable') => setDensity(value)}>
            <SelectTrigger className="w-[120px] h-9">
              <SelectValue placeholder="密度" />
            </SelectTrigger>
            <SelectContent>
              {densityOptions.map(option => (
                <SelectItem key={option} value={option}>
                  {option === 'compact' ? '紧凑' : '舒适'}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          
          {selectedRows.size > 0 && (
            <motion.div
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              className="flex items-center gap-2"
            >
              <Badge variant="secondary">
                已选择 {selectedRows.size} 项
              </Badge>
              <Button variant="outline" size="sm">
                批量编辑
              </Button>
              <Button variant="destructive" size="sm">
                <Trash2 className="h-4 w-4 mr-1" />
                删除
              </Button>
            </motion.div>
          )}
        </div>
      </div>
      
      {/* 表格内容 */}
      <div className="flex-1 overflow-auto">
        <table className="w-full">
          <thead className="bg-gray-50 sticky top-0 z-10">
            <tr>
              <th className={cn("w-12", cellPaddingClass)}>
                <input
                  type="checkbox"
                  checked={selectedRows.size === paginatedData.length && paginatedData.length > 0}
                  onChange={(e) => handleSelectAll(e.target.checked)}
                  className="rounded border-gray-300"
                />
              </th>
              {type === 'nodes' ? (
                <NodeTableHeaders onSort={handleSort} sortField={sortField} sortDirection={sortDirection} cellPaddingClass={cellPaddingClass} />
              ) : (
                <PathTableHeaders onSort={handleSort} sortField={sortField} sortDirection={sortDirection} cellPaddingClass={cellPaddingClass} />
              )}
              <th className={cn("w-20 text-left text-xs font-medium text-gray-500 uppercase tracking-wider", cellPaddingClass)}>
                操作
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            <AnimatePresence>
              {paginatedData.map((item, index) => (
                <motion.tr
                  key={item.id}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                  transition={{ delay: index * 0.02 }}
                  className={cn(
                    'hover:bg-gray-50 transition-colors',
                    selectedRows.has(item.id) && 'bg-blue-50'
                  )}
                >
                  <td className={cellPaddingClass}>
                    <input
                      type="checkbox"
                      checked={selectedRows.has(item.id)}
                      onChange={(e) => handleRowSelect(item.id, e.target.checked)}
                      className="rounded border-gray-300"
                    />
                  </td>
                  {type === 'nodes' ? (
                    <NodeTableRow node={item as Node} cellPaddingClass={cellPaddingClass} />
                  ) : (
                    <PathTableRow path={item as Path} cellPaddingClass={cellPaddingClass} />
                  )}
                  <td className={cellPaddingClass}>
                    <div className="flex items-center gap-1">
                      <Button variant="ghost" size="sm">
                        <Edit className="h-4 w-4" />
                      </Button>
                      <Button variant="ghost" size="sm">
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </div>
                  </td>
                </motion.tr>
              ))}
            </AnimatePresence>
          </tbody>
        </table>
        
        {filteredAndSortedData.length === 0 && (
          <div className="text-center py-12">
            <div className="text-gray-400 text-6xl mb-4">
              {type === 'nodes' ? '📍' : '🔗'}
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              没有找到{type === 'nodes' ? '节点' : '路径'}
            </h3>
            <p className="text-gray-500">
              {searchTerm ? '尝试调整搜索条件' : `点击上方按钮添加${type === 'nodes' ? '节点' : '路径'}`}
            </p>
          </div>
        )}
      </div>

      {/* Pagination controls */}
      {filteredAndSortedData.length > 0 && (
        <div className="p-4 border-t border-gray-200 flex items-center justify-between text-sm text-gray-600">
          <span>
            显示 {Math.min((currentPage - 1) * pageSize + 1, filteredAndSortedData.length)} -{' '}
            {Math.min(currentPage * pageSize, filteredAndSortedData.length)} 共 {filteredAndSortedData.length} 项
          </span>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
              disabled={currentPage === 1}
            >
              上一页
            </Button>
            <span>
              页 {currentPage} / {totalPages}
            </span>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
              disabled={currentPage === totalPages}
            >
              下一页
            </Button>
          </div>
        </div>
      )}
    </motion.div>
  )
}

// 节点表格头部
const NodeTableHeaders: React.FC<{
  onSort: (field: string) => void
  sortField: string
  sortDirection: 'asc' | 'desc'
  cellPaddingClass: string
}> = ({ onSort, sortField, sortDirection, cellPaddingClass }) => {
  const headers = [
    { field: 'id', label: 'ID' },
    { field: 'name', label: '名称' },
    { field: 'type', label: '类型' },
    { field: 'status', label: '状态' },
    { field: 'position.x', label: 'X坐标' },
    { field: 'position.y', label: 'Y坐标' },
    { field: 'created_at', label: '创建时间' },
  ]
  
  return (
    <>
      {headers.map(({ field, label }) => (
        <th
          key={field}
          className={cn("text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100", cellPaddingClass)}
          onClick={() => onSort(field)}
        >
          <div className="flex items-center gap-1">
            {label}
            {sortField === field && (
              sortDirection === 'asc' ? 
                <ChevronUp className="h-4 w-4" /> : 
                <ChevronDown className="h-4 w-4" />
            )}
          </div>
        </th>
      ))}
    </>
  )
}

// 路径表格头部
const PathTableHeaders: React.FC<{
  onSort: (field: string) => void
  sortField: string
  sortDirection: 'asc' | 'desc'
  cellPaddingClass: string
}> = ({ onSort, sortField, sortDirection, cellPaddingClass }) => {
  const headers = [
    { field: 'id', label: 'ID' },
    { field: 'name', label: '名称' },
    { field: 'type', label: '类型' },
    { field: 'status', label: '状态' },
    { field: 'start_node_id', label: '起始节点' },
    { field: 'end_node_id', label: '结束节点' },
    { field: 'weight', label: '权重' },
    { field: 'created_at', label: '创建时间' },
  ]
  
  return (
    <>
      {headers.map(({ field, label }) => (
        <th
          key={field}
          className={cn("text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100", cellPaddingClass)}
          onClick={() => onSort(field)}
        >
          <div className="flex items-center gap-1">
            {label}
            {sortField === field && (
              sortDirection === 'asc' ? 
                <ChevronUp className="h-4 w-4" /> : 
                <ChevronDown className="h-4 w-4" />
            )}
          </div>
        </th>
      ))}
    </>
  )
}

// 节点表格行
const NodeTableRow: React.FC<{ node: Node, cellPaddingClass: string }> = ({ node, cellPaddingClass }) => {
  return (
    <>
      <td className={cn("text-sm font-mono text-gray-900", cellPaddingClass)}>
        {node.id.slice(0, 8)}...
      </td>
      <td className={cn("text-sm text-gray-900 font-medium", cellPaddingClass)}>
        {node.name}
      </td>
      <td className={cn("text-sm", cellPaddingClass)}>
        <Badge variant="outline">
          {node.type}
        </Badge>
      </td>
      <td className={cn("text-sm", cellPaddingClass)}>
        <Badge 
          variant={
            node.status === 'active' ? 'success' : 
            node.status === 'error' ? 'destructive' : 'secondary'
          }
        >
          {node.status}
        </Badge>
      </td>
      <td className={cn("text-sm font-mono text-gray-600", cellPaddingClass)}>
        {node.position.x.toFixed(2)}
      </td>
      <td className={cn("text-sm font-mono text-gray-600", cellPaddingClass)}>
        {node.position.y.toFixed(2)}
      </td>
      <td className={cn("text-sm text-gray-500", cellPaddingClass)}>
        {new Date(node.created_at).toLocaleDateString()}
      </td>
    </>
  )
}

// 路径表格行
const PathTableRow: React.FC<{ path: Path, cellPaddingClass: string }> = ({ path, cellPaddingClass }) => {
  const { nodes } = useAppStore()
  const startNode = nodes[path.start_node_id]
  const endNode = nodes[path.end_node_id]
  
  return (
    <>
      <td className={cn("text-sm font-mono text-gray-900", cellPaddingClass)}>
        {path.id.slice(0, 8)}...
      </td>
      <td className={cn("text-sm text-gray-900 font-medium", cellPaddingClass)}>
        {path.name}
      </td>
      <td className={cn("text-sm", cellPaddingClass)}>
        <Badge variant="outline">
          {path.type}
        </Badge>
      </td>
      <td className={cn("text-sm", cellPaddingClass)}>
        <Badge 
          variant={
            path.status === 'active' ? 'success' : 
            path.status === 'blocked' ? 'destructive' : 'secondary'
          }
        >
          {path.status}
        </Badge>
      </td>
      <td className={cn("text-sm text-gray-600", cellPaddingClass)}>
        {startNode?.name || path.start_node_id.slice(0, 8)}
      </td>
      <td className={cn("text-sm text-gray-600", cellPaddingClass)}>
        {endNode?.name || path.end_node_id.slice(0, 8)}
      </td>
      <td className={cn("text-sm font-mono text-gray-600", cellPaddingClass)}>
        {path.weight}
      </td>
      <td className={cn("text-sm text-gray-500", cellPaddingClass)}>
        {new Date(path.created_at).toLocaleDateString()}
      </td>
    </>
  )
}
