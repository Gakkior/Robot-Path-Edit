#!/bin/bash

# 机器人路径编辑器部署脚本
# 支持多种部署方式: Docker, PM2, Systemd

set -e

# 配置
PROJECT_NAME="robot-path-editor"
VERSION="${VERSION:-latest}"
DEPLOY_MODE="${DEPLOY_MODE:-docker}"
SERVICE_PORT="${SERVICE_PORT:-8080}"
CONFIG_FILE="${CONFIG_FILE:-configs/config.yaml}"

# 路径配置
DEPLOY_DIR="/opt/${PROJECT_NAME}"
LOG_DIR="/var/log/${PROJECT_NAME}"
CONFIG_DIR="/etc/${PROJECT_NAME}"
DATA_DIR="/var/lib/${PROJECT_NAME}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

# 检查权限
check_permissions() {
    if [ "$EUID" -ne 0 ] && [ "$DEPLOY_MODE" != "docker" ]; then
        error "系统级部署需要root权限，请使用sudo运行此脚本"
        exit 1
    fi
}

# 创建目录结构
create_directories() {
    info "创建目录结构..."
    
    case "$DEPLOY_MODE" in
        "systemd"|"pm2"|"binary")
            mkdir -p "$DEPLOY_DIR"
            mkdir -p "$LOG_DIR"
            mkdir -p "$CONFIG_DIR"
            mkdir -p "$DATA_DIR"
            ;;
        "docker")
            mkdir -p "./data"
            mkdir -p "./logs"
            ;;
    esac
    
    success "目录创建完成"
}

# 下载/复制文件
deploy_files() {
    info "部署应用文件..."
    
    case "$DEPLOY_MODE" in
        "systemd"|"pm2"|"binary")
            # 复制二进制文件
            if [ -f "build/${PROJECT_NAME}" ]; then
                cp "build/${PROJECT_NAME}" "$DEPLOY_DIR/"
                chmod +x "$DEPLOY_DIR/${PROJECT_NAME}"
            else
                error "未找到二进制文件: build/${PROJECT_NAME}"
                exit 1
            fi
            
            # 复制web资源
            if [ -d "web" ]; then
                cp -r web "$DEPLOY_DIR/"
            fi
            
            # 复制配置文件
            if [ -f "$CONFIG_FILE" ]; then
                cp "$CONFIG_FILE" "$CONFIG_DIR/config.yaml"
            else
                # 创建默认配置
                create_default_config "$CONFIG_DIR/config.yaml"
            fi
            ;;
        "docker")
            # Docker部署通过docker-compose处理
            ;;
    esac
    
    success "文件部署完成"
}

# 创建默认配置
create_default_config() {
    local config_path=$1
    info "创建默认配置文件: $config_path"
    
    cat > "$config_path" << EOF
server:
  address: "0.0.0.0:${SERVICE_PORT}"
  mode: "release"

database:
  type: "sqlite"
  dsn: "${DATA_DIR}/robot_path_editor.db"

logger:
  level: "info"
  output: "${LOG_DIR}/app.log"
  format: "json"
EOF
    
    success "默认配置创建完成"
}

# Docker部署
deploy_docker() {
    info "使用Docker部署..."
    
    # 创建docker-compose.yml
    cat > docker-compose.yml << EOF
version: '3.8'

services:
  robot-path-editor:
    image: ${PROJECT_NAME}:${VERSION}
    container_name: ${PROJECT_NAME}
    restart: unless-stopped
    ports:
      - "${SERVICE_PORT}:8080"
    volumes:
      - ./data:/var/lib/${PROJECT_NAME}
      - ./logs:/var/log/${PROJECT_NAME}
      - ./configs:/etc/${PROJECT_NAME}
    environment:
      - GIN_MODE=release
      - CONFIG_PATH=/etc/${PROJECT_NAME}/config.yaml
    networks:
      - robot-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  # 可选: 添加数据库服务
  postgres:
    image: postgres:15-alpine
    container_name: ${PROJECT_NAME}-db
    restart: unless-stopped
    environment:
      POSTGRES_DB: robot_path_editor
      POSTGRES_USER: robot_user
      POSTGRES_PASSWORD: robot_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - robot-network
    ports:
      - "5432:5432"

  # 可选: 添加Redis缓存
  redis:
    image: redis:7-alpine
    container_name: ${PROJECT_NAME}-redis
    restart: unless-stopped
    volumes:
      - redis_data:/data
    networks:
      - robot-network
    ports:
      - "6379:6379"

  # 可选: 添加Prometheus监控
  prometheus:
    image: prom/prometheus:latest
    container_name: ${PROJECT_NAME}-prometheus
    restart: unless-stopped
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    networks:
      - robot-network

  # 可选: 添加Grafana可视化
  grafana:
    image: grafana/grafana:latest
    container_name: ${PROJECT_NAME}-grafana
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana
      - ./monitoring/grafana:/etc/grafana/provisioning
    networks:
      - robot-network

networks:
  robot-network:
    driver: bridge

volumes:
  postgres_data:
  redis_data:
  prometheus_data:
  grafana_data:
EOF
    
    # 创建监控配置
    mkdir -p monitoring
    cat > monitoring/prometheus.yml << EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: '${PROJECT_NAME}'
    static_configs:
      - targets: ['robot-path-editor:8080']
    metrics_path: '/metrics'
    scrape_interval: 5s
EOF
    
    # 启动服务
    docker-compose up -d
    
    success "Docker部署完成"
}

# Systemd部署
deploy_systemd() {
    info "使用Systemd部署..."
    
    # 创建systemd服务文件
    cat > "/etc/systemd/system/${PROJECT_NAME}.service" << EOF
[Unit]
Description=Robot Path Editor Service
After=network.target

[Service]
Type=simple
User=nobody
Group=nogroup
WorkingDirectory=${DEPLOY_DIR}
ExecStart=${DEPLOY_DIR}/${PROJECT_NAME}
ExecReload=/bin/kill -HUP \$MAINPID
KillMode=mixed
KillSignal=SIGTERM
RestartSec=5
Restart=always
StandardOutput=journal
StandardError=journal
SyslogIdentifier=${PROJECT_NAME}

# 环境变量
Environment=CONFIG_PATH=${CONFIG_DIR}/config.yaml
Environment=GIN_MODE=release

# 安全设置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=${DATA_DIR} ${LOG_DIR}

[Install]
WantedBy=multi-user.target
EOF
    
    # 重新加载systemd并启动服务
    systemctl daemon-reload
    systemctl enable "${PROJECT_NAME}.service"
    systemctl start "${PROJECT_NAME}.service"
    
    success "Systemd部署完成"
}

# PM2部署
deploy_pm2() {
    info "使用PM2部署..."
    
    # 检查PM2是否安装
    if ! command -v pm2 &> /dev/null; then
        info "安装PM2..."
        npm install -g pm2
    fi
    
    # 创建PM2配置文件
    cat > "${DEPLOY_DIR}/ecosystem.config.js" << EOF
module.exports = {
  apps: [{
    name: '${PROJECT_NAME}',
    script: '${DEPLOY_DIR}/${PROJECT_NAME}',
    cwd: '${DEPLOY_DIR}',
    instances: 1,
    autorestart: true,
    watch: false,
    max_memory_restart: '1G',
    env: {
      CONFIG_PATH: '${CONFIG_DIR}/config.yaml',
      GIN_MODE: 'release'
    },
    log_file: '${LOG_DIR}/combined.log',
    out_file: '${LOG_DIR}/out.log',
    error_file: '${LOG_DIR}/error.log',
    log_date_format: 'YYYY-MM-DD HH:mm Z'
  }]
};
EOF
    
    # 停止旧实例（如果存在）
    pm2 delete "${PROJECT_NAME}" 2>/dev/null || true
    
    # 启动新实例
    pm2 start "${DEPLOY_DIR}/ecosystem.config.js"
    pm2 save
    pm2 startup
    
    success "PM2部署完成"
}

# 二进制部署
deploy_binary() {
    info "二进制文件部署..."
    
    # 创建启动脚本
    cat > "${DEPLOY_DIR}/start.sh" << 'EOF'
#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_NAME="robot-path-editor"
PID_FILE="/var/run/${PROJECT_NAME}.pid"

start() {
    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo "服务已在运行"
        return
    fi
    
    echo "启动 $PROJECT_NAME..."
    nohup "$SCRIPT_DIR/$PROJECT_NAME" > /dev/null 2>&1 &
    echo $! > "$PID_FILE"
    echo "服务已启动，PID: $(cat $PID_FILE)"
}

stop() {
    if [ ! -f "$PID_FILE" ]; then
        echo "服务未运行"
        return
    fi
    
    PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        echo "停止 $PROJECT_NAME (PID: $PID)..."
        kill "$PID"
        sleep 2
        
        if kill -0 "$PID" 2>/dev/null; then
            echo "强制停止..."
            kill -9 "$PID"
        fi
        
        rm -f "$PID_FILE"
        echo "服务已停止"
    else
        echo "服务未运行"
        rm -f "$PID_FILE"
    fi
}

restart() {
    stop
    sleep 2
    start
}

status() {
    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo "服务正在运行，PID: $(cat $PID_FILE)"
    else
        echo "服务未运行"
    fi
}

case "$1" in
    start)   start ;;
    stop)    stop ;;
    restart) restart ;;
    status)  status ;;
    *)       echo "用法: $0 {start|stop|restart|status}" ;;
esac
EOF
    
    chmod +x "${DEPLOY_DIR}/start.sh"
    
    # 创建符号链接
    ln -sf "${DEPLOY_DIR}/start.sh" "/usr/local/bin/${PROJECT_NAME}"
    
    success "二进制部署完成，使用 '${PROJECT_NAME} start' 启动服务"
}

# 创建nginx配置
setup_nginx() {
    info "配置Nginx反向代理..."
    
    if ! command -v nginx &> /dev/null; then
        warning "Nginx未安装，跳过配置"
        return
    fi
    
    cat > "/etc/nginx/sites-available/${PROJECT_NAME}" << EOF
server {
    listen 80;
    server_name localhost;
    
    # 反向代理到应用
    location / {
        proxy_pass http://127.0.0.1:${SERVICE_PORT};
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        
        # WebSocket支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # 静态文件直接服务
    location /static/ {
        alias ${DEPLOY_DIR}/web/static/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # 健康检查
    location /health {
        access_log off;
        proxy_pass http://127.0.0.1:${SERVICE_PORT}/api/v1/health;
    }
    
    # 日志配置
    access_log ${LOG_DIR}/nginx-access.log;
    error_log ${LOG_DIR}/nginx-error.log;
}
EOF
    
    # 启用站点
    ln -sf "/etc/nginx/sites-available/${PROJECT_NAME}" "/etc/nginx/sites-enabled/"
    
    # 测试配置
    nginx -t && systemctl reload nginx
    
    success "Nginx配置完成"
}

# 健康检查
health_check() {
    info "执行健康检查..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f -s "http://localhost:${SERVICE_PORT}/api/v1/health" > /dev/null; then
            success "服务健康检查通过"
            return 0
        fi
        
        info "等待服务启动... (${attempt}/${max_attempts})"
        sleep 2
        ((attempt++))
    done
    
    error "服务健康检查失败"
    return 1
}

# 显示状态
show_status() {
    info "部署状态:"
    echo "  - 项目名称: $PROJECT_NAME"
    echo "  - 版本: $VERSION"
    echo "  - 部署模式: $DEPLOY_MODE"
    echo "  - 服务端口: $SERVICE_PORT"
    echo "  - 访问地址: http://localhost:${SERVICE_PORT}"
    
    case "$DEPLOY_MODE" in
        "systemd")
            echo "  - 服务状态: $(systemctl is-active ${PROJECT_NAME} 2>/dev/null || echo 'unknown')"
            echo "  - 管理命令: systemctl {start|stop|restart|status} ${PROJECT_NAME}"
            ;;
        "pm2")
            echo "  - PM2状态:"
            pm2 show "$PROJECT_NAME" 2>/dev/null || echo "    未运行"
            ;;
        "docker")
            echo "  - Docker状态:"
            docker-compose ps 2>/dev/null || echo "    未运行"
            ;;
        "binary")
            echo "  - 管理命令: ${PROJECT_NAME} {start|stop|restart|status}"
            ;;
    esac
}

# 卸载
uninstall() {
    info "卸载 $PROJECT_NAME..."
    
    case "$DEPLOY_MODE" in
        "systemd")
            systemctl stop "${PROJECT_NAME}" 2>/dev/null || true
            systemctl disable "${PROJECT_NAME}" 2>/dev/null || true
            rm -f "/etc/systemd/system/${PROJECT_NAME}.service"
            systemctl daemon-reload
            ;;
        "pm2")
            pm2 delete "${PROJECT_NAME}" 2>/dev/null || true
            ;;
        "docker")
            docker-compose down -v 2>/dev/null || true
            docker rmi "${PROJECT_NAME}:${VERSION}" 2>/dev/null || true
            ;;
        "binary")
            "${DEPLOY_DIR}/start.sh" stop 2>/dev/null || true
            rm -f "/usr/local/bin/${PROJECT_NAME}"
            ;;
    esac
    
    # 删除文件
    rm -rf "$DEPLOY_DIR"
    rm -rf "$CONFIG_DIR"
    rm -f "/etc/nginx/sites-enabled/${PROJECT_NAME}"
    rm -f "/etc/nginx/sites-available/${PROJECT_NAME}"
    
    success "卸载完成"
}

# 显示帮助
show_help() {
    cat << EOF
机器人路径编辑器部署脚本

用法: $0 [命令] [选项]

命令:
    deploy      部署应用
    status      显示部署状态
    health      执行健康检查
    nginx       配置Nginx反向代理
    uninstall   卸载应用
    help        显示帮助信息

环境变量:
    DEPLOY_MODE     部署模式 (docker|systemd|pm2|binary) 默认: docker
    VERSION         应用版本 默认: latest
    SERVICE_PORT    服务端口 默认: 8080
    CONFIG_FILE     配置文件路径 默认: configs/config.yaml

示例:
    # Docker部署
    $0 deploy

    # Systemd部署
    DEPLOY_MODE=systemd $0 deploy

    # PM2部署
    DEPLOY_MODE=pm2 $0 deploy

    # 检查状态
    $0 status

    # 配置Nginx
    $0 nginx

EOF
}

# 主函数
main() {
    case "${1:-deploy}" in
        "deploy")
            check_permissions
            create_directories
            deploy_files
            
            case "$DEPLOY_MODE" in
                "docker")    deploy_docker ;;
                "systemd")   deploy_systemd ;;
                "pm2")       deploy_pm2 ;;
                "binary")    deploy_binary ;;
                *)           error "不支持的部署模式: $DEPLOY_MODE"; exit 1 ;;
            esac
            
            health_check
            show_status
            ;;
        "status")
            show_status
            ;;
        "health")
            health_check
            ;;
        "nginx")
            setup_nginx
            ;;
        "uninstall")
            uninstall
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