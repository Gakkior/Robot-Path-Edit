// Package web 鎻愪緵鍐呭祵鐨勫墠绔祫婧?
package web

import (
	"embed"
	"net/http"
)

//go:embed static/*
var StaticFiles embed.FS

// IndexHTML 涓婚〉闈TML
var IndexHTML = []byte(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>鏈哄櫒浜鸿矾寰勭紪杈戝櫒</title>
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
        <div class="logo">馃</div>
        <h1>鏈哄櫒浜鸿矾寰勭紪杈戝櫒</h1>
        <p class="subtitle">鐜颁唬鍖栫殑涓夌鍏煎璺緞绠＄悊宸ュ叿</p>
        
        <div class="status">鉁?鏈嶅姟杩愯涓?/div>
        
        <div class="features">
            <div class="feature">
                <h3>馃搳 鍙鍖栫紪杈?/h3>
                <p>鎷栨嫿寮忕敾甯冪紪杈戯紝鐩磋鐨勮妭鐐瑰拰璺緞绠＄悊</p>
            </div>
            <div class="feature">
                <h3>馃梽锔?鏁版嵁搴撴敮鎸?/h3>
                <p>鏀寔SQLite銆丮ySQL绛夊绉嶆暟鎹簱</p>
            </div>
            <div class="feature">
                <h3>馃摫 涓夌鍏煎</h3>
                <p>鏀寔Windows銆丩inux銆丄ndroid骞冲彴</p>
            </div>
            <div class="feature">
                <h3>馃帹 鐜颁唬璁捐</h3>
                <p>缇庤娴佺晠鐨勭幇浠ｅ寲鐣岄潰璁捐</p>
            </div>
        </div>
        
        <div class="api-info">
            <h3>API 绔偣</h3>
            <p>RESTful API 鏈嶅姟宸插惎鍔紝鍙闂互涓嬬鐐癸細</p>
            <div><code class="api-endpoint">GET /api/v1/nodes</code> - 鑾峰彇鑺傜偣鍒楄〃</div>
            <div><code class="api-endpoint">GET /api/v1/paths</code> - 鑾峰彇璺緞鍒楄〃</div>
            <div><code class="api-endpoint">GET /health</code> - 鍋ュ悍妫€鏌?/div>
            <div><code class="api-endpoint">GET /metrics</code> - 绯荤粺鎸囨爣</div>
        </div>
    </div>
    
    <script>
        // 绠€鍗曠殑API娴嬭瘯
        fetch('/health')
            .then(response => response.json())
            .then(data => {
                console.log('Health check:', data);
            })
            .catch(error => {
                console.error('Error:', error);
            });
            
        // 瀹炴椂鐘舵€佹洿鏂?
        setInterval(() => {
            fetch('/health')
                .then(response => {
                    if (response.ok) {
                        document.querySelector('.status').textContent = '鉁?鏈嶅姟杩愯涓?;
                        document.querySelector('.status').style.background = '#2ecc71';
                    } else {
                        document.querySelector('.status').textContent = '鉂?鏈嶅姟寮傚父';
                        document.querySelector('.status').style.background = '#e74c3c';
                    }
                })
                .catch(() => {
                    document.querySelector('.status').textContent = '鉂?杩炴帴澶辫触';
                    document.querySelector('.status').style.background = '#e74c3c';
                });
        }, 10000); // 姣?0绉掓鏌ヤ竴娆?
    </script>
</body>
</html>`)

// ServeStatic 鎻愪緵闈欐€佹枃浠舵湇鍔?
func ServeStatic() http.FileSystem {
	return http.FS(StaticFiles)
}
