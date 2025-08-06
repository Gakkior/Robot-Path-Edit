// src/components/Toolbar/Toolbar.tsx

/**
 * 工具栏组件
 * 
 * 提供画布编辑的各种工具和操作
 * 包括编辑模式切换、布局算法、路径生成等功能
 * 
 * 设计参考：
 * - Figma的工具栏设计
 * - Adobe Creative Suite的工具面板
 * - VS Code的活动栏
 */

import React from 'react'
import { motion } from 'framer-motion'
import { MousePointer, Plus, Trash2, GitBranch, Grid3X3, Circle, Zap, RotateCcw, RotateCw, ZoomIn, ZoomOut, Maximize, Save, Download, Upload, LayoutTemplate } from 'lucide-react'
import { Button } from '@/components/ui/Button'
import { Separator } from '@/components/ui/Separator'
import { Badge } from '@/components/ui/Badge'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/DropdownMenu'
import { useAppStore, useCommandActions } from '@/stores/useAppStore'
import { useApplyLayout, useGeneratePaths, useExportData } from '@/services/api'
import { cn } from '@/utils/cn'

interface ToolbarProps {
  onOpenTemplateModal: () => void;
}

export const Toolbar: React.FC<ToolbarProps> = ({ onOpenTemplateModal }) => {
  const { editMode, setEditMode, viewport, setViewport, isDirty } = useAppStore()
  const { undo, redo, canUndo, canRedo } = useCommandActions();
  const applyLayoutMutation = useApplyLayout()
  const generatePathsMutation = useGeneratePaths()
  const exportDataMutation = useExportData();
  
  // 编辑模式按钮配置
  const editModes = [
    { mode: 'select' as const, icon: MousePointer, label: '选择', shortcut: 'V' },
    { mode: 'add-node' as const, icon: Plus, label: '添加节点', shortcut: 'N' },
    { mode: 'add-path' as const, icon: GitBranch, label: '添加路径', shortcut: 'P' },
    { mode: 'delete' as const, icon: Trash2, label: '删除', shortcut: 'D' },
  ]
  
  // 布局算法按钮配置
  const layoutAlgorithms = [
    { algorithm: 'force-directed' as const, icon: Zap, label: '力导向布局' },
    { algorithm: 'grid' as const, icon: Grid3X3, label: '网格布局' },
    { algorithm: 'circular' as const, icon: Circle, label: '圆形布局' },
  ]
  
  // 处理编辑模式切换
  const handleEditModeChange = (mode: typeof editMode) => {
    setEditMode(mode)
  }
  
  // 处理布局算法应用
  const handleApplyLayout = (algorithm: string) => {
    applyLayoutMutation.mutate({ algorithm: algorithm as any })
  }
  
  // 处理路径生成
  const handleGeneratePaths = (algorithm: string) => {
    generatePathsMutation.mutate({ algorithm: algorithm as any })
  }
  
  // 处理缩放
  const handleZoomIn = () => {
    setViewport({ scale: Math.min(viewport.scale * 1.2, 5) })
  }
  
  const handleZoomOut = () => {
    setViewport({ scale: Math.max(viewport.scale / 1.2, 0.1) })
  }
  
  const handleZoomFit = () => {
    setViewport({ scale: 1, x: 0, y: 0 })
  }

  // Handle export
  const handleExport = (type: 'nodes' | 'paths' | 'all', format: 'csv' | 'xlsx') => {
    exportDataMutation.mutate({ type, format });
  };
  
  return (
    <motion.div
      initial={{ y: -20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      className="bg-white border-b border-gray-200 px-4 py-2 flex items-center gap-2 shadow-sm"
    >
      {/* 编辑模式组 */}
      <div className="flex items-center gap-1 bg-gray-50 rounded-lg p-1">
        {editModes.map(({ mode, icon: Icon, label, shortcut }) => (
          <Button
            key={mode}
            variant={editMode === mode ? 'default' : 'ghost'}
            size="sm"
            onClick={() => handleEditModeChange(mode)}
            className={cn(
              'relative',
              editMode === mode && 'shadow-sm'
            )}
            title={`${label} (${shortcut})`}
          >
            <Icon className="h-4 w-4" />
            <span className="ml-1 hidden sm:inline">{label}</span>
            {editMode === mode && (
              <motion.div
                layoutId="activeMode"
                className="absolute inset-0 bg-primary-600 rounded-md -z-10"
                transition={{ type: 'spring', bounce: 0.2, duration: 0.6 }}
              />
            )}
          </Button>
        ))}
      </div>
      
      <Separator orientation="vertical" className="h-6" />

      {/* Undo/Redo Group */}
      <div className="flex items-center gap-1">
        <Button
          variant="outline"
          size="sm"
          onClick={undo}
          disabled={!canUndo}
          title="撤销 (Ctrl+Z)"
        >
          <RotateCcw className="h-4 w-4" />
          <span className="ml-1 hidden sm:inline">撤销</span>
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={redo}
          disabled={!canRedo}
          title="重做 (Ctrl+Y)"
        >
          <RotateCw className="h-4 w-4" />
          <span className="ml-1 hidden sm:inline">重做</span>
        </Button>
      </div>
      
      <Separator orientation="vertical" className="h-6" />
      
      {/* 布局算法组 */}
      <div className="flex items-center gap-1">
        <span className="text-sm text-gray-600 mr-2">布局:</span>
        {layoutAlgorithms.map(({ algorithm, icon: Icon, label }) => (
          <Button
            key={algorithm}
            variant="outline"
            size="sm"
            onClick={() => handleApplyLayout(algorithm)}
            disabled={applyLayoutMutation.isPending}
            title={label}
          >
            <Icon className="h-4 w-4" />
            <span className="ml-1 hidden md:inline">{label}</span>
          </Button>
        ))}
      </div>
      
      <Separator orientation="vertical" className="h-6" />
      
      {/* 路径生成组 */}
      <div className="flex items-center gap-1">
        <span className="text-sm text-gray-600 mr-2">路径:</span>
        <Button
          variant="outline"
          size="sm"
          onClick={() => handleGeneratePaths('nearest-neighbor')}
          disabled={generatePathsMutation.isPending}
          title="最近邻路径"
        >
          <GitBranch className="h-4 w-4" />
          <span className="ml-1 hidden md:inline">最近邻</span>
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={() => handleGeneratePaths('full-connectivity')}
          disabled={generatePathsMutation.isPending}
          title="完全连通"
        >
          <Circle className="h-4 w-4" />
          <span className="ml-1 hidden md:inline">完全连通</span>
        </Button>
      </div>
      
      <Separator orientation="vertical" className="h-6" />
      
      {/* 视图控制组 */}
      <div className="flex items-center gap-1">
        <Button
          variant="outline"
          size="sm"
          onClick={handleZoomOut}
          title="缩小"
        >
          <ZoomOut className="h-4 w-4" />
        </Button>
        <div className="px-2 py-1 text-sm text-gray-600 min-w-[60px] text-center">
          {Math.round(viewport.scale * 100)}%
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={handleZoomIn}
          title="放大"
        >
          <ZoomIn className="h-4 w-4" />
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={handleZoomFit}
          title="适应画布"
        >
          <Maximize className="h-4 w-4" />
        </Button>
      </div>
      
      <div className="flex-1" />
      
      {/* 右侧操作组 */}
      <div className="flex items-center gap-2">
        {/* 保存状态指示 */}
        {isDirty && (
          <Badge variant="warning" className="animate-pulse">
            未保存
          </Badge>
        )}
        
        {/* 操作按钮 */}
        <Button variant="outline" size="sm" title="导入">
          <Upload className="h-4 w-4" />
          <span className="ml-1 hidden sm:inline">导入</span>
        </Button>
        
        {/* Export Dropdown */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="sm" title="导出">
              <Download className="h-4 w-4" />
              <span className="ml-1 hidden sm:inline">导出</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => handleExport('nodes', 'csv')}>
              导出节点 (CSV)
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => handleExport('paths', 'csv')}>
              导出路径 (CSV)
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => handleExport('all', 'csv')}>
              导出所有 (CSV)
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => handleExport('nodes', 'xlsx')}>
              导出节点 (Excel)
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => handleExport('paths', 'xlsx')}>
              导出路径 (Excel)
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => handleExport('all', 'xlsx')}>
              导出所有 (Excel)
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        {/* Template Dropdown */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="sm" title="模板">
              <LayoutTemplate className="h-4 w-4" />
              <span className="ml-1 hidden sm:inline">模板</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={onOpenTemplateModal}>
              管理模板
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => console.log('Save current as template')}>
              保存当前为模板
            </DropdownMenuItem>
            {/* Add more template actions here */}
          </DropdownMenuContent>
        </DropdownMenu>
        
        <Button 
          size="sm" 
          disabled={!isDirty}
          title="保存更改"
        >
          <Save className="h-4 w-4" />
          <span className="ml-1 hidden sm:inline">保存</span>
        </Button>
      </div>
    </motion.div>
  )
}
