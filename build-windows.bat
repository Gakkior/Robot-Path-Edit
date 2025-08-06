@echo off
chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

REM 机器人路径编辑器 Windows专用构建脚本

set PROJECT_NAME=robot-path-editor
set VERSION=%VERSION%
if "%VERSION%"=="" set VERSION=v1.0.0

set BUILD_TIME=%date% %time%
set BUILD_DIR=build

REM 颜色定义
for /F %%a in ('echo prompt $E ^| cmd') do set "ESC=%%a"
set "RED=%ESC%[31m"
set "GREEN=%ESC%[32m"
set "YELLOW=%ESC%[33m"
set "BLUE=%ESC%[34m"
set "NC=%ESC%[0m"

echo %BLUE%[信息]%NC% 机器人路径编辑器 Windows专用构建脚本
echo %BLUE%[信息]%NC% 项目版本: %VERSION%
echo %BLUE%[信息]%NC% 构建时间: %BUILD_TIME%

REM 检查Go环境
go version >nul 2>&1
if errorlevel 1 (
    echo %RED%[错误]%NC% Go环境未安装或不在PATH中
    exit /b 1
)

REM 创建构建目录
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

echo %BLUE%[信息]%NC% 开始编译Windows版本...

REM 编译当前架构 (通常是64位)
echo %BLUE%[信息]%NC% 编译当前Windows架构...
set CGO_ENABLED=0
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%.exe" cmd/server/main.go

REM 编译Windows 64位
echo %BLUE%[信息]%NC% 编译Windows 64位版本...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-win64.exe" cmd/server/main.go

REM 编译Windows 32位
echo %BLUE%[信息]%NC% 编译Windows 32位版本...
set GOARCH=386
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-win32.exe" cmd/server/main.go

REM 检查编译结果
if exist "%BUILD_DIR%\%PROJECT_NAME%.exe" (
    echo %GREEN%[成功]%NC% 当前架构编译完成: %BUILD_DIR%\%PROJECT_NAME%.exe
) else (
    echo %RED%[错误]%NC% 当前架构编译失败
)

if exist "%BUILD_DIR%\%PROJECT_NAME%-win64.exe" (
    echo %GREEN%[成功]%NC% Windows 64位编译完成: %BUILD_DIR%\%PROJECT_NAME%-win64.exe
) else (
    echo %RED%[错误]%NC% Windows 64位编译失败
)

if exist "%BUILD_DIR%\%PROJECT_NAME%-win32.exe" (
    echo %GREEN%[成功]%NC% Windows 32位编译完成: %BUILD_DIR%\%PROJECT_NAME%-win32.exe
) else (
    echo %RED%[错误]%NC% Windows 32位编译失败
)

echo.
echo %GREEN%[完成]%NC% Windows专用编译完成！
echo.
echo 编译文件位置:
if exist "%BUILD_DIR%\%PROJECT_NAME%.exe" echo   - %BUILD_DIR%\%PROJECT_NAME%.exe (当前架构)
if exist "%BUILD_DIR%\%PROJECT_NAME%-win64.exe" echo   - %BUILD_DIR%\%PROJECT_NAME%-win64.exe (64位)
if exist "%BUILD_DIR%\%PROJECT_NAME%-win32.exe" echo   - %BUILD_DIR%\%PROJECT_NAME%-win32.exe (32位)
echo.
echo 使用方法:
echo   %BUILD_DIR%\%PROJECT_NAME%.exe
echo   或者
echo   %BUILD_DIR%\%PROJECT_NAME%-win64.exe
echo.
echo %YELLOW%[提示]%NC% 如需其他平台编译，请使用: .\build.bat build-all

endlocal