@echo off
echo 构建前端项目...

cd frontend

echo 安装依赖...
call npm install

echo 构建生产版本...
call npm run build

echo 复制构建文件到web目录...
if exist "..\web\static\new-frontend" rmdir /s /q "..\web\static\new-frontend"
mkdir "..\web\static\new-frontend"
xcopy /e /y "dist\*" "..\web\static\new-frontend\"

echo 前端构建完成！
echo 访问: http://localhost:8080/app/new
pause