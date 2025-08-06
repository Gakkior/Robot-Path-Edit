import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
 Node, Path, ListResponse, DatabaseConnection, TableMapping,
 HealthStatus, SystemStats, Template, SaveAsTemplateRequest, ApplyTemplateResponse, ExportDataRequest,
 CreatePathRequest, UpdatePathRequest
} from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

const fetcher = async <T>(url: string, options?: RequestInit): Promise<T> => {
 const response = await fetch(`${API_BASE_URL}${url}`, options);
 if (!response.ok) {
   const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
   throw new Error(errorData.error || response.statusText);
 }
 return response.json();
};

const poster = async <T, U>(url: string, data: U, options?: RequestInit): Promise<T> => {
 return fetcher<T>(url, {
   method: 'POST',
   headers: {
     'Content-Type': 'application/json',
     ...options?.headers,
   },
   body: JSON.stringify(data),
   ...options,
 });
};

const putter = async <T, U>(url: string, data: U, options?: RequestInit): Promise<T> => {
 return fetcher<T>(url, {
   method: 'PUT',
   headers: {
     'Content-Type': 'application/json',
     ...options?.headers,
   },
   body: JSON.stringify(data),
   ...options,
 });
};

const deleter = async <T>(url: string, options?: RequestInit): Promise<T> => {
 return fetcher<T>(url, {
   method: 'DELETE',
   ...options,
 });
};

// --- System API ---
export const systemApi = {
 health: () => fetcher<HealthStatus>('/health'),
 stats: () => fetcher<SystemStats>('/stats'),
};

export const useHealth = () => {
 return useQuery<HealthStatus, Error>({
   queryKey: ['health'],
   queryFn: systemApi.health,
   refetchInterval: 5000, // Poll every 5 seconds
   staleTime: 1000,
   retry: false, // Don't retry on network errors for health check
 });
};

export const useStats = () => {
 return useQuery<SystemStats, Error>({
   queryKey: ['stats'],
   queryFn: systemApi.stats,
   refetchInterval: 10000, // Poll every 10 seconds
   staleTime: 5000,
 });
};

// --- Node API ---
export const nodeApi = {
 list: () => fetcher<{ nodes: Node[] }>('/nodes').then(res => res.nodes),
 create: (node: Omit<Node, 'id' | 'created_at' | 'updated_at' | 'status'>) => poster<{ node: Node }, Omit<Node, 'id' | 'created_at' | 'updated_at' | 'status'>>('/nodes', node).then(res => res.node),
 update: (node: Partial<Node> & { id: string }) => putter<{ node: Node }, Partial<Node>>(`/nodes/${node.id}`, node).then(res => res.node),
 delete: (id: string) => deleter<{ message: string }>(`/nodes/${id}`),
 batchCreate: (nodes: Omit<Node, 'id' | 'created_at' | 'updated_at' | 'status'>[]) => poster<{ nodes: Node[] }, Omit<Node, 'id' | 'created_at' | 'updated_at' | 'status'>[]>('/nodes/batch', nodes).then(res => res.nodes),
 batchUpdate: (nodes: (Partial<Node> & { id: string })[]) => putter<{ nodes: Node[] }, (Partial<Node> & { id: string })[]>('/nodes/batch', nodes).then(res => res.nodes),
 batchDelete: (ids: string[]) => poster<{ message: string }, { ids: string[] }>('/nodes/batch-delete', { ids }),
};

export const useNodes = () => useQuery<Node[], Error>({ queryKey: ['nodes'], queryFn: nodeApi.list });
export const useCreateNode = () => {
 const queryClient = useQueryClient();
 return useMutation<Node, Error, Omit<Node, 'id' | 'created_at' | 'updated_at' | 'status'>>({
   mutationFn: nodeApi.create,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['nodes'] });
   },
 });
};
export const useUpdateNode = () => {
 const queryClient = useQueryClient();
 return useMutation<Node, Error, Partial<Node> & { id: string }>({
   mutationFn: nodeApi.update,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['nodes'] });
   },
 });
};
export const useDeleteNode = () => {
 const queryClient = useQueryClient();
 return useMutation<any, Error, string>({
   mutationFn: nodeApi.delete,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['nodes'] });
   },
 });
};

// --- Path API ---
export const pathApi = {
 list: () => fetcher<{ paths: Path[] }>('/paths').then(res => res.paths),
 create: (path: CreatePathRequest) => poster<{ path: Path }, CreatePathRequest>('/paths', path).then(res => res.path),
 update: (path: UpdatePathRequest) => putter<{ path: Path }, UpdatePathRequest>(`/paths/${path.id}`, path).then(res => res.path),
 delete: (id: string) => deleter<{ message: string }>(`/paths/${id}`),
};

export const usePaths = () => useQuery<Path[], Error>({ queryKey: ['paths'], queryFn: pathApi.list });
export const useCreatePath = () => {
 const queryClient = useQueryClient();
 return useMutation<Path, Error, CreatePathRequest>({
   mutationFn: pathApi.create,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['paths'] });
   },
 });
};
export const useUpdatePath = () => {
 const queryClient = useQueryClient();
 return useMutation<Path, Error, UpdatePathRequest>({
   mutationFn: pathApi.update,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['paths'] });
   },
 });
};
export const useDeletePath = () => {
 const queryClient = useQueryClient();
 return useMutation<any, Error, string>({
   mutationFn: pathApi.delete,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['paths'] });
   },
 });
};

// --- Layout API ---
export const layoutApi = {
 apply: (algorithm: string) => poster<any, { algorithm: string }>('/layout/apply', { algorithm }),
};

export const useApplyLayout = () => {
 const queryClient = useQueryClient();
 return useMutation<any, Error, { algorithm: string }>({
   mutationFn: layoutApi.apply,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['nodes'] });
     queryClient.invalidateQueries({ queryKey: ['paths'] });
   },
 });
};

// --- Path Generation API ---
export const pathGenerationApi = {
 generateNearestNeighbor: () => poster<any, {}>(`/path-generation/nearest-neighbor`, {}),
 generateFullConnectivity: () => poster<any, {}>(`/path-generation/full-connectivity`, {}),
};

export const useGeneratePaths = () => {
 const queryClient = useQueryClient();
 return useMutation<any, Error, { algorithm: string }>({
   mutationFn: ({ algorithm }) => {
     if (algorithm === 'nearest-neighbor') {
       return pathGenerationApi.generateNearestNeighbor();
     } else if (algorithm === 'full-connectivity') {
       return pathGenerationApi.generateFullConnectivity();
     }
     return Promise.reject(new Error('Unknown path generation algorithm'));
   },
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['paths'] });
   },
 });
};

// --- Database Connection API ---
export const databaseApi = {
 listConnections: () => fetcher<{ connections: DatabaseConnection[] }>('/database/connections').then(res => res.connections),
 createConnection: (data: Omit<DatabaseConnection, 'id' | 'createdAt' | 'updatedAt'>) => poster<{ connection: DatabaseConnection }, Omit<DatabaseConnection, 'id' | 'createdAt' | 'updatedAt'>>('/database/connections', data).then(res => res.connection),
 updateConnection: (id: string, data: Partial<DatabaseConnection>) => putter<{ connection: DatabaseConnection }, Partial<DatabaseConnection>>(`/database/connections/${id}`, data).then(res => res.connection),
 deleteConnection: (id: string) => deleter<{ message: string }>(`/database/connections/${id}`),
 testConnection: (id: string) => poster<{ message: string }, {}>(`/database/connections/${id}/test`, {}),

 listTableMappings: () => fetcher<{ mappings: TableMapping[] }>('/database/mappings').then(res => res.mappings),
 createTableMapping: (data: Omit<TableMapping, 'id' | 'createdAt' | 'updatedAt'>) => poster<{ mapping: TableMapping }, Omit<TableMapping, 'id' | 'createdAt' | 'updatedAt'>>('/database/mappings', data).then(res => res.mapping),
 updateTableMapping: (id: string, data: Partial<TableMapping>) => putter<{ mapping: TableMapping }, Partial<TableMapping>>(`/database/mappings/${id}`, data).then(res => res.mapping),
 deleteTableMapping: (id: string) => deleter<{ message: string }>(`/database/mappings/${id}`),

 syncNodesFromExternal: (mappingId: string) => poster<any, {}>(`/sync/mappings/${mappingId}/nodes`, {}),
 syncPathsFromExternal: (mappingId: string) => poster<any, {}>(`/sync/mappings/${mappingId}/paths`, {}),
 syncAllDataFromExternal: (mappingId: string) => poster<any, {}>(`/sync/mappings/${mappingId}/all`, {}),
 validateExternalTable: (connectionId: string, tableName: string) => fetcher<any>(`/database/validate-table?connection_id=${connectionId}&table_name=${tableName}`),
};

export const useDatabaseConnections = () => useQuery<DatabaseConnection[], Error>({ queryKey: ['databaseConnections'], queryFn: databaseApi.listConnections });
export const useCreateDatabaseConnection = () => {
 const queryClient = useQueryClient();
 return useMutation<DatabaseConnection, Error, Omit<DatabaseConnection, 'id' | 'createdAt' | 'updatedAt'>>({
   mutationFn: databaseApi.createConnection,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['databaseConnections'] });
   },
 });
};
export const useUpdateDatabaseConnection = () => {
 const queryClient = useQueryClient();
 return useMutation<DatabaseConnection, Error, { id: string; data: Partial<DatabaseConnection> }>({
   mutationFn: ({ id, data }) => databaseApi.updateConnection(id, data),
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['databaseConnections'] });
   },
 });
};
export const useDeleteDatabaseConnection = () => {
 const queryClient = useQueryClient();
 return useMutation<any, Error, string>({
   mutationFn: databaseApi.deleteConnection,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['databaseConnections'] });
   },
 });
};
export const useTestDatabaseConnection = () => {
 return useMutation<any, Error, string>({
   mutationFn: databaseApi.testConnection,
 });
};

export const useTableMappings = () => useQuery<TableMapping[], Error>({ queryKey: ['tableMappings'], queryFn: databaseApi.listTableMappings });
export const useCreateTableMapping = () => {
 const queryClient = useQueryClient();
 return useMutation<TableMapping, Error, Omit<TableMapping, 'id' | 'createdAt' | 'updatedAt'>>({
   mutationFn: databaseApi.createTableMapping,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['tableMappings'] });
   },
 });
};
export const useUpdateTableMapping = () => {
 const queryClient = useQueryClient();
 return useMutation<TableMapping, Error, { id: string; data: Partial<TableMapping> }>({
   mutationFn: ({ id, data }) => databaseApi.updateTableMapping(id, data),
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['tableMappings'] });
   },
 });
};
export const useDeleteTableMapping = () => {
 const queryClient = useQueryClient();
 return useMutation<any, Error, string>({
   mutationFn: databaseApi.deleteTableMapping,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['tableMappings'] });
   },
 });
};

export const useSyncNodesFromExternal = () => {
 const queryClient = useQueryClient();
 return useMutation<any, Error, string>({
   mutationFn: databaseApi.syncNodesFromExternal,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['nodes'] });
   },
 });
};
export const useSyncPathsFromExternal = () => {
 const queryClient = useQueryClient();
 return useMutation<any, Error, string>({
   mutationFn: databaseApi.syncPathsFromExternal,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['paths'] });
   },
 });
};
export const useSyncAllDataFromExternal = () => {
 const queryClient = useQueryClient();
 return useMutation<any, Error, string>({
   mutationFn: databaseApi.syncAllDataFromExternal,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['nodes'] });
     queryClient.invalidateQueries({ queryKey: ['paths'] });
   },
 });
};
export const useValidateExternalTable = () => {
 return useMutation<any, Error, { connectionId: string; tableName: string }>({
   mutationFn: ({ connectionId, tableName }) => databaseApi.validateExternalTable(connectionId, tableName),
 });
};

// --- Template API ---
export const templateApi = {
 listTemplates: () => fetcher<{ templates: Template[] }>('/templates').then(res => res.templates),
 createTemplate: (data: SaveAsTemplateRequest) => poster<{ template: Template }, SaveAsTemplateRequest>('/templates', data).then(res => res.template),
 saveAsTemplate: (data: SaveAsTemplateRequest) => poster<{ template: Template }, SaveAsTemplateRequest>('/templates/save-as', data).then(res => res.template),
 applyTemplate: (id: string, canvasConfig: { width: number; height: number }) => poster<{ result: ApplyTemplateResponse }, { width: number; height: number }>(`/templates/${id}/apply`, canvasConfig).then(res => res.result),
 deleteTemplate: (id: string) => deleter<{ message: string }>(`/templates/${id}`),
 importTemplate: (data: { content: string }) => poster<{ template: Template }, { content: string }>('/templates/import', data).then(res => res.template),
 exportTemplate: (id: string) => fetcher<{ export: string }>(`/templates/${id}/export`).then(res => res.export),
 getPublicTemplates: () => fetcher<{ templates: Template[] }>('/templates/public').then(res => res.templates),
 getTemplatesByCategory: (category: string) => fetcher<{ templates: Template[] }>(`/templates/category/${category}`).then(res => res.templates),
 searchTemplates: (query: string) => fetcher<{ templates: Template[] }>(`/templates/search?q=${query}`).then(res => res.templates),
};

export const useTemplates = () => useQuery<Template[], Error>({ queryKey: ['templates'], queryFn: templateApi.listTemplates });
export const useCreateTemplate = () => {
 const queryClient = useQueryClient();
 return useMutation<Template, Error, SaveAsTemplateRequest>({
   mutationFn: templateApi.createTemplate,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['templates'] });
   },
 });
};
export const useSaveAsTemplate = () => {
 const queryClient = useQueryClient();
 return useMutation<Template, Error, SaveAsTemplateRequest>({
   mutationFn: templateApi.saveAsTemplate,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['templates'] });
   },
 });
};
export const useApplyTemplate = () => {
 const queryClient = useQueryClient();
 return useMutation<ApplyTemplateResponse, Error, { id: string; canvasConfig: { width: number; height: number } }>({
   mutationFn: ({ id, canvasConfig }) => templateApi.applyTemplate(id, canvasConfig),
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['nodes'] });
     queryClient.invalidateQueries({ queryKey: ['paths'] });
   },
 });
};
export const useDeleteTemplate = () => {
 const queryClient = useQueryClient();
 return useMutation<any, Error, string>({
   mutationFn: templateApi.deleteTemplate,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['templates'] });
   },
 });
};
export const useImportTemplate = () => {
 const queryClient = useQueryClient();
 return useMutation<Template, Error, { content: string }>({
   mutationFn: templateApi.importTemplate,
   onSuccess: () => {
     queryClient.invalidateQueries({ queryKey: ['templates'] });
   },
 });
};
export const useExportTemplate = () => {
 return useMutation<string, Error, string>({
   mutationFn: templateApi.exportTemplate,
 });
};
export const usePublicTemplates = () => useQuery<Template[], Error>({ queryKey: ['publicTemplates'], queryFn: templateApi.getPublicTemplates });
export const useTemplatesByCategory = (category: string) => useQuery<Template[], Error>({ queryKey: ['templatesByCategory', category], queryFn: () => templateApi.getTemplatesByCategory(category) });
export const useSearchTemplates = (query: string) => useQuery<Template[], Error>({ queryKey: ['searchTemplates', query], queryFn: () => templateApi.searchTemplates(query) });


// --- Data Export API ---
export const exportApi = {
 exportData: (req: ExportDataRequest) => poster<Blob, ExportDataRequest>('/export', req, {
   headers: {
     'Content-Type': 'application/json',
     'Accept': req.format === 'csv' ? 'text/csv' : 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
   },
   responseType: 'blob', // Indicate that the response is a blob
 }),
};

export const useExportData = () => {
 return useMutation<Blob, Error, ExportDataRequest>({
   mutationFn: async (req) => {
     const response = await fetch(`${API_BASE_URL}/export`, {
       method: 'POST',
       headers: {
         'Content-Type': 'application/json',
         'Accept': req.format === 'csv' ? 'text/csv' : 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
       },
       body: JSON.stringify(req),
     });

     if (!response.ok) {
       const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
       throw new Error(errorData.error || response.statusText);
     }

     const blob = await response.blob();
     const filename = `data.${req.format}`;
     const url = window.URL.createObjectURL(blob);
     const a = document.createElement('a');
     a.href = url;
     a.download = filename;
     document.body.appendChild(a);
     a.click();
     a.remove();
     window.URL.revokeObjectURL(url);
     return blob;
   },
 });
};
