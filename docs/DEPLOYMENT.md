# 🚀 部署指南

## 快速开始

### 方式一：演示版 (推荐)
最快的体验方式，内置SQLite数据库：

```bash
# 下载并运行
go run cmd/demo/main.go

# 或使用预编译版本
./build/robot-path-editor-demo.exe  # Windows
./build/robot-path-editor-demo      # Linux/macOS
```

访问 **http://localhost:8080** 开始使用！

### 方式二：Docker 部署 (推荐生产环境)

```bash
# 1. 克隆项目
git clone <your-repo-url>
cd robot-path-editor

# 2. 启动应用 (SQLite版本)
docker-compose up -d app

# 3. 启动应用 (MySQL版本)
docker-compose --profile mysql up -d
```

访问 **http://localhost:8080**

### 方式三：手动编译部署

```bash
# 1. 环境要求
Go 1.21+
Node.js 18+ (开发时需要)

# 2. 编译后端
./scripts/build.sh        # Linux/macOS
./scripts/build.bat       # Windows

# 3. 编译前端 (可选，已内置)
cd frontend
npm install
npm run build

# 4. 启动服务
./build/robot-path-editor
```

## 开发环境

### 启动开发服务器

```bash
# 方式一：自动启动脚本
./scripts/dev-start.sh        # Linux/macOS
./scripts/dev-start.bat       # Windows

# 方式二：手动启动
# 终端1: 启动后端
go run cmd/server/main.go

# 终端2: 启动前端 (如需开发前端)
cd frontend
npm install
npm run dev
```

**访问地址：**
- 现代前端 (开发): http://localhost:5173
- 经典前端: http://localhost:8080/app
- 现代前端 (生产): http://localhost:8080/app/new
- API文档: http://localhost:8080/api/v1

## 配置说明

### 数据库配置

创建 `configs/config.yaml` 文件：

```yaml
# SQLite (推荐)
database:
  type: "sqlite"
  dsn: "./data/robot-path-editor.db"

# MySQL
database:
  type: "mysql"
  dsn: "user:password@tcp(localhost:3306)/robot_paths"

# 服务器配置
server:
  host: "0.0.0.0"
  port: 8080

# 日志配置
logger:
  level: "info"
```

### 环境变量

```bash
# 数据库设置
export DATABASE_TYPE=sqlite
export DATABASE_DSN=./data/robot-path-editor.db

# 服务器设置
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080

# 日志级别
export LOG_LEVEL=info
```

## 故障排除

### 常见问题

**1. 端口被占用**
```bash
# 检查端口占用
netstat -tlnp | grep :8080

# 修改端口
export SERVER_PORT=8081
```

**2. 数据库连接失败**
```bash
# 检查SQLite数据库文件权限
ls -la data/

# 检查MySQL连接
mysql -h localhost -u root -p
```

**3. 前端资源404**
```bash
# 确保web目录存在
ls -la web/static/

# 重新构建前端
cd frontend && npm run build
./scripts/build-frontend.sh
```

**4. Go编译失败**
```bash
# 检查Go版本
go version

# 清理模块缓存
go clean -modcache
go mod download
```

## 性能优化

### 生产环境建议

1. **使用MySQL/PostgreSQL替代SQLite**
2. **启用Nginx反向代理**
3. **配置HTTPS证书**
4. **设置日志轮转**
5. **监控资源使用**

### Nginx配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 更新部署

### Docker环境更新

```bash
# 1. 停止服务
docker-compose down

# 2. 更新代码
git pull

# 3. 重新构建和启动
docker-compose up -d --build
```

### 手动环境更新

```bash
# 1. 备份数据
cp -r data data.backup

# 2. 更新代码
git pull

# 3. 重新编译
./scripts/build.sh

# 4. 重启服务
./build/robot-path-editor
```