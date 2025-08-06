// src/components/Layout/AppLayout.tsx

'use client'

import React, { useState } from 'react'
import { Canvas } from '@/components/Canvas/Canvas'
import { Toolbar } from '@/components/Toolbar/Toolbar'
import { PropertyPanel } from '@/components/Sidebar/PropertyPanel'
import { DataTable } from '@/components/Table/DataTable'
import { DataManagementModal } from '@/components/DataManagement/DataManagementModal'
import { TemplateManagementModal } from '@/components/TemplateManagement/TemplateManagementModal'
import { useAppStore } from '@/stores/useAppStore'
import { useHealth } from '@/services/api'
import { Wifi, WifiOff, Database } from 'lucide-react'
import { Button } from '@/components/ui/Button'
import { Separator } from '@/components/ui/Separator'
import { Badge } from '@/components/ui/Badge'

const AppLayout: React.FC = () => { // Changed to a regular const declaration
  const { selectedNodeId, selectedPathId, editMode, isConnecting, connectingNodeId } = useAppStore()
  const { data: health, isLoading: healthLoading, isError: healthError } = useHealth();

  const [isDataManagementModalOpen, setIsDataManagementModalOpen] = useState(false);
  const [isTemplateManagementModalOpen, setIsTemplateManagementModalOpen] = useState(false);

  const [activeView, setActiveView] = useState<'canvas' | 'table'>('canvas'); // State to manage active view

  const handleOpenDataManagementModal = () => {
    setIsDataManagementModalOpen(true);
  };

  const handleCloseDataManagementModal = () => {
    setIsDataManagementModalOpen(false);
  };

  const handleOpenTemplateModal = () => {
    setIsTemplateManagementModalOpen(true);
  };

  const handleCloseTemplateModal = () => {
    setIsTemplateManagementModalOpen(false);
  };

  const getSelectModeIndicator = () => {
    if (editMode === 'add-path' && isConnecting) {
      return `连接模式: 从节点 ${connectingNodeId?.substring(0, 4)}...`;
    }
    switch (editMode) {
      case 'select':
        return '选择模式';
      case 'add-node':
        return '添加节点模式';
      case 'add-path':
        return '添加路径模式';
      case 'delete':
        return '删除模式';
      default:
        return '';
    }
  };

  const healthStatusColor = healthError ? 'bg-red-500' : (health?.status === 'ok' ? 'bg-green-500' : 'bg-yellow-500');
  const healthStatusText = healthError ? '断开连接' : (health?.status === 'ok' ? '已连接' : '连接中...');
  const HealthIcon = healthError ? WifiOff : Wifi;

  return (
    <div className="flex flex-col h-screen bg-gray-50">
      {/* Header */}
      <header className="flex items-center justify-between p-4 bg-white border-b border-gray-200 shadow-sm">
        <h1 className="text-xl font-bold text-gray-800">🤖 机器人路径编辑器</h1>
        <div className="flex items-center gap-4">
          {/* Connection Status */}
          <div className="flex items-center gap-2 text-sm text-gray-600">
            <HealthIcon className={`h-5 w-5 ${healthStatusColor} rounded-full p-1 text-white`} />
            <span>后端: {healthStatusText}</span>
          </div>
          <Button variant="outline" size="sm" onClick={handleOpenDataManagementModal}>
            <Database className="h-4 w-4 mr-2" />
            数据源
          </Button>
          <Button variant="outline" size="sm" onClick={() => setActiveView(activeView === 'canvas' ? 'table' : 'canvas')}>
            {activeView === 'canvas' ? '表格视图' : '画布视图'}
          </Button>
        </div>
      </header>

      {/* Toolbar */}
      <Toolbar onOpenTemplateModal={handleOpenTemplateModal} />

      {/* Main Content Area */}
      <main className="flex flex-1 overflow-hidden">
        {/* Canvas / Table View */}
        <div className="flex-1 relative">
          {activeView === 'canvas' ? <Canvas /> : <DataTable />}
          {/* Select Mode Indicator */}
          <div className="absolute bottom-4 left-4 bg-white px-3 py-1.5 rounded-full shadow-md text-sm font-medium text-gray-700">
            {getSelectModeIndicator()}
          </div>
        </div>

        {/* Property Panel */}
        {(selectedNodeId || selectedPathId) && (
          <PropertyPanel />
        )}
      </main>

      {/* Modals */}
      <DataManagementModal
        isOpen={isDataManagementModalOpen}
        onClose={handleCloseDataManagementModal}
      />
      <TemplateManagementModal
        isOpen={isTemplateManagementModalOpen}
        onClose={handleCloseTemplateModal}
      />
    </div>
  );
};

export default AppLayout; // Changed to default export
