@echo off
chcp 65001 >nul 2>&1

echo ========================================
echo   机器人路径编辑器 - 开发环境启动
echo ========================================
echo.

echo [1/2] 检查Node.js环境...
node --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Node.js未安装，请先安装Node.js 18+
    pause
    exit /b 1
)
for /f %%i in ('node --version') do echo ✅ Node.js环境正常: %%i

echo.
echo [2/2] 启动前端开发服务器...
cd frontend

if not exist "node_modules" (
    echo 首次运行，正在安装依赖...
    npm install
    if errorlevel 1 (
        echo ❌ 依赖安装失败
        pause
        exit /b 1
    )
)

echo.
echo 🚀 前端开发服务器: http://localhost:5173
echo 🔗 后端API服务器: http://localhost:8080
echo.
echo 💡 提示: 请在另一个终端运行 'go run cmd/server/main.go' 启动后端
echo.

npm run dev
pause