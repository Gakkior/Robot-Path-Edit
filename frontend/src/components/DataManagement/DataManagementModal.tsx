/**
 * 数据管理模态框
 * 
 * 提供数据库连接和表映射的配置界面
 */

import React, { useState } from 'react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/Dialog'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/Tabs'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/Select'
import { Separator } from '@/components/ui/Separator'
import { Badge } from '@/components/ui/Badge'
import {
  useDatabaseConnections,
  useCreateConnection,
  useUpdateConnection,
  useDeleteConnection,
  useTestConnection,
  useTableMappings,
  useCreateMapping,
  useUpdateMapping,
  useDeleteMapping,
  useSyncAllDataFromExternal, // Corrected import name
  useValidateExternalTable, // Corrected import name
} from '@/services' // Import from the barrel file
import type { DatabaseType, CreateConnectionRequest, CreateTableMappingRequest, TableMappingType } from '@/types'
import { Loader2, CheckCircle, XCircle, Plus, Trash2, RefreshCw, Database, Link, Play } from 'lucide-react'
import { useToast } from '@/hooks/use-toast' // Assuming useToast is available

interface DataManagementModalProps {
  isOpen: boolean
  onClose: () => void
}

export const DataManagementModal: React.FC<DataManagementModalProps> = ({ isOpen, onClose }) => {
  const [activeTab, setActiveTab] = useState('connections')

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-3xl h-[80vh] flex flex-col">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Database className="h-6 w-6 text-primary-600" />
            数据管理
          </DialogTitle>
          <DialogDescription>
            配置数据库连接和表映射，并同步数据到画布。
          </DialogDescription>
        </DialogHeader>

        <Tabs value={activeTab} onValueChange={setActiveTab} className="flex flex-col flex-1 overflow-hidden">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="connections">数据库连接</TabsTrigger>
            <TabsTrigger value="mappings">表映射</TabsTrigger>
          </TabsList>

          <TabsContent value="connections" className="flex-1 overflow-y-auto p-2">
            <DatabaseConnectionsTab />
          </TabsContent>
          <TabsContent value="mappings" className="flex-1 overflow-y-auto p-2">
            <TableMappingsTab />
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  )
}

const DatabaseConnectionsTab: React.FC = () => {
  const { data: connections, isLoading, error, refetch } = useDatabaseConnections()
  const createConnectionMutation = useCreateConnection()
  const deleteConnectionMutation = useDeleteConnection()
  const testConnectionMutation = useTestConnection()
  const { toast } = useToast()

  const [newConnection, setNewConnection] = useState<CreateConnectionRequest>({
    name: '',
    type: 'mysql',
    host: '',
    port: 3306,
    database: '',
    username: '',
    password: '',
  })

  const handleCreateConnection = async () => {
    try {
      await createConnectionMutation.mutateAsync(newConnection)
      toast({
        title: '连接创建成功',
        description: `数据库连接 "${newConnection.name}" 已创建。`,
        variant: 'success',
      })
      setNewConnection({ name: '', type: 'mysql', host: '', port: 3306, database: '', username: '', password: '' })
    } catch (err: any) {
      toast({
        title: '连接创建失败',
        description: err.message || '创建连接时发生错误。',
        variant: 'destructive',
      })
    }
  }

  const handleDeleteConnection = async (id: string) => {
    if (window.confirm('确定要删除此连接吗？')) {
      try {
        await deleteConnectionMutation.mutateAsync(id)
        toast({
          title: '连接删除成功',
          description: '数据库连接已删除。',
          variant: 'success',
        })
      } catch (err: any) {
        toast({
          title: '连接删除失败',
          description: err.message || '删除连接时发生错误。',
          variant: 'destructive',
        })
      }
    }
  }

  const handleTestConnection = async (id: string) => {
    try {
      const res = await testConnectionMutation.mutateAsync(id)
      if (res.success) {
        toast({
          title: '连接测试成功',
          description: '数据库连接正常。',
          variant: 'success',
        })
      } else {
        toast({
          title: '连接测试失败',
          description: res.message || '无法连接到数据库。',
          variant: 'destructive',
        })
      }
    } catch (err: any) {
      toast({
        title: '连接测试失败',
        description: err.message || '测试连接时发生错误。',
        variant: 'destructive',
      })
    }
  }

  if (isLoading) return <div className="text-center py-8 text-gray-500">加载连接中...</div>
  if (error) return <div className="text-center py-8 text-red-500">加载连接失败: {error.message}</div>

  return (
    <div className="space-y-6">
      <h3 className="text-lg font-semibold text-gray-900 flex items-center justify-between">
        现有数据库连接
        <Button variant="outline" size="sm" onClick={() => refetch()}>
          <RefreshCw className="h-4 w-4 mr-2" /> 刷新
        </Button>
      </h3>
      {connections && connections.length > 0 ? (
        <div className="space-y-4">
          {connections.map((conn) => (
            <div key={conn.id} className="border border-gray-200 rounded-lg p-4 shadow-sm flex items-center justify-between">
              <div>
                <div className="font-medium text-gray-900">{conn.name} <Badge variant="secondary">{conn.type}</Badge></div>
                <div className="text-sm text-gray-600">{conn.username}@{conn.host}:{conn.port}/{conn.database}</div>
              </div>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" onClick={() => handleTestConnection(conn.id)} disabled={testConnectionMutation.isPending}>
                  {testConnectionMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Play className="h-4 w-4" />}
                  <span className="ml-1">测试</span>
                </Button>
                <Button variant="destructive" size="sm" onClick={() => handleDeleteConnection(conn.id)} disabled={deleteConnectionMutation.isPending}>
                  {deleteConnectionMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
                  <span className="ml-1">删除</span>
                </Button>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <p className="text-gray-500 text-center py-4">暂无数据库连接。</p>
      )}

      <Separator />

      <h3 className="text-lg font-semibold text-gray-900">添加新连接</h3>
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">连接名称</label>
          <Input value={newConnection.name} onChange={(e) => setNewConnection({ ...newConnection, name: e.target.value })} />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">数据库类型</label>
          <Select value={newConnection.type} onValueChange={(value: DatabaseType) => setNewConnection({ ...newConnection, type: value })}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="mysql">MySQL</SelectItem>
              <SelectItem value="sqlite">SQLite</SelectItem>
              <SelectItem value="postgresql">PostgreSQL</SelectItem>
              <SelectItem value="sqlserver">SQL Server</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">主机</label>
          <Input value={newConnection.host} onChange={(e) => setNewConnection({ ...newConnection, host: e.target.value })} />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">端口</label>
          <Input type="number" value={newConnection.port} onChange={(e) => setNewConnection({ ...newConnection, port: parseInt(e.target.value) || 0 })} />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">数据库名</label>
          <Input value={newConnection.database} onChange={(e) => setNewConnection({ ...newConnection, database: e.target.value })} />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">用户名</label>
          <Input value={newConnection.username} onChange={(e) => setNewConnection({ ...newConnection, username: e.target.value })} />
        </div>
        <div className="col-span-2">
          <label className="block text-sm font-medium text-gray-700 mb-1">密码</label>
          <Input type="password" value={newConnection.password} onChange={(e) => setNewConnection({ ...newConnection, password: e.target.value })} />
        </div>
      </div>
      <Button onClick={handleCreateConnection} disabled={createConnectionMutation.isPending || !newConnection.name || !newConnection.host || !newConnection.database || !newConnection.username}>
        {createConnectionMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <Plus className="h-4 w-4 mr-2" />}
        创建连接
      </Button>
    </div>
  )
}

const TableMappingsTab: React.FC = () => {
  const { data: connections, isLoading: connectionsLoading } = useDatabaseConnections()
  const { data: mappings, isLoading: mappingsLoading, error: mappingsError, refetch } = useTableMappings()
  const createMappingMutation = useCreateMapping()
  const deleteMappingMutation = useDeleteMapping()
  const syncAllDataMutation = useSyncAllDataFromExternal() // Corrected hook name
  const validateTableMutation = useValidateExternalTable() // Corrected hook name
  const { toast } = useToast()

  const [newMapping, setNewMapping] = useState<CreateTableMappingRequest>({
    connection_id: '',
    table_name: '',
    type: 'node',
    node_mapping: { id_field: '', name_field: '', x_field: '', y_field: '' },
    path_mapping: undefined,
  })
  const [validationResult, setValidationResult] = useState<any>(null)
  const [validationLoading, setValidationLoading] = useState(false)

  const availableConnections = connections || []

  const handleCreateMapping = async () => {
    try {
      await createMappingMutation.mutateAsync(newMapping)
      toast({
        title: '映射创建成功',
        description: `表映射 "${newMapping.table_name}" 已创建。`,
        variant: 'success',
      })
      setNewMapping({
        connection_id: '',
        table_name: '',
        type: 'node',
        node_mapping: { id_field: '', name_field: '', x_field: '', y_field: '' },
        path_mapping: undefined,
      })
      setValidationResult(null)
    } catch (err: any) {
      toast({
        title: '映射创建失败',
        description: err.message || '创建映射时发生错误。',
        variant: 'destructive',
      })
    }
  }

  const handleDeleteMapping = async (id: string) => {
    if (window.confirm('确定要删除此表映射吗？')) {
      try {
        await deleteMappingMutation.mutateAsync(id)
        toast({
          title: '映射删除成功',
          description: '表映射已删除。',
          variant: 'success',
        })
      } catch (err: any) {
        toast({
          title: '映射删除失败',
          description: err.message || '删除映射时发生错误。',
          variant: 'destructive',
        })
      }
    }
  }

  const handleSyncData = async (mappingId: string) => {
    try {
      const res = await syncAllDataMutation.mutateAsync(mappingId)
      toast({
        title: '数据同步成功',
        description: res.message || `已同步 ${res.data?.nodes_synced || 0} 个节点和 ${res.data?.paths_synced || 0} 条路径。`,
        variant: 'success',
      })
    } catch (err: any) {
      toast({
        title: '数据同步失败',
        description: err.message || '同步数据时发生错误。',
        variant: 'destructive',
      })
    }
  }

  const handleValidateTable = async () => {
    if (!newMapping.connection_id || !newMapping.table_name) {
      toast({
        title: '验证失败',
        description: '请选择数据库连接和输入表名。',
        variant: 'warning',
      })
      return
    }
    setValidationLoading(true)
    try {
      const res = await validateTableMutation.mutateAsync({
        connectionId: newMapping.connection_id,
        tableName: newMapping.table_name,
      })
      setValidationResult(res.data)
      if (res.data?.is_valid) {
        toast({
          title: '表结构验证成功',
          description: res.data.message,
          variant: 'success',
        })
      } else {
        toast({
          title: '表结构验证失败',
          description: res.data?.message || '表结构不符合要求。',
          variant: 'destructive',
        })
      }
    } catch (err: any) {
      setValidationResult({ is_valid: false, message: err.message || '验证请求失败。' })
      toast({
        title: '表结构验证失败',
        description: err.message || '验证请求失败。',
        variant: 'destructive',
      })
    } finally {
      setValidationLoading(false)
    }
  }

  if (connectionsLoading || mappingsLoading) return <div className="text-center py-8 text-gray-500">加载映射中...</div>
  if (mappingsError) return <div className="text-center py-8 text-red-500">加载映射失败: {mappingsError.message}</div>

  return (
    <div className="space-y-6">
      <h3 className="text-lg font-semibold text-gray-900 flex items-center justify-between">
        现有表映射
        <Button variant="outline" size="sm" onClick={() => refetch()}>
          <RefreshCw className="h-4 w-4 mr-2" /> 刷新
        </Button>
      </h3>
      {mappings && mappings.length > 0 ? (
        <div className="space-y-4">
          {mappings.map((mapping) => (
            <div key={mapping.id} className="border border-gray-200 rounded-lg p-4 shadow-sm flex items-center justify-between">
              <div>
                <div className="font-medium text-gray-900">
                  {mapping.table_name} <Badge variant="secondary">{mapping.type === 'node' ? '节点表' : '路径表'}</Badge>
                </div>
                <div className="text-sm text-gray-600">
                  连接: {availableConnections.find(c => c.id === mapping.connection_id)?.name || mapping.connection_id.slice(0, 8)}...
                </div>
              </div>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" onClick={() => handleSyncData(mapping.id)} disabled={syncAllDataMutation.isPending}>
                  {syncAllDataMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />}
                  <span className="ml-1">同步数据</span>
                </Button>
                <Button variant="destructive" size="sm" onClick={() => handleDeleteMapping(mapping.id)} disabled={deleteMappingMutation.isPending}>
                  {deleteMappingMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
                  <span className="ml-1">删除</span>
                </Button>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <p className="text-gray-500 text-center py-4">暂无表映射。</p>
      )}

      <Separator />

      <h3 className="text-lg font-semibold text-gray-900">创建新映射</h3>
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">数据库连接</label>
          <Select value={newMapping.connection_id} onValueChange={(value) => setNewMapping({ ...newMapping, connection_id: value })}>
            <SelectTrigger>
              <SelectValue placeholder="选择一个连接" />
            </SelectTrigger>
            <SelectContent>
              {availableConnections.map(conn => (
                <SelectItem key={conn.id} value={conn.id}>{conn.name}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">表名</label>
          <Input value={newMapping.table_name} onChange={(e) => setNewMapping({ ...newMapping, table_name: e.target.value })} />
        </div>
        <div className="col-span-2">
          <Button variant="outline" size="sm" onClick={handleValidateTable} disabled={validationLoading || !newMapping.connection_id || !newMapping.table_name}>
            {validationLoading ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <CheckCircle className="h-4 w-4 mr-2" />}
            验证表结构
          </Button>
          {validationResult && (
            <div className={`mt-2 p-2 rounded-md text-sm ${validationResult.is_valid ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'}`}>
              {validationResult.is_valid ? <CheckCircle className="inline h-4 w-4 mr-1" /> : <XCircle className="inline h-4 w-4 mr-1" />}
              {validationResult.message}
              {validationResult.columns && validationResult.columns.length > 0 && (
                <div className="mt-1">
                  <span className="font-semibold">检测到的字段:</span> {validationResult.columns.map(col => col.name).join(', ')}
                </div>
              )}
              {validationResult.missing_fields && validationResult.missing_fields.length > 0 && (
                <div className="mt-1">
                  <span className="font-semibold">缺失的字段:</span> {validationResult.missing_fields.join(', ')}
                </div>
              )}
            </div>
          )}
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">映射类型</label>
          <Select value={newMapping.type} onValueChange={(value: TableMappingType) => {
            setNewMapping(prev => ({
              ...prev,
              type: value,
              node_mapping: value === 'node' ? { id_field: '', name_field: '', x_field: '', y_field: '' } : undefined,
              path_mapping: value === 'path' ? { id_field: '', start_node_field: '', end_node_field: '' } : undefined,
            }))
          }}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="node">节点表</SelectItem>
              <SelectItem value="path">路径表</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {newMapping.type === 'node' && newMapping.node_mapping && (
        <div className="space-y-4 mt-4 p-4 border border-gray-200 rounded-lg bg-gray-50">
          <h4 className="text-md font-semibold text-gray-800">节点字段映射</h4>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">ID 字段</label>
              <Input value={newMapping.node_mapping.id_field} onChange={(e) => setNewMapping(prev => ({ ...prev, node_mapping: { ...prev.node_mapping!, id_field: e.target.value } }))} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">名称 字段 (可选)</label>
              <Input value={newMapping.node_mapping.name_field || ''} onChange={(e) => setNewMapping(prev => ({ ...prev, node_mapping: { ...prev.node_mapping!, name_field: e.target.value } }))} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">类型 字段 (可选)</label>
              <Input value={newMapping.node_mapping.type_field || ''} onChange={(e) => setNewMapping(prev => ({ ...prev, node_mapping: { ...prev.node_mapping!, type_field: e.target.value } }))} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">X 坐标字段</label>
              <Input value={newMapping.node_mapping.x_field} onChange={(e) => setNewMapping(prev => ({ ...prev, node_mapping: { ...prev.node_mapping!, x_field: e.target.value } }))} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Y 坐标字段</label>
              <Input value={newMapping.node_mapping.y_field} onChange={(e) => setNewMapping(prev => ({ ...prev, node_mapping: { ...prev.node_mapping!, y_field: e.target.value } }))} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Z 坐标字段 (可选)</label>
              <Input value={newMapping.node_mapping.z_field || ''} onChange={(e) => setNewMapping(prev => ({ ...prev, node_mapping: { ...prev.node_mapping!, z_field: e.target.value } }))} />
            </div>
          </div>
        </div>
      )}

      {newMapping.type === 'path' && newMapping.path_mapping && (
        <div className="space-y-4 mt-4 p-4 border border-gray-200 rounded-lg bg-gray-50">
          <h4 className="text-md font-semibold text-gray-800">路径字段映射</h4>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">ID 字段</label>
              <Input value={newMapping.path_mapping.id_field} onChange={(e) => setNewMapping(prev => ({ ...prev, path_mapping: { ...prev.path_mapping!, id_field: e.target.value } }))} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">名称 字段 (可选)</label>
              <Input value={newMapping.path_mapping.name_field || ''} onChange={(e) => setNewMapping(prev => ({ ...prev, path_mapping: { ...prev.path_mapping!, name_field: e.target.value } }))} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">起始节点 ID 字段</label>
              <Input value={newMapping.path_mapping.start_node_field} onChange={(e) => setNewMapping(prev => ({ ...prev, path_mapping: { ...prev.path_mapping!, start_node_field: e.target.value } }))} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">结束节点 ID 字段</label>
              <Input value={newMapping.path_mapping.end_node_field} onChange={(e) => setNewMapping(prev => ({ ...prev, path_mapping: { ...prev.path_mapping!, end_node_field: e.target.value } }))} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">权重 字段 (可选)</label>
              <Input value={newMapping.path_mapping.weight_field || ''} onChange={(e) => setNewMapping(prev => ({ ...prev, path_mapping: { ...prev.path_mapping!, weight_field: e.target.value } }))} />
            </div>
          </div>
        </div>
      )}

      <Button onClick={handleCreateMapping} disabled={createMappingMutation.isPending || !newMapping.connection_id || !newMapping.table_name || (newMapping.type === 'node' && (!newMapping.node_mapping?.id_field || !newMapping.node_mapping?.x_field || !newMapping.node_mapping?.y_field)) || (newMapping.type === 'path' && (!newMapping.path_mapping?.id_field || !newMapping.path_mapping?.start_node_field || !newMapping.path_mapping?.end_node_field))}>
        {createMappingMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <Plus className="h-4 w-4 mr-2" />}
        创建映射
      </Button>
    </div>
  )
}
