#!/bin/bash

# 机器人路径编辑器构建脚本
# 支持跨平台构建和打包

set -e

# 项目配置
PROJECT_NAME="robot-path-editor"
VERSION="${VERSION:-v1.0.0}"
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DIR="build"
DIST_DIR="dist"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印信息
info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 清理函数
cleanup() {
    info "清理构建目录..."
    rm -rf ${BUILD_DIR}
    rm -rf ${DIST_DIR}
}

# 准备构建环境
prepare() {
    info "准备构建环境..."
    mkdir -p ${BUILD_DIR}
    mkdir -p ${DIST_DIR}
    
    # 检查Go环境
    if ! command -v go &> /dev/null; then
        error "Go环境未安装或未在PATH中"
        exit 1
    fi
    
    info "Go版本: $(go version)"
    info "项目版本: ${VERSION}"
    info "构建时间: ${BUILD_TIME}"
    info "Git提交: ${GIT_COMMIT}"
}

# 运行测试
test() {
    info "运行测试..."
    
    # 单元测试
    if [ -d "tests/unit" ]; then
        info "运行单元测试..."
        go test ./tests/unit/... -v -cover -coverprofile=${BUILD_DIR}/unit-coverage.out
    fi
    
    # 集成测试
    if [ -d "tests/integration" ]; then
        info "运行集成测试..."
        go test ./tests/integration/... -v -cover -coverprofile=${BUILD_DIR}/integration-coverage.out
    fi
    
    # 生成覆盖率报告
    if [ -f "${BUILD_DIR}/unit-coverage.out" ] || [ -f "${BUILD_DIR}/integration-coverage.out" ]; then
        info "生成覆盖率报告..."
        go tool cover -html=${BUILD_DIR}/unit-coverage.out -o ${BUILD_DIR}/coverage.html 2>/dev/null || true
    fi
    
    success "测试完成"
}

# 构建二进制文件
build() {
    local target_os=$1
    local target_arch=$2
    local output_name=$3
    
    info "构建 ${target_os}/${target_arch}..."
    
    # 设置构建标志
    LDFLAGS="-s -w"
    LDFLAGS="${LDFLAGS} -X 'main.Version=${VERSION}'"
    LDFLAGS="${LDFLAGS} -X 'main.BuildTime=${BUILD_TIME}'"
    LDFLAGS="${LDFLAGS} -X 'main.GitCommit=${GIT_COMMIT}'"
    
    # 设置环境变量
    export CGO_ENABLED=0
    export GOOS=${target_os}
    export GOARCH=${target_arch}
    
    # 构建主服务器
    go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${output_name} cmd/server/main.go
    
    # 构建演示版
    go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${output_name}-demo cmd/demo/main.go
    
    success "构建完成: ${BUILD_DIR}/${output_name}"
}

# 构建所有平台
build_all() {
    info "开始跨平台构建..."
    
    # Windows
    build "windows" "amd64" "${PROJECT_NAME}-windows-amd64.exe"
    build "windows" "386" "${PROJECT_NAME}-windows-386.exe"
    
    # Linux
    build "linux" "amd64" "${PROJECT_NAME}-linux-amd64"
    build "linux" "386" "${PROJECT_NAME}-linux-386"
    build "linux" "arm64" "${PROJECT_NAME}-linux-arm64"
    
    # macOS
    build "darwin" "amd64" "${PROJECT_NAME}-darwin-amd64"
    build "darwin" "arm64" "${PROJECT_NAME}-darwin-arm64"
    
    success "所有平台构建完成"
}

# 打包静态资源
package_assets() {
    info "打包静态资源..."
    
    # 确保web目录存在
    if [ ! -d "web/static" ]; then
        warning "web/static目录不存在，跳过资源打包"
        return
    fi
    
    # 复制web资源
    cp -r web ${BUILD_DIR}/
    
    # 复制配置文件
    if [ -f "configs/config.yaml" ]; then
        cp configs/config.yaml ${BUILD_DIR}/
    fi
    
    success "资源打包完成"
}

# 创建发布包
create_release() {
    info "创建发布包..."
    
    # 打包静态资源
    package_assets
    
    # 创建各平台的发布包
    for binary in ${BUILD_DIR}/${PROJECT_NAME}-*; do
        if [ -f "$binary" ]; then
            platform=$(basename "$binary" | sed "s/${PROJECT_NAME}-//")
            package_name="${PROJECT_NAME}-${VERSION}-${platform}"
            
            info "创建发布包: ${package_name}"
            
            # 创建临时目录
            temp_dir="${BUILD_DIR}/temp-${platform}"
            mkdir -p "$temp_dir"
            
            # 复制文件
            cp "$binary" "$temp_dir/${PROJECT_NAME}"
            cp "${BUILD_DIR}/${PROJECT_NAME}-demo"* "$temp_dir/" 2>/dev/null || true
            
            # 复制资源文件
            if [ -d "${BUILD_DIR}/web" ]; then
                cp -r "${BUILD_DIR}/web" "$temp_dir/"
            fi
            if [ -f "${BUILD_DIR}/config.yaml" ]; then
                cp "${BUILD_DIR}/config.yaml" "$temp_dir/"
            fi
            
            # 创建README
            cat > "$temp_dir/README.txt" << EOF
${PROJECT_NAME} ${VERSION}

构建时间: ${BUILD_TIME}
Git提交: ${GIT_COMMIT}

使用方法:
1. 配置 config.yaml 文件
2. 运行 ./${PROJECT_NAME} 启动服务
3. 运行 ./${PROJECT_NAME}-demo 启动演示版
4. 访问 http://localhost:8080

文档: https://github.com/your-org/${PROJECT_NAME}
EOF
            
            # 创建压缩包
            if command -v zip &> /dev/null; then
                (cd "${BUILD_DIR}" && zip -r "../${DIST_DIR}/${package_name}.zip" "temp-${platform}/")
            else
                tar -czf "${DIST_DIR}/${package_name}.tar.gz" -C "${BUILD_DIR}" "temp-${platform}/"
            fi
            
            # 清理临时目录
            rm -rf "$temp_dir"
            
            success "发布包创建完成: ${DIST_DIR}/${package_name}"
        fi
    done
}

# Docker构建
build_docker() {
    info "构建Docker镜像..."
    
    if ! command -v docker &> /dev/null; then
        warning "Docker未安装，跳过Docker构建"
        return
    fi
    
    # 创建Dockerfile
    cat > Dockerfile << EOF
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o robot-path-editor cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/robot-path-editor .
COPY --from=builder /app/web ./web
COPY --from=builder /app/configs ./configs
EXPOSE 8080
CMD ["./robot-path-editor"]
EOF
    
    # 构建镜像
    docker build -t "${PROJECT_NAME}:${VERSION}" .
    docker build -t "${PROJECT_NAME}:latest" .
    
    # 清理Dockerfile
    rm -f Dockerfile
    
    success "Docker镜像构建完成"
}

# 性能基准测试
benchmark() {
    info "运行性能基准测试..."
    
    # 运行基准测试
    go test -bench=. -benchmem ./tests/unit/... > ${BUILD_DIR}/benchmark.txt 2>&1 || true
    
    if [ -f "${BUILD_DIR}/benchmark.txt" ]; then
        success "基准测试完成，结果保存到: ${BUILD_DIR}/benchmark.txt"
    fi
}

# 代码质量检查
lint() {
    info "运行代码质量检查..."
    
    # golangci-lint
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run --out-format colored-line-number > ${BUILD_DIR}/lint.txt 2>&1 || true
        success "代码检查完成，结果保存到: ${BUILD_DIR}/lint.txt"
    else
        warning "golangci-lint未安装，跳过代码检查"
    fi
    
    # go vet
    go vet ./... > ${BUILD_DIR}/vet.txt 2>&1 || true
    
    # go fmt检查
    unformatted=$(gofmt -l . | grep -v vendor || true)
    if [ -n "$unformatted" ]; then
        warning "以下文件需要格式化:"
        echo "$unformatted"
    fi
}

# 生成文档
docs() {
    info "生成项目文档..."
    
    # 生成Go文档
    if command -v godoc &> /dev/null; then
        info "生成Go文档..."
        mkdir -p ${BUILD_DIR}/docs
        # 这里可以添加文档生成逻辑
    fi
    
    # 生成API文档
    if [ -f "api/openapi.yaml" ]; then
        info "发现OpenAPI规范文件"
        # 这里可以添加API文档生成逻辑
    fi
}

# 安装依赖
install_deps() {
    info "安装构建依赖..."
    
    # 安装golangci-lint
    if ! command -v golangci-lint &> /dev/null; then
        info "安装golangci-lint..."
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
    fi
    
    success "依赖安装完成"
}

# 显示帮助信息
show_help() {
    cat << EOF
机器人路径编辑器构建脚本

用法: $0 [命令]

命令:
    build           构建当前平台的二进制文件
    build-all       构建所有平台的二进制文件
    test            运行测试
    benchmark       运行性能基准测试
    lint            运行代码质量检查
    docs            生成文档
    docker          构建Docker镜像
    release         创建发布包
    clean           清理构建文件
    install-deps    安装构建依赖
    help            显示此帮助信息

示例:
    $0 build-all    # 构建所有平台
    $0 test         # 运行测试
    $0 release      # 创建发布包

环境变量:
    VERSION         设置版本号 (默认: v1.0.0)
    BUILD_DIR       设置构建目录 (默认: build)
    DIST_DIR        设置发布目录 (默认: dist)

EOF
}

# 主函数
main() {
    case "${1:-build}" in
        "build")
            prepare
            build $(go env GOOS) $(go env GOARCH) "${PROJECT_NAME}"
            ;;
        "build-all")
            prepare
            build_all
            ;;
        "test")
            prepare
            test
            ;;
        "benchmark")
            prepare
            benchmark
            ;;
        "lint")
            prepare
            lint
            ;;
        "docs")
            prepare
            docs
            ;;
        "docker")
            prepare
            build_docker
            ;;
        "release")
            prepare
            test
            lint
            build_all
            create_release
            ;;
        "clean")
            cleanup
            ;;
        "install-deps")
            install_deps
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@"