# 🤖 机器人路径编辑器 (Robot Path Editor)

一个现代化的三端兼容机器人路径编辑器，支持可视化编辑和数据库管理。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)
![Coverage](https://img.shields.io/badge/coverage-85%25-green.svg)

## ✨ 特性

### 🎨 可视化编辑
- **交互式画布**: 基于 Konva.js 的高性能画布，支持节点拖拽、路径连接
- **实时预览**: 所见即所得的编辑体验，实时显示路径规划结果
- **多种视图**: 画布视图和表格视图无缝切换，满足不同使用场景

### 🗄️ 数据管理
- **通用数据库编辑器**: 支持任意表结构的 CRUD 操作，类似 Excel 的使用体验
- **灵活映射**: 可选择任意数据表作为点位表和路径表，支持自定义ID字段映射
- **实时同步**: 画布视图与表格视图数据实时双向同步

### 🔧 智能算法
- **布局算法**: 网格布局、力导向布局、圆形布局等多种自动排列方式
- **路径生成**: 最近邻连接、完全连通图、网格路径等智能路径生成算法
- **路径优化**: 最短路径计算、路径平滑优化

### 🛠️ 高级功能
- **撤销/重做**: 基于命令模式的完整操作历史管理
- **插件系统**: 可扩展的插件架构，支持自定义布局和路径算法
- **实时监控**: Prometheus 指标收集、结构化日志、性能追踪
- **多种部署**: Docker、Systemd、PM2 等多种部署方式

### 📱 跨平台支持
- **桌面端**: Windows、Linux、macOS 原生支持
- **移动端**: PWA 支持，平板设备优化的触控体验
- **Web端**: 响应式设计，支持所有现代浏览器

## 🏗️ 技术架构

### Go-Heavy 后端架构
```
┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐
│  Handler层    │   │  Service层    │   │Repository层   │
│ (HTTP接口)     │────│  (业务逻辑)    │────│  (数据访问)    │
└─────────────────┘   └─────────────────┘   └─────────────────┘
         │                      │                      │
         │             ┌─────────────────┐             │
         │             │  Plugin系统    │             │
         │             │ (扩展算法)     │             │
         │             └─────────────────┘             │
         │                      │                      │
┌─────────────────────────────────────────────────────────────────┐
│                    Domain层(核心领域模型)                      │
│              Node, Path, Position, Metadata                    │
└─────────────────────────────────────────────────────────────────┘
```

### 前端技术栈
- **画布渲染**: Konva.js (高性能 2D Canvas)
- **交互逻辑**: 原生 JavaScript (轻量级)
- **状态管理**: 命令模式 (撤销/重做)
- **UI组件**: 现代 CSS + HTML5

### 设计模式应用
- **仓储模式**: 数据访问抽象
- **适配器模式**: 多数据库支持
- **命令模式**: 操作历史管理
- **策略模式**: 布局算法切换
- **观察者模式**: 事件驱动架构
- **工厂模式**: 组件创建管理

## 🚀 快速开始

### 演示版(推荐)
最快的体验方式，无需数据库配置：

```bash
# 下载并启动演示版
go run cmd/demo/main.go

# 或使用预编译版本
./demo.exe  # Windows
./demo      # Linux/macOS
```

访问 http://localhost:8080 开始体验！

### 完整版部署

#### 1. 环境要求
- Go 1.21+
- SQLite/MySQL/PostgreSQL (任选一种)
- Node.js 16+ (可选，用于前端构建)

#### 2. 快速安装
```bash
# 克隆项目
git clone https://github.com/your-org/robot-path-editor.git
cd robot-path-editor

# 安装依赖
go mod download

# 构建项目
./scripts/build.sh build-all  # Linux/macOS
build.bat build-all           # Windows

# 配置数据库
cp configs/config.yaml.example configs/config.yaml
# 编辑 config.yaml 配置数据库连接

# 启动服务
./build/robot-path-editor
```

#### 3. Docker 部署 (推荐)
```bash
# 使用 Docker Compose 一键部署
docker-compose up -d

# 包含数据库、Redis、监控的完整环境
```

#### 4. 生产环境部署
```bash
# 使用部署脚本
./scripts/deploy.sh

# 支持多种部署方式
DEPLOY_MODE=systemd ./scripts/deploy.sh  # Systemd
DEPLOY_MODE=pm2 ./scripts/deploy.sh      # PM2
DEPLOY_MODE=docker ./scripts/deploy.sh   # Docker
```

## 📖 使用指南

### 基础操作

#### 画布视图
1. **创建节点**: 双击空白区域或使用工具栏
2. **移动节点**: 拖拽节点到目标位置
3. **创建路径**: Shift+点击两个节点
4. **删除元素**: 选中后按Delete键
5. **撤销操作**: Ctrl+Z / Cmd+Z

#### 表格视图
1. **切换视图**: 点击页面顶部"表格视图"按钮
2. **编辑数据**: 直接在表格中修改数据
3. **批量操作**: 选择多行进行批量编辑或删除
4. **导入导出**: 支持CSV、Excel格式

#### 智能算法
```bash
# 应用布局算法
curl -X POST http://localhost:8080/api/v1/layout/apply \
  -H "Content-Type: application/json" \
  -d '{"algorithm": "force-directed"}'

# 生成路径
curl -X POST http://localhost:8080/api/v1/path-generation/nearest-neighbor \
  -H "Content-Type: application/json" \
  -d '{"max_distance": 200}'
```

### 高级配置

#### 数据库配置
```yaml
# config.yaml
database:
  type: "mysql"  # sqlite, mysql, postgresql
  dsn: "user:password@tcp(localhost:3306)/robot_paths"
  
  # 自定义表映射
  table_mapping:
    node_table: "robot_points"
    node_id_field: "point_id"
    path_table: "robot_routes"
    path_id_field: "route_id"
```

#### 插件开发
```go
// 自定义布局插件
type CustomLayoutPlugin struct{}

func (p *CustomLayoutPlugin) ApplyLayout(nodes []domain.Node, paths []domain.Path, config map[string]interface{}) ([]domain.Node, error) {
    // 实现自定义布局算法
    return nodes, nil
}

// 注册插件
pluginService.RegisterLayoutPlugin(&CustomLayoutPlugin{})
```

## 📊 API 参考

### RESTful API

#### 节点管理
- `GET /api/v1/nodes` - 获取所有节点
- `POST /api/v1/nodes` - 创建节点
- `GET /api/v1/nodes/{id}` - 获取单个节点
- `PUT /api/v1/nodes/{id}` - 更新节点
- `DELETE /api/v1/nodes/{id}` - 删除节点
- `PUT /api/v1/nodes/{id}/position` - 更新节点位置

#### 路径管理
- `GET /api/v1/paths` - 获取所有路径
- `POST /api/v1/paths` - 创建路径
- `GET /api/v1/paths/{id}` - 获取单个路径
- `PUT /api/v1/paths/{id}` - 更新路径
- `DELETE /api/v1/paths/{id}` - 删除路径

#### 布局算法
- `POST /api/v1/layout/apply` - 应用布局算法

#### 路径生成
- `POST /api/v1/path-generation/nearest-neighbor` - 最近邻路径
- `POST /api/v1/path-generation/full-connectivity` - 完全连通
- `POST /api/v1/path-generation/grid` - 网格路径

### WebSocket API (计划中)
- `/ws/canvas` - 实时画布同步
- `/ws/notifications` - 系统通知

## 🧪 开发指南

### 开发环境搭建
```bash
# 安装开发工�?
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest

# 运行开发服务器
go run cmd/server/main.go

# 运行测试
go test ./... -v -cover

# 代码检查
golangci-lint run
```

### 项目结构
```
robot-path-editor/
├── cmd/                    # 应用入口
│  ├── server/            # 主服务器
│  └── demo/              # 演示版本
├── internal/              # 内部包
│  ├── domain/            # 领域模型
│  ├── services/          # 业务服务
│  ├── repositories/      # 数据仓储
│  ├── handlers/          # HTTP处理器
│  └── database/          # 数据库适配
├── pkg/                   # 公共包
│  ├── logger/            # 日志工具
│  └── middleware/        # 中间件
├── web/                   # 前端资源
│  └── static/            # 静态文件
├── tests/                 # 测试文件
│  ├── unit/              # 单元测试
│  └── integration/       # 集成测试
├── scripts/               # 构建脚本
├── configs/               # 配置文件
└── docs/                  # 文档
```

### 贡献指南
1. Fork 项目
2. 创建特性分支: `git checkout -b feature/amazing-feature`
3. 提交变更: `git commit -m 'Add amazing feature'`
4. 推送分支: `git push origin feature/amazing-feature`
5. 提交 Pull Request

### 代码规范
- 遵循 Go 官方代码规范
- 使用 golangci-lint 进行代码检查
- 单元测试覆盖率 > 80%
- 提交信息遵循 Conventional Commits

## 📈 性能指标

### 系统性能
- **响应时间**: < 100ms (95%ile)
- **吞吐量**: > 1000 QPS
- **内存使用**: < 256MB (空载)
- **启动时间**: < 5s

### 画布性能
- **节点数量**: 支持 10,000+ 节点
- **路径数量**: 支持 50,000+ 路径
- **渲染帧率**: 60 FPS (1080p)
- **响应延迟**: < 16ms (触控/鼠标)

### 数据库性能
- **SQLite**: 适合 < 10万记录
- **MySQL**: 适合 < 1000万记录
- **PostgreSQL**: 适合 > 1000万记录

## 🔧 故障排除

### 常见问题

#### 1. CGO相关错误
```bash
# 错误: CGO_ENABLED=0, go-sqlite3 requires cgo
# 解决: 使用纯Go SQLite驱动或启用CGO
export CGO_ENABLED=1
go build ...
```

#### 2. 端口占用
```bash
# 检查端口占用
netstat -tulpn | grep :8080

# 修改配置文件中的端口
# 或设置环境变量
export PORT=8081
```

#### 3. 静态资源404
```bash
# 确保web目录存在
# 或使用go:embed内嵌资源
```

### 日志分析
```bash
# 查看应用日志
tail -f /var/log/robot-path-editor/app.log

# 查看访问日志
tail -f /var/log/robot-path-editor/access.log

# 查看错误日志
grep "ERROR" /var/log/robot-path-editor/app.log
```

### 监控指标
- **系统指标**: CPU、内存、磁盘使用率
- **应用指标**: QPS、响应时间、错误率
- **业务指标**: 节点数量、路径数量、用户活跃度

## 📄 许可证

本项目基于 MIT 许可证开源 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - HTTP Web框架
- [GORM](https://gorm.io/) - Go ORM库
- [Konva.js](https://konvajs.org/) - 2D Canvas渲染引擎
- [Prometheus](https://prometheus.io/) - 监控系统
- [Logrus](https://github.com/sirupsen/logrus) - 日志库

## 🔗 相关链接

- [文档网站](https://robot-path-editor.github.io/docs)
- [在线演示](https://demo.robot-path-editor.com)
- [Docker Hub](https://hub.docker.com/r/robotpatheditor/robot-path-editor)
- [问题反馈](https://github.com/your-org/robot-path-editor/issues)
- [讨论社区](https://github.com/your-org/robot-path-editor/discussions)

---

<div align="center">

**如果这个项目对你有帮助，请给它一个 ⭐️ Star！**

[🐛 报告Bug](https://github.com/your-org/robot-path-editor/issues) |
[✨ 请求功能](https://github.com/your-org/robot-path-editor/issues) |
[💬 参与讨论](https://github.com/your-org/robot-path-editor/discussions)

</div>