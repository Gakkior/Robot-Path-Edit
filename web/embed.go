// Package web 提供内嵌的前端资源
package web

import (
	"embed"
	"net/http"
)

//go:embed static/*
var StaticFiles embed.FS

// IndexHTML 主页面HTML
var IndexHTML = []byte(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>机器人路径编辑器</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        
        .container {
            text-align: center;
            padding: 2rem;
            background: rgba(255, 255, 255, 0.95);
            border-radius: 20px;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
            backdrop-filter: blur(10px);
            max-width: 600px;
        }
        
        .logo {
            font-size: 3rem;
            margin-bottom: 1rem;
        }
        
        h1 {
            color: #2c3e50;
            margin-bottom: 1rem;
            font-size: 2.5rem;
            font-weight: 300;
        }
        
        .subtitle {
            color: #7f8c8d;
            margin-bottom: 2rem;
            font-size: 1.2rem;
        }
        
        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
            margin: 2rem 0;
        }
        
        .feature {
            padding: 1rem;
            background: rgba(52, 152, 219, 0.1);
            border-radius: 10px;
            border: 1px solid rgba(52, 152, 219, 0.2);
        }
        
        .feature h3 {
            color: #3498db;
            margin-bottom: 0.5rem;
        }
        
        .status {
            display: inline-block;
            padding: 0.5rem 1rem;
            background: #2ecc71;
            color: white;
            border-radius: 25px;
            font-weight: 500;
            margin: 1rem 0;
        }
        
        .api-info {
            margin-top: 2rem;
            padding: 1rem;
            background: rgba(241, 196, 15, 0.1);
            border-radius: 10px;
            border: 1px solid rgba(241, 196, 15, 0.2);
        }
        
        .api-info h3 {
            color: #f39c12;
            margin-bottom: 0.5rem;
        }
        
        .api-endpoint {
            font-family: 'Monaco', 'Menlo', monospace;
            background: rgba(0, 0, 0, 0.05);
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            margin: 0.25rem 0;
            display: inline-block;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">🤖</div>
        <h1>机器人路径编辑器</h1>
        <p class="subtitle">现代化的三端兼容路径管理工具</p>
        
        <div class="status">✅ 服务运行中</div>
        
        <div class="features">
            <div class="feature">
                <h3>📊 可视化编辑</h3>
                <p>拖拽式画布编辑，直观的节点和路径管理</p>
            </div>
            <div class="feature">
                <h3>🗄️ 数据库支持</h3>
                <p>支持SQLite、MySQL等多种数据库</p>
            </div>
            <div class="feature">
                <h3>📱 三端兼容</h3>
                <p>支持Windows、Linux、Android平台</p>
            </div>
            <div class="feature">
                <h3>🎨 现代设计</h3>
                <p>美观流畅的现代化界面设计</p>
            </div>
        </div>
        
        <div class="api-info">
            <h3>🚀 开始使用</h3>
            <p>点击下面的按钮进入编辑器界面：</p>
            <div style="margin: 1rem 0;">
                <a href="/app" style="display: inline-block; background: #667eea; color: white; text-decoration: none; padding: 0.75rem 1.5rem; border-radius: 8px; font-weight: 500; transition: background-color 0.2s;" onmouseover="this.style.background='#5a67d8'" onmouseout="this.style.background='#667eea'">
                    🎨 打开编辑器
                </a>
            </div>
            <hr style="margin: 1.5rem 0; border: none; border-top: 1px solid #e1e8ed;">
            <h3>API 端点</h3>
            <p>RESTful API 服务已启动，可访问以下端点：</p>
            <div><code class="api-endpoint">GET /api/v1/nodes</code> - 获取节点列表</div>
            <div><code class="api-endpoint">GET /api/v1/paths</code> - 获取路径列表</div>
            <div><code class="api-endpoint">GET /health</code> - 健康检查</div>
            <div><code class="api-endpoint">GET /metrics</code> - 系统指标</div>
        </div>
    </div>
    
    <script>
        // 简单的API测试
        fetch('/health')
            .then(response => response.json())
            .then(data => {
                console.log('Health check:', data);
            })
            .catch(error => {
                console.error('Error:', error);
            });
            
        // 实时状态更新
        setInterval(() => {
            fetch('/health')
                .then(response => {
                    if (response.ok) {
                        document.querySelector('.status').textContent = '✅ 服务运行中';
                        document.querySelector('.status').style.background = '#2ecc71';
                    } else {
                        document.querySelector('.status').textContent = '❌ 服务异常';
                        document.querySelector('.status').style.background = '#e74c3c';
                    }
                })
                .catch(() => {
                    document.querySelector('.status').textContent = '❌ 连接失败';
                    document.querySelector('.status').style.background = '#e74c3c';
                });
        }, 10000); // 10秒检查一次
    </script>
</body>
</html>`)

// AppHTML 主应用界面HTML
var AppHTML = []byte(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>机器人路径编辑器 - 应用界面</title>
    <link rel="icon" type="image/svg+xml" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>🤖</text></svg>">
    <script src="https://unpkg.com/konva@9/konva.min.js"></script>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f5f7fa;
            height: 100vh;
            overflow: hidden;
        }
        
        .app-container {
            display: flex;
            flex-direction: column;
            height: 100vh;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 0.75rem 1rem;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            display: flex;
            justify-content: between;
            align-items: center;
            z-index: 1000;
        }
        
        .header h1 {
            font-size: 1.25rem;
            font-weight: 500;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
        
        .header-controls {
            display: flex;
            gap: 1rem;
            align-items: center;
            margin-left: auto;
        }
        
        .view-toggle {
            display: flex;
            background: rgba(255,255,255,0.2);
            border-radius: 8px;
            overflow: hidden;
        }
        
        .view-toggle button {
            background: none;
            border: none;
            color: white;
            padding: 0.5rem 1rem;
            cursor: pointer;
            transition: background-color 0.2s;
        }
        
        .view-toggle button.active {
            background: rgba(255,255,255,0.3);
        }
        
        .view-toggle button:hover {
            background: rgba(255,255,255,0.2);
        }
        
        .main-content {
            flex: 1;
            display: flex;
            overflow: hidden;
        }
        
        .sidebar {
            width: 300px;
            background: white;
            border-right: 1px solid #e1e8ed;
            display: flex;
            flex-direction: column;
            box-shadow: 2px 0 4px rgba(0,0,0,0.05);
        }
        
        .sidebar-header {
            padding: 1rem;
            border-bottom: 1px solid #e1e8ed;
            background: #f8f9fa;
        }
        
        .sidebar-content {
            flex: 1;
            padding: 1rem;
            overflow-y: auto;
        }
        
        .work-area {
            flex: 1;
            display: flex;
            flex-direction: column;
            position: relative;
        }
        
        .toolbar {
            background: white;
            border-bottom: 1px solid #e1e8ed;
            padding: 0.5rem 1rem;
            display: flex;
            gap: 0.5rem;
            align-items: center;
            box-shadow: 0 1px 3px rgba(0,0,0,0.05);
        }
        
        .toolbar button {
            background: #667eea;
            color: white;
            border: none;
            padding: 0.5rem 0.75rem;
            border-radius: 6px;
            cursor: pointer;
            font-size: 0.875rem;
            transition: background-color 0.2s;
            display: flex;
            align-items: center;
            gap: 0.25rem;
        }
        
        .toolbar button:hover {
            background: #5a67d8;
        }
        
        .toolbar .separator {
            width: 1px;
            height: 24px;
            background: #e1e8ed;
            margin: 0 0.5rem;
        }
        
        .canvas-container {
            flex: 1;
            background: white;
            position: relative;
            overflow: hidden;
        }
        
        #canvas-stage {
            background: 
                radial-gradient(circle, #e2e8f0 1px, transparent 1px);
            background-size: 20px 20px;
            cursor: crosshair;
        }
        
        .table-container {
            flex: 1;
            background: white;
            display: none;
            flex-direction: column;
        }
        
        .table-wrapper {
            flex: 1;
            overflow: auto;
            padding: 1rem;
        }
        
        .data-table {
            width: 100%;
            border-collapse: collapse;
            background: white;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        
        .data-table th {
            background: #f8f9fa;
            padding: 0.75rem;
            text-align: left;
            font-weight: 600;
            border-bottom: 2px solid #e1e8ed;
            position: sticky;
            top: 0;
            z-index: 10;
        }
        
        .data-table td {
            padding: 0.75rem;
            border-bottom: 1px solid #e1e8ed;
        }
        
        .data-table tr:hover {
            background: #f8f9fa;
        }
        
        .status-bar {
            background: #f8f9fa;
            border-top: 1px solid #e1e8ed;
            padding: 0.5rem 1rem;
            font-size: 0.875rem;
            color: #6c757d;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .property-panel {
            background: white;
            border: 1px solid #e1e8ed;
            border-radius: 8px;
            margin-bottom: 1rem;
        }
        
        .property-panel h3 {
            background: #f8f9fa;
            padding: 0.75rem;
            margin: 0;
            border-bottom: 1px solid #e1e8ed;
            font-size: 0.875rem;
            font-weight: 600;
        }
        
        .property-content {
            padding: 1rem;
        }
        
        .form-group {
            margin-bottom: 1rem;
        }
        
        .form-group label {
            display: block;
            margin-bottom: 0.25rem;
            font-size: 0.875rem;
            font-weight: 500;
            color: #374151;
        }
        
        .form-group input,
        .form-group select,
        .form-group textarea {
            width: 100%;
            padding: 0.5rem;
            border: 1px solid #d1d5db;
            border-radius: 4px;
            font-size: 0.875rem;
        }
        
        .form-group input:focus,
        .form-group select:focus,
        .form-group textarea:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
        }
        
        .btn-group {
            display: flex;
            gap: 0.5rem;
            margin-top: 1rem;
        }
        
        .btn {
            padding: 0.5rem 1rem;
            border: none;
            border-radius: 4px;
            font-size: 0.875rem;
            cursor: pointer;
            transition: all 0.2s;
        }
        
        .btn-primary {
            background: #667eea;
            color: white;
        }
        
        .btn-primary:hover {
            background: #5a67d8;
        }
        
        .btn-secondary {
            background: #e2e8f0;
            color: #374151;
        }
        
        .btn-secondary:hover {
            background: #cbd5e0;
        }
        
        /* 响应式设计 */
        @media (max-width: 768px) {
            .sidebar {
                position: absolute;
                left: -300px;
                top: 0;
                height: 100%;
                z-index: 999;
                transition: left 0.3s;
            }
            
            .sidebar.open {
                left: 0;
            }
            
            .header-controls {
                gap: 0.5rem;
            }
            
            .view-toggle button {
                padding: 0.5rem;
                font-size: 0.75rem;
            }
        }
        
        /* 加载动画 */
        .loading {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 3px solid rgba(255,255,255,.3);
            border-radius: 50%;
            border-top-color: #fff;
            animation: spin 1s ease-in-out infinite;
        }
        
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="app-container">
        <!-- 顶部标题栏 -->
        <header class="header">
            <h1>
                <span>🤖</span>
                机器人路径编辑器
            </h1>
            <div class="header-controls">
                <div class="view-toggle">
                    <button id="canvas-view-btn" class="active">📊 画布视图</button>
                    <button id="table-view-btn">📋 表格视图</button>
                </div>
                <button class="btn btn-primary" id="save-btn">💾 保存</button>
            </div>
        </header>
        
        <div class="main-content">
            <!-- 左侧边栏 -->
            <aside class="sidebar">
                <div class="sidebar-header">
                    <h2>属性面板</h2>
                </div>
                <div class="sidebar-content" id="sidebar-content">
                    <!-- 节点属性面板 -->
                    <div class="property-panel" id="node-properties" style="display: none;">
                        <h3>📍 节点属性</h3>
                        <div class="property-content">
                            <div class="form-group">
                                <label>节点ID</label>
                                <input type="text" id="node-id" readonly>
                            </div>
                            <div class="form-group">
                                <label>节点名称</label>
                                <input type="text" id="node-name" placeholder="输入节点名称">
                            </div>
                            <div class="form-group">
                                <label>X坐标</label>
                                <input type="number" id="node-x" step="0.1">
                            </div>
                            <div class="form-group">
                                <label>Y坐标</label>
                                <input type="number" id="node-y" step="0.1">
                            </div>
                            <div class="form-group">
                                <label>节点类型</label>
                                <select id="node-type">
                                    <option value="normal">普通节点</option>
                                    <option value="start">起始节点</option>
                                    <option value="end">结束节点</option>
                                    <option value="waypoint">路径点</option>
                                </select>
                            </div>
                            <div class="btn-group">
                                <button class="btn btn-primary" id="update-node-btn">更新</button>
                                <button class="btn btn-secondary" id="delete-node-btn">删除</button>
                            </div>
                        </div>
                    </div>
                    
                    <!-- 路径属性面板 -->
                    <div class="property-panel" id="path-properties" style="display: none;">
                        <h3>🔗 路径属性</h3>
                        <div class="property-content">
                            <div class="form-group">
                                <label>路径ID</label>
                                <input type="text" id="path-id" readonly>
                            </div>
                            <div class="form-group">
                                <label>起始节点</label>
                                <input type="text" id="path-start" readonly>
                            </div>
                            <div class="form-group">
                                <label>目标节点</label>
                                <input type="text" id="path-end" readonly>
                            </div>
                            <div class="form-group">
                                <label>路径权重</label>
                                <input type="number" id="path-weight" step="0.1" value="1.0">
                            </div>
                            <div class="form-group">
                                <label>路径类型</label>
                                <select id="path-type">
                                    <option value="normal">普通路径</option>
                                    <option value="bidirectional">双向路径</option>
                                    <option value="one-way">单向路径</option>
                                </select>
                            </div>
                            <div class="btn-group">
                                <button class="btn btn-primary" id="update-path-btn">更新</button>
                                <button class="btn btn-secondary" id="delete-path-btn">删除</button>
                            </div>
                        </div>
                    </div>
                    
                    <!-- 默认提示 -->
                    <div class="property-panel" id="default-panel">
                        <h3>👈 选择元素</h3>
                        <div class="property-content">
                            <p style="color: #6c757d; font-size: 0.875rem;">
                                点击画布上的节点或路径来查看和编辑属性
                            </p>
                        </div>
                    </div>
                </div>
            </aside>
            
            <!-- 主工作区 -->
            <main class="work-area">
                <!-- 工具栏 -->
                <div class="toolbar">
                    <button id="add-node-btn">➕ 添加节点</button>
                    <button id="delete-mode-btn">🗑️ 删除模式</button>
                    <div class="separator"></div>
                    <button id="layout-force-btn">🔀 力导向布局</button>
                    <button id="layout-grid-btn">⚏ 网格布局</button>
                    <button id="layout-circle-btn">⭕ 圆形布局</button>
                    <div class="separator"></div>
                    <button id="generate-paths-btn">🔗 生成路径</button>
                    <button id="clear-paths-btn">🧹 清空路径</button>
                    <div class="separator"></div>
                    <button id="zoom-fit-btn">🔍 适应画布</button>
                    <button id="zoom-reset-btn">📐 重置缩放</button>
                </div>
                
                <!-- 画布容器 -->
                <div class="canvas-container" id="canvas-container">
                    <div id="canvas-stage"></div>
                </div>
                
                <!-- 表格容器 -->
                <div class="table-container" id="table-container">
                    <div class="toolbar">
                        <button id="add-row-btn">➕ 添加行</button>
                        <button id="delete-row-btn">🗑️ 删除行</button>
                        <div class="separator"></div>
                        <button id="import-csv-btn">📥 导入CSV</button>
                        <button id="export-csv-btn">📤 导出CSV</button>
                    </div>
                    <div class="table-wrapper">
                        <table class="data-table" id="data-table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>名称</th>
                                    <th>X坐标</th>
                                    <th>Y坐标</th>
                                    <th>类型</th>
                                    <th>操作</th>
                                </tr>
                            </thead>
                            <tbody id="table-body">
                                <!-- 表格数据将动态加载 -->
                            </tbody>
                        </table>
                    </div>
                </div>
            </main>
        </div>
        
        <!-- 状态栏 -->
        <footer class="status-bar">
            <div>
                <span id="node-count">节点: 0</span>
                <span style="margin-left: 1rem;" id="path-count">路径: 0</span>
            </div>
            <div>
                <span id="connection-status">已连接到数据库</span>
            </div>
        </footer>
    </div>
    
    <!-- 加载JavaScript文件 -->
    <script src="/static/app.js"></script>
    <script src="/static/canvas.js"></script>
    <script src="/static/table.js"></script>
    
    <script>
        // 应用初始化
        document.addEventListener('DOMContentLoaded', function() {
            console.log('🤖 机器人路径编辑器启动中...');
            
            // 初始化画布和表格视图
            if (typeof initCanvas === 'function') {
                initCanvas();
            }
            if (typeof initTable === 'function') {
                initTable();
            }
            
            // 视图切换
            const canvasViewBtn = document.getElementById('canvas-view-btn');
            const tableViewBtn = document.getElementById('table-view-btn');
            const canvasContainer = document.getElementById('canvas-container');
            const tableContainer = document.getElementById('table-container');
            
            canvasViewBtn.addEventListener('click', () => {
                canvasViewBtn.classList.add('active');
                tableViewBtn.classList.remove('active');
                canvasContainer.style.display = 'block';
                tableContainer.style.display = 'none';
            });
            
            tableViewBtn.addEventListener('click', () => {
                tableViewBtn.classList.add('active');
                canvasViewBtn.classList.remove('active');
                canvasContainer.style.display = 'none';
                tableContainer.style.display = 'flex';
            });
            
            console.log('✅ 应用界面初始化完成');
        });
    </script>
</body>
</html>`)

// ServeStatic 提供静态文件服务
func ServeStatic() http.FileSystem {
	return http.FS(StaticFiles)
}
