#!/bin/bash
set -e

echo "构建前端项目..."

cd frontend

echo "安装依赖..."
npm install

echo "构建生产版本..."
npm run build

echo "复制构建文件到web目录..."
rm -rf ../web/static/new-frontend
mkdir -p ../web/static/new-frontend
cp -r dist/* ../web/static/new-frontend/

echo "前端构建完成！"
echo "访问: http://localhost:8080/app/new"