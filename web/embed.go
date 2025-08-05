// Package web 提供内嵌的前端资�?
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
        
        <div class="status">�?服务运行�?/div>
        
        <div class="features">
            <div class="feature">
                <h3>📊 可视化编�?/h3>
                <p>拖拽式画布编辑，直观的节点和路径管理</p>
            </div>
            <div class="feature">
                <h3>🗄�?数据库支�?/h3>
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
            <h3>API 端点</h3>
            <p>RESTful API 服务已启动，可访问以下端点：</p>
            <div><code class="api-endpoint">GET /api/v1/nodes</code> - 获取节点列表</div>
            <div><code class="api-endpoint">GET /api/v1/paths</code> - 获取路径列表</div>
            <div><code class="api-endpoint">GET /health</code> - 健康检�?/div>
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
            
        // 实时状态更�?
        setInterval(() => {
            fetch('/health')
                .then(response => {
                    if (response.ok) {
                        document.querySelector('.status').textContent = '�?服务运行�?;
                        document.querySelector('.status').style.background = '#2ecc71';
                    } else {
                        document.querySelector('.status').textContent = '�?服务异常';
                        document.querySelector('.status').style.background = '#e74c3c';
                    }
                })
                .catch(() => {
                    document.querySelector('.status').textContent = '�?连接失败';
                    document.querySelector('.status').style.background = '#e74c3c';
                });
        }, 10000); // �?0秒检查一�?
    </script>
</body>
</html>`)

// ServeStatic 提供静态文件服�?
func ServeStatic() http.FileSystem {
	return http.FS(StaticFiles)
}
