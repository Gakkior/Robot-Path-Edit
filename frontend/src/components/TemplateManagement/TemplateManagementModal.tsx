// src/components/TemplateManagement/TemplateManagementModal.tsx

'use client'

import React, { useState } from 'react'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/Dialog'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/Tabs'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Textarea } from '@/components/ui/Textarea' // Corrected import
import { Label } from '@/components/ui/Label'
import { useAppStore } from '@/stores/useAppStore'
import { useTemplates, useSaveAsTemplate, useApplyTemplate, useDeleteTemplate, useImportTemplate, useExportTemplate } from '@/services/api'
import { Template, SaveAsTemplateRequest } from '@/types'
import { toast } from '@/hooks/use-toast'
import { Loader2, Plus, Trash2, Download, Upload, Search } from 'lucide-react'

interface TemplateManagementModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const TemplateManagementModal: React.FC<TemplateManagementModalProps> = ({ isOpen, onClose }) => {
  const { nodes, paths, setNodes, setPaths } = useAppStore();
  const { data: templates, isLoading, isError, refetch } = useTemplates();
  const saveAsTemplateMutation = useSaveAsTemplate();
  const applyTemplateMutation = useApplyTemplate();
  const deleteTemplateMutation = useDeleteTemplate();
  const importTemplateMutation = useImportTemplate();
  const exportTemplateMutation = useExportTemplate();

  const [newTemplateName, setNewTemplateName] = useState('');
  const [newTemplateDescription, setNewTemplateDescription] = useState('');
  const [newTemplateCategory, setNewTemplateCategory] = useState('');
  const [newTemplateLayoutType, setNewTemplateLayoutType] = useState<Template['layoutType']>('custom');
  const [importContent, setImportContent] = useState('');
  const [searchTerm, setSearchTerm] = useState('');

  const handleSaveAsTemplate = async () => {
    if (!newTemplateName.trim()) {
      toast({
        title: '错误',
        description: '模板名称不能为空。',
        variant: 'destructive',
      });
      return;
    }

    const templateData: SaveAsTemplateRequest = {
      name: newTemplateName,
      description: newTemplateDescription,
      category: newTemplateCategory,
      layoutType: newTemplateLayoutType,
      nodes: nodes,
      paths: paths,
      isPublic: false, // Default to private
    };

    try {
      await saveAsTemplateMutation.mutateAsync(templateData);
      toast({
        title: '成功',
        description: '模板保存成功！',
      });
      setNewTemplateName('');
      setNewTemplateDescription('');
      setNewTemplateCategory('');
      refetch(); // Refresh template list
    } catch (error: any) {
      toast({
        title: '保存失败',
        description: error.message || '保存模板时发生错误。',
        variant: 'destructive',
      });
    }
  };

  const handleApplyTemplate = async (templateId: string) => {
    try {
      // Assuming canvas dimensions are fixed or can be retrieved
      const canvasWidth = window.innerWidth;
      const canvasHeight = window.innerHeight;
      const result = await applyTemplateMutation.mutateAsync({
        id: templateId,
        canvasConfig: { width: canvasWidth, height: canvasHeight },
      });
      setNodes(result.nodes);
      setPaths(result.paths);
      toast({
        title: '成功',
        description: '模板应用成功！画布已更新。',
      });
      onClose();
    } catch (error: any) {
      toast({
        title: '应用失败',
        description: error.message || '应用模板时发生错误。',
        variant: 'destructive',
      });
    }
  };

  const handleDeleteTemplate = async (templateId: string) => {
    if (!window.confirm('确定要删除此模板吗？')) {
      return;
    }
    try {
      await deleteTemplateMutation.mutateAsync(templateId);
      toast({
        title: '成功',
        description: '模板删除成功！',
      });
      refetch();
    } catch (error: any) {
      toast({
        title: '删除失败',
        description: error.message || '删除模板时发生错误。',
        variant: 'destructive',
      });
    }
  };

  const handleImportTemplate = async () => {
    if (!importContent.trim()) {
      toast({
        title: '错误',
        description: '导入内容不能为空。',
        variant: 'destructive',
      });
      return;
    }
    try {
      await importTemplateMutation.mutateAsync({ content: importContent });
      toast({
        title: '成功',
        description: '模板导入成功！',
      });
      setImportContent('');
      refetch();
    } catch (error: any) {
      toast({
        title: '导入失败',
        description: error.message || '导入模板时发生错误。',
        variant: 'destructive',
      });
    }
  };

  const handleExportTemplate = async (templateId: string) => {
    try {
      const exportedContent = await exportTemplateMutation.mutateAsync(templateId);
      // You might want to offer a download or display the content
      const blob = new Blob([exportedContent], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `template-${templateId}.json`;
      document.body.appendChild(a);
      a.click();
      a.remove();
      URL.revokeObjectURL(url);
      toast({
        title: '成功',
        description: '模板已导出并下载。',
      });
    } catch (error: any) {
      toast({
        title: '导出失败',
        description: error.message || '导出模板时发生错误。',
        variant: 'destructive',
      });
    }
  };

  const filteredTemplates = templates?.filter(template =>
    template.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    template.description?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    template.category?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[800px] h-[600px] flex flex-col">
        <DialogHeader>
          <DialogTitle>模板管理</DialogTitle>
          <DialogDescription>
            管理您的机器人路径模板，包括保存、应用、导入和导出。
          </DialogDescription>
        </DialogHeader>
        <Tabs defaultValue="my-templates" className="flex-1 flex flex-col">
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="my-templates">我的模板</TabsTrigger>
            <TabsTrigger value="save-template">保存当前为模板</TabsTrigger>
            <TabsTrigger value="import-export">导入/导出</TabsTrigger>
          </TabsList>
          <TabsContent value="my-templates" className="flex-1 flex flex-col overflow-hidden">
            <div className="relative mb-4">
              <Search className="absolute left-2 top-2.5 h-4 w-4 text-gray-500" />
              <Input
                placeholder="搜索模板..."
                className="pl-8"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </div>
            {isLoading && <div className="flex justify-center items-center flex-1"><Loader2 className="h-8 w-8 animate-spin" /> 正在加载模板...</div>}
            {isError && <div className="text-red-500 text-center flex-1 flex items-center justify-center">加载模板失败。</div>}
            {!isLoading && !isError && (
              <div className="flex-1 overflow-y-auto pr-2">
                {filteredTemplates && filteredTemplates.length > 0 ? (
                  <div className="grid gap-4">
                    {filteredTemplates.map((template) => (
                      <div key={template.id} className="border rounded-lg p-4 flex items-center justify-between shadow-sm">
                        <div>
                          <h3 className="font-semibold">{template.name}</h3>
                          <p className="text-sm text-gray-500">{template.description}</p>
                          <p className="text-xs text-gray-400">分类: {template.category || '无'}</p>
                          <p className="text-xs text-gray-400">布局: {template.layoutType}</p>
                        </div>
                        <div className="flex gap-2">
                          <Button size="sm" onClick={() => handleApplyTemplate(template.id)} disabled={applyTemplateMutation.isPending}>
                            {applyTemplateMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : '应用'}
                          </Button>
                          <Button size="sm" variant="outline" onClick={() => handleExportTemplate(template.id)} disabled={exportTemplateMutation.isPending}>
                            {exportTemplateMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Download className="h-4 w-4" />}
                          </Button>
                          <Button size="sm" variant="destructive" onClick={() => handleDeleteTemplate(template.id)} disabled={deleteTemplateMutation.isPending}>
                            {deleteTemplateMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-center text-gray-500 mt-8">没有找到模板。</p>
                )}
              </div>
            )}
          </TabsContent>
          <TabsContent value="save-template" className="flex-1 flex flex-col p-4 overflow-y-auto">
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="templateName" className="text-right">
                  模板名称
                </Label>
                <Input
                  id="templateName"
                  value={newTemplateName}
                  onChange={(e) => setNewTemplateName(e.target.value)}
                  className="col-span-3"
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="templateDescription" className="text-right">
                  描述
                </Label>
                <Textarea
                  id="templateDescription"
                  value={newTemplateDescription}
                  onChange={(e) => setNewTemplateDescription(e.target.value)}
                  className="col-span-3"
                  rows={3}
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="templateCategory" className="text-right">
                  分类
                </Label>
                <Input
                  id="templateCategory"
                  value={newTemplateCategory}
                  onChange={(e) => setNewTemplateCategory(e.target.value)}
                  className="col-span-3"
                  placeholder="例如: factory, warehouse"
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="templateLayoutType" className="text-right">
                  布局类型
                </Label>
                <select
                  id="templateLayoutType"
                  value={newTemplateLayoutType}
                  onChange={(e) => setNewTemplateLayoutType(e.target.value as Template['layoutType'])}
                  className="col-span-3 border rounded-md p-2"
                >
                  <option value="custom">自定义</option>
                  <option value="tree">树形图</option>
                  <option value="grid">网格</option>
                  <option value="circular">圆形</option>
                  <option value="force-directed">力导向</option>
                  <option value="pipeline">管道</option>
                  <option value="hierarchical">层次</option>
                  <option value="radial">径向</option>
                </select>
              </div>
            </div>
            <div className="flex justify-end mt-auto">
              <Button onClick={handleSaveAsTemplate} disabled={saveAsTemplateMutation.isPending}>
                {saveAsTemplateMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus className="h-4 w-4 mr-2" />}
                保存为模板
              </Button>
            </div>
          </TabsContent>
          <TabsContent value="import-export" className="flex-1 flex flex-col p-4 overflow-y-auto">
            <div className="grid gap-4 py-4">
              <h4 className="font-semibold">导入模板 (JSON)</h4>
              <Textarea
                placeholder="粘贴模板JSON内容到这里..."
                rows={8}
                value={importContent}
                onChange={(e) => setImportContent(e.target.value)}
              />
              <Button onClick={handleImportTemplate} disabled={importTemplateMutation.isPending}>
                {importTemplateMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Upload className="h-4 w-4 mr-2" />}
                导入模板
              </Button>
            </div>
            <Separator className="my-4" />
            <div className="grid gap-4 py-4">
              <h4 className="font-semibold">导出现有模板</h4>
              <p className="text-sm text-gray-500">从 "我的模板" 列表中选择一个模板进行导出。</p>
            </div>
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  );
};
