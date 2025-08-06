@echo off
chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

REM 机器人路径编辑器 Windows构建脚本

set PROJECT_NAME=robot-path-editor
set VERSION=%VERSION%
if "%VERSION%"=="" set VERSION=v1.0.0

set BUILD_TIME=%date% %time%
set BUILD_DIR=build
set DIST_DIR=dist

REM 颜色定义
for /F %%a in ('echo prompt $E ^| cmd') do set "ESC=%%a"
set "RED=%ESC%[31m"
set "GREEN=%ESC%[32m"
set "YELLOW=%ESC%[33m"
set "BLUE=%ESC%[34m"
set "NC=%ESC%[0m"

echo %BLUE%[信息]%NC% 机器人路径编辑器构建脚本
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
if not exist "%DIST_DIR%" mkdir "%DIST_DIR%"

REM 解析命令行参数
set COMMAND=%1
if "%COMMAND%"=="" set COMMAND=build

if /i "%COMMAND%"=="clean" goto :clean
if /i "%COMMAND%"=="test" goto :test
if /i "%COMMAND%"=="build" goto :build
if /i "%COMMAND%"=="build-all" goto :build_all
if /i "%COMMAND%"=="release" goto :release
if /i "%COMMAND%"=="help" goto :help

echo %RED%[错误]%NC% 未知命令: %COMMAND%
goto :help

:clean
echo %BLUE%[信息]%NC% 清理构建目录...
if exist "%BUILD_DIR%" rmdir /s /q "%BUILD_DIR%"
if exist "%DIST_DIR%" rmdir /s /q "%DIST_DIR%"
echo %GREEN%[成功]%NC% 清理完成
goto :end

:test
echo %BLUE%[信息]%NC% 运行测试...
if exist "tests" (
    go test ./tests/... -v
    if errorlevel 1 (
        echo %RED%[错误]%NC% 测试失败
        exit /b 1
    )
) else (
    echo %YELLOW%[警告]%NC% 未找到测试目录，跳过测试
)
echo %GREEN%[成功]%NC% 测试通过
goto :end

:build
echo %BLUE%[信息]%NC% 为当前平台构建...
set CGO_ENABLED=0
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%.exe" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo.exe" cmd/demo/main.go

if exist "%BUILD_DIR%\%PROJECT_NAME%.exe" (
    echo %GREEN%[成功]%NC% 构建完成: %BUILD_DIR%\%PROJECT_NAME%.exe
) else (
    echo %RED%[错误]%NC% 构建失败
    exit /b 1
)
goto :end

:build_all
echo %BLUE%[信息]%NC% 开始跨平台构建...

REM Windows
echo %BLUE%[信息]%NC% 构建 Windows/amd64...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-windows-amd64.exe" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo-windows-amd64.exe" cmd/demo/main.go

echo %BLUE%[信息]%NC% 构建 Windows/386...
set GOARCH=386
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-windows-386.exe" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo-windows-386.exe" cmd/demo/main.go

REM Linux
echo %BLUE%[信息]%NC% 构建 Linux/amd64...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-linux-amd64" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo-linux-amd64" cmd/demo/main.go

echo %BLUE%[信息]%NC% 构建 Linux/arm64...
set GOARCH=arm64
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-linux-arm64" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo-linux-arm64" cmd/demo/main.go

REM macOS
echo %BLUE%[信息]%NC% 构建 macOS/amd64...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-darwin-amd64" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo-darwin-amd64" cmd/demo/main.go

echo %BLUE%[信息]%NC% 构建 macOS/arm64...
set GOARCH=arm64
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-darwin-arm64" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo-darwin-arm64" cmd/demo/main.go

echo %GREEN%[成功]%NC% 跨平台构建完成
goto :end

:release
echo %BLUE%[信息]%NC% 创建发布包...
call :test
if errorlevel 1 exit /b 1

call :build_all
if errorlevel 1 exit /b 1

REM 复制资源文件
if exist "web" xcopy /e /i "web" "%BUILD_DIR%\web"
if exist "configs\config.yaml" copy "configs\config.yaml" "%BUILD_DIR%\"

REM 创建发布包
for %%f in ("%BUILD_DIR%\%PROJECT_NAME%-*.exe") do (
    set "filename=%%~nf"
    set "platform=!filename:%PROJECT_NAME%-=!"
    
    mkdir "%BUILD_DIR%\temp-!platform!"
    copy "%%f" "%BUILD_DIR%\temp-!platform!\%PROJECT_NAME%.exe"
    if exist "%BUILD_DIR%\%PROJECT_NAME%-demo-!platform!.exe" copy "%BUILD_DIR%\%PROJECT_NAME%-demo-!platform!.exe" "%BUILD_DIR%\temp-!platform!\%PROJECT_NAME%-demo.exe"
    if exist "%BUILD_DIR%\web" xcopy /e /i "%BUILD_DIR%\web" "%BUILD_DIR%\temp-!platform!\web"
    if exist "%BUILD_DIR%\config.yaml" copy "%BUILD_DIR%\config.yaml" "%BUILD_DIR%\temp-!platform!\"
    
    echo 机器人路径编辑器 %VERSION% > "%BUILD_DIR%\temp-!platform!\README.txt"
    echo. >> "%BUILD_DIR%\temp-!platform!\README.txt"
    echo 构建时间: %BUILD_TIME% >> "%BUILD_DIR%\temp-!platform!\README.txt"
    echo. >> "%BUILD_DIR%\temp-!platform!\README.txt"
    echo 使用方法: >> "%BUILD_DIR%\temp-!platform!\README.txt"
    echo 1. 运行 %PROJECT_NAME%.exe 启动服务器 >> "%BUILD_DIR%\temp-!platform!\README.txt"
    echo 2. 运行 %PROJECT_NAME%-demo.exe 启动演示版 >> "%BUILD_DIR%\temp-!platform!\README.txt"
    echo 3. 访问 http://localhost:8080 >> "%BUILD_DIR%\temp-!platform!\README.txt"
    
    if exist "C:\Program Files\7-Zip\7z.exe" (
        "C:\Program Files\7-Zip\7z.exe" a "%DIST_DIR%\%PROJECT_NAME%-%VERSION%-!platform!.zip" "%BUILD_DIR%\temp-!platform!\*"
    ) else (
        powershell -command "Compress-Archive -Path '%BUILD_DIR%\temp-!platform!\*' -DestinationPath '%DIST_DIR%\%PROJECT_NAME%-%VERSION%-!platform!.zip'"
    )
    
    rmdir /s /q "%BUILD_DIR%\temp-!platform!"
    echo %GREEN%[成功]%NC% 发布包创建完成: %DIST_DIR%\%PROJECT_NAME%-%VERSION%-!platform!.zip
)

echo %GREEN%[成功]%NC% 所有发布包创建完成
goto :end

:help
echo.
echo 机器人路径编辑器 Windows构建脚本
echo.
echo 用法: %0 [命令]
echo.
echo 命令:
echo   build        为当前平台构建二进制文件
echo   build-all    为所有平台构建二进制文件
echo   test         运行测试
echo   release      创建发布包
echo   clean        清理构建文件
echo   help         显示此帮助信息
echo.
echo 示例:
echo   %0 build-all    # 为所有平台构建
echo   %0 test         # 运行测试
echo   %0 release      # 创建发布包
echo.
echo 环境变量:
echo   VERSION         设置版本号 (默认: v1.0.0)
echo.
goto :end

:end
endlocal 