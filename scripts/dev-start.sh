#!/bin/bash
set -e

echo "========================================"
echo "  机器人路径编辑器 - 开发环境启动"
echo "========================================"
echo

echo "[1/3] 检查Node.js环境..."
if ! command -v node &> /dev/null; then
    echo "❌ Node.js未安装，请先安装Node.js 18+"
    exit 1
fi
echo "✅ Node.js环境正常: $(node --version)"

echo
echo "[2/3] 安装前端依赖..."
cd frontend
if [ ! -d "node_modules" ]; then
    echo "首次运行，正在安装依赖..."
    npm install
else
    echo "✅ 依赖已存在"
fi

echo
echo "[3/3] 启动开发服务器..."
echo
echo "🚀 前端开发服务器: http://localhost:5173"
echo "🔗 后端API服务器: http://localhost:8080"
echo
echo "提示: 请在另一个终端运行 'go run cmd/server/main.go' 启动后端"
echo

# 在后台启动后端
(cd .. && go run cmd/server/main.go) &

# 启动前端
npm run dev