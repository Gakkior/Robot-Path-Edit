@echo off
setlocal enabledelayedexpansion

REM Robot Path Editor Windows Build Script (English Version)

set PROJECT_NAME=robot-path-editor
set VERSION=%VERSION%
if "%VERSION%"=="" set VERSION=v1.0.0

set BUILD_TIME=%date% %time%
set BUILD_DIR=build

REM Color definitions
for /F %%a in ('echo prompt $E ^| cmd') do set "ESC=%%a"
set "RED=%ESC%[31m"
set "GREEN=%ESC%[32m"
set "YELLOW=%ESC%[33m"
set "BLUE=%ESC%[34m"
set "NC=%ESC%[0m"

echo %BLUE%[INFO]%NC% Robot Path Editor Windows Build Script
echo %BLUE%[INFO]%NC% Project Version: %VERSION%
echo %BLUE%[INFO]%NC% Build Time: %BUILD_TIME%

REM Check Go environment
go version >nul 2>&1
if errorlevel 1 (
    echo %RED%[ERROR]%NC% Go environment not installed or not in PATH
    exit /b 1
)

REM Create build directory
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

echo %BLUE%[INFO]%NC% Starting Windows build...

REM Build current architecture (usually 64-bit)
echo %BLUE%[INFO]%NC% Building current Windows architecture...
set CGO_ENABLED=0
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%.exe" cmd/server/main.go

REM Build Windows 64-bit
echo %BLUE%[INFO]%NC% Building Windows 64-bit...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-win64.exe" cmd/server/main.go

REM Build Windows 32-bit
echo %BLUE%[INFO]%NC% Building Windows 32-bit...
set GOARCH=386
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-win32.exe" cmd/server/main.go

REM Check build results
if exist "%BUILD_DIR%\%PROJECT_NAME%.exe" (
    echo %GREEN%[SUCCESS]%NC% Current architecture build completed: %BUILD_DIR%\%PROJECT_NAME%.exe
) else (
    echo %RED%[ERROR]%NC% Current architecture build failed
)

if exist "%BUILD_DIR%\%PROJECT_NAME%-win64.exe" (
    echo %GREEN%[SUCCESS]%NC% Windows 64-bit build completed: %BUILD_DIR%\%PROJECT_NAME%-win64.exe
) else (
    echo %RED%[ERROR]%NC% Windows 64-bit build failed
)

if exist "%BUILD_DIR%\%PROJECT_NAME%-win32.exe" (
    echo %GREEN%[SUCCESS]%NC% Windows 32-bit build completed: %BUILD_DIR%\%PROJECT_NAME%-win32.exe
) else (
    echo %RED%[ERROR]%NC% Windows 32-bit build failed
)

echo.
echo %GREEN%[COMPLETE]%NC% Windows build completed!
echo.
echo Build files location:
if exist "%BUILD_DIR%\%PROJECT_NAME%.exe" echo   - %BUILD_DIR%\%PROJECT_NAME%.exe (current architecture)
if exist "%BUILD_DIR%\%PROJECT_NAME%-win64.exe" echo   - %BUILD_DIR%\%PROJECT_NAME%-win64.exe (64-bit)
if exist "%BUILD_DIR%\%PROJECT_NAME%-win32.exe" echo   - %BUILD_DIR%\%PROJECT_NAME%-win32.exe (32-bit)
echo.
echo Usage:
echo   %BUILD_DIR%\%PROJECT_NAME%.exe
echo   or
echo   %BUILD_DIR%\%PROJECT_NAME%-win64.exe
echo.
echo %YELLOW%[TIP]%NC% For cross-platform builds, use: .\build.bat build-all

endlocal