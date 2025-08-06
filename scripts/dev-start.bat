@echo off
echo ========================================
echo   机器人路径编辑器 - 开发环境启动
echo ========================================
echo.

echo [1/3] 检查Node.js环境...
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Node.js未安装，请先安装Node.js 18+
    pause
    exit /b 1
)
echo ✅ Node.js环境正常

echo.
echo [2/3] 安装前端依赖...
cd frontend
if not exist node_modules (
    echo 首次运行，正在安装依赖...
    call npm install
    if %errorlevel% neq 0 (
        echo ❌ 依赖安装失败
        pause
        exit /b 1
    )
) else (
    echo ✅ 依赖已存在
)

echo.
echo [3/3] 启动开发服务器...
echo.
echo 🚀 前端开发服务器: http://localhost:5173
echo 🔗 后端API服务器: http://localhost:8080
echo.
echo 提示: 请在另一个终端运行 'go run cmd/server/main.go' 启动后端
echo.

start cmd /k "cd .. && echo 启动后端服务器... && go run cmd/server/main.go"
call npm run dev

pause