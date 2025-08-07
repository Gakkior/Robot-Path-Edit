@echo off
chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

REM 机器人路径编辑器 Windows构建脚本
set PROJECT_NAME=robot-path-editor
set VERSION=%VERSION%
if "%VERSION%"=="" set VERSION=v1.0.0

echo [信息] 机器人路径编辑器构建脚本
echo [信息] 项目版本: %VERSION%

REM 检查Go环境
go version >nul 2>&1
if errorlevel 1 (
    echo [错误] Go环境未安装或不在PATH中
    exit /b 1
)

REM 创建构建目录
if not exist "build" mkdir "build"

REM 解析命令行参数
set COMMAND=%1
if "%COMMAND%"=="" set COMMAND=build

if /i "%COMMAND%"=="clean" goto :clean
if /i "%COMMAND%"=="test" goto :test
if /i "%COMMAND%"=="build" goto :build
if /i "%COMMAND%"=="help" goto :help

echo [错误] 未知命令: %COMMAND%
goto :help

:clean
echo [信息] 清理构建目录...
if exist "build" rmdir /s /q "build"
echo [成功] 清理完成
goto :end

:test
echo [信息] 运行测试...
go test ./... -v
if errorlevel 1 (
    echo [错误] 测试失败
    exit /b 1
)
echo [成功] 测试通过
goto :end

:build
echo [信息] 构建应用程序...
set CGO_ENABLED=0
go build -ldflags "-s -w" -o "build\%PROJECT_NAME%.exe" cmd/server/main.go
go build -ldflags "-s -w" -o "build\%PROJECT_NAME%-demo.exe" cmd/demo/main.go

if exist "build\%PROJECT_NAME%.exe" (
    echo [成功] 构建完成: build\%PROJECT_NAME%.exe
    echo [成功] 演示版本: build\%PROJECT_NAME%-demo.exe
) else (
    echo [错误] 构建失败
    exit /b 1
)
goto :end

:help
echo.
echo 用法: %0 [命令]
echo.
echo 命令:
echo   build        构建应用程序
echo   test         运行测试
echo   clean        清理构建文件
echo   help         显示此帮助信息
echo.
goto :end

:end
endlocal