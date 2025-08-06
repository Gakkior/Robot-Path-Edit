@echo off
setlocal enabledelayedexpansion

REM Robot Path Editor Windows Build Script

set PROJECT_NAME=robot-path-editor
set VERSION=%VERSION%
if "%VERSION%"=="" set VERSION=v1.0.0

set BUILD_TIME=%date% %time%
set BUILD_DIR=build
set DIST_DIR=dist

REM Color definitions
for /F %%a in ('echo prompt $E ^| cmd') do set "ESC=%%a"
set "RED=%ESC%[31m"
set "GREEN=%ESC%[32m"
set "YELLOW=%ESC%[33m"
set "BLUE=%ESC%[34m"
set "NC=%ESC%[0m"

echo %BLUE%[INFO]%NC% Robot Path Editor Build Script
echo %BLUE%[INFO]%NC% Project Version: %VERSION%
echo %BLUE%[INFO]%NC% Build Time: %BUILD_TIME%

REM Check Go environment
go version >nul 2>&1
if errorlevel 1 (
    echo %RED%[ERROR]%NC% Go environment not installed or not in PATH
    exit /b 1
)

REM Create build directories
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"
if not exist "%DIST_DIR%" mkdir "%DIST_DIR%"

REM Parse command line arguments
set COMMAND=%1
if "%COMMAND%"=="" set COMMAND=build

if /i "%COMMAND%"=="clean" goto :clean
if /i "%COMMAND%"=="test" goto :test
if /i "%COMMAND%"=="build" goto :build
if /i "%COMMAND%"=="build-all" goto :build_all
if /i "%COMMAND%"=="release" goto :release
if /i "%COMMAND%"=="help" goto :help

echo %RED%[ERROR]%NC% Unknown command: %COMMAND%
goto :help

:clean
echo %BLUE%[INFO]%NC% Cleaning build directories...
if exist "%BUILD_DIR%" rmdir /s /q "%BUILD_DIR%"
if exist "%DIST_DIR%" rmdir /s /q "%DIST_DIR%"
echo %GREEN%[SUCCESS]%NC% Clean completed
goto :end

:test
echo %BLUE%[INFO]%NC% Running tests...
if exist "tests" (
    go test ./tests/... -v
    if errorlevel 1 (
        echo %RED%[ERROR]%NC% Tests failed
        exit /b 1
    )
) else (
    echo %YELLOW%[WARNING]%NC% No tests directory found, skipping tests
)
echo %GREEN%[SUCCESS]%NC% Tests passed
goto :end

:build
echo %BLUE%[INFO]%NC% Building for current platform...
set CGO_ENABLED=0
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%.exe" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo.exe" cmd/demo/main.go

if exist "%BUILD_DIR%\%PROJECT_NAME%.exe" (
    echo %GREEN%[SUCCESS]%NC% Build completed: %BUILD_DIR%\%PROJECT_NAME%.exe
) else (
    echo %RED%[ERROR]%NC% Build failed
    exit /b 1
)
goto :end

:build_all
echo %BLUE%[INFO]%NC% Starting cross-platform build...

REM Windows
echo %BLUE%[INFO]%NC% Building Windows/amd64...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-windows-amd64.exe" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo-windows-amd64.exe" cmd/demo/main.go

echo %BLUE%[INFO]%NC% Building Windows/386...
set GOARCH=386
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-windows-386.exe" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo-windows-386.exe" cmd/demo/main.go

REM Linux
echo %BLUE%[INFO]%NC% Building Linux/amd64...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-linux-amd64" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo-linux-amd64" cmd/demo/main.go

echo %BLUE%[INFO]%NC% Building Linux/arm64...
set GOARCH=arm64
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-linux-arm64" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo-linux-arm64" cmd/demo/main.go

REM macOS
echo %BLUE%[INFO]%NC% Building macOS/amd64...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-darwin-amd64" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo-darwin-amd64" cmd/demo/main.go

echo %BLUE%[INFO]%NC% Building macOS/arm64...
set GOARCH=arm64
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-darwin-arm64" cmd/server/main.go
go build -ldflags "-s -w" -o "%BUILD_DIR%\%PROJECT_NAME%-demo-darwin-arm64" cmd/demo/main.go

echo %GREEN%[SUCCESS]%NC% Cross-platform build completed
goto :end

:release
echo %BLUE%[INFO]%NC% Creating release packages...
call :test
if errorlevel 1 exit /b 1

call :build_all
if errorlevel 1 exit /b 1

REM Copy resource files
if exist "web" xcopy /e /i "web" "%BUILD_DIR%\web"
if exist "configs\config.yaml" copy "configs\config.yaml" "%BUILD_DIR%\"

REM Create release packages
for %%f in ("%BUILD_DIR%\%PROJECT_NAME%-*.exe") do (
    set "filename=%%~nf"
    set "platform=!filename:%PROJECT_NAME%-=!"
    
    mkdir "%BUILD_DIR%\temp-!platform!"
    copy "%%f" "%BUILD_DIR%\temp-!platform!\%PROJECT_NAME%.exe"
    if exist "%BUILD_DIR%\%PROJECT_NAME%-demo-!platform!.exe" copy "%BUILD_DIR%\%PROJECT_NAME%-demo-!platform!.exe" "%BUILD_DIR%\temp-!platform!\%PROJECT_NAME%-demo.exe"
    if exist "%BUILD_DIR%\web" xcopy /e /i "%BUILD_DIR%\web" "%BUILD_DIR%\temp-!platform!\web"
    if exist "%BUILD_DIR%\config.yaml" copy "%BUILD_DIR%\config.yaml" "%BUILD_DIR%\temp-!platform!\"
    
    echo %PROJECT_NAME% %VERSION% > "%BUILD_DIR%\temp-!platform!\README.txt"
    echo. >> "%BUILD_DIR%\temp-!platform!\README.txt"
    echo Build Time: %BUILD_TIME% >> "%BUILD_DIR%\temp-!platform!\README.txt"
    echo. >> "%BUILD_DIR%\temp-!platform!\README.txt"
    echo Usage: >> "%BUILD_DIR%\temp-!platform!\README.txt"
    echo 1. Run %PROJECT_NAME%.exe to start server >> "%BUILD_DIR%\temp-!platform!\README.txt"
    echo 2. Run %PROJECT_NAME%-demo.exe to start demo >> "%BUILD_DIR%\temp-!platform!\README.txt"
    echo 3. Visit http://localhost:8080 >> "%BUILD_DIR%\temp-!platform!\README.txt"
    
    if exist "C:\Program Files\7-Zip\7z.exe" (
        "C:\Program Files\7-Zip\7z.exe" a "%DIST_DIR%\%PROJECT_NAME%-%VERSION%-!platform!.zip" "%BUILD_DIR%\temp-!platform!\*"
    ) else (
        powershell -command "Compress-Archive -Path '%BUILD_DIR%\temp-!platform!\*' -DestinationPath '%DIST_DIR%\%PROJECT_NAME%-%VERSION%-!platform!.zip'"
    )
    
    rmdir /s /q "%BUILD_DIR%\temp-!platform!"
    echo %GREEN%[SUCCESS]%NC% Release package created: %DIST_DIR%\%PROJECT_NAME%-%VERSION%-!platform!.zip
)

echo %GREEN%[SUCCESS]%NC% Release packages created
goto :end

:help
echo.
echo Robot Path Editor Windows Build Script
echo.
echo Usage: %0 [command]
echo.
echo Commands:
echo   build        Build binary for current platform
echo   build-all    Build binaries for all platforms
echo   test         Run tests
echo   release      Create release packages
echo   clean        Clean build files
echo   help         Show this help information
echo.
echo Examples:
echo   %0 build-all    # Build for all platforms
echo   %0 test         # Run tests
echo   %0 release      # Create release packages
echo.
echo Environment Variables:
echo   VERSION         Set version number (default: v1.0.0)
echo.
goto :end

:end
endlocal 