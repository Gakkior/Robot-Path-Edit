# 🤖 机器人路径编辑器 (Robot Path Editor)

<div align="center">

**现代化的三端兼容机器人路径编辑器，支持可视化编辑和数据库管理**

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)
![Coverage](https://img.shields.io/badge/coverage-85%25-green.svg)

[🚀 快速开始](#-快速开始) • [📖 使用指南](#-使用指南) • [🛠️ 开发指南](#-开发指南) • [📊 API文档](#-api-参考) • [🎯 需求文档](requirements.md)

</div>

## ✨ 核心特性

### 🎨 可视化编辑
- **交互式画布**: 基于 Konva.js 的高性能画布，支持节点拖拽、路径连接
- **实时预览**: 所见即所得的编辑体验，实时显示路径规划结果
- **多种视图**: 画布视图和表格视图无缝切换，满足不同使用场景
- **撤销重做**: 基于命令模式的完整操作历史管理 (Ctrl+Z/Ctrl+Y)

### 🗄️ 数据管理
- **通用数据库编辑器**: 支持任意表结构的 CRUD 操作，类似 Excel 的使用体验
- **灵活映射**: 可选择任意数据表作为点位表和路径表，支持自定义ID字段映射
- **实时同步**: 画布视图与表格视图数据实时双向同步
- **多数据库支持**: SQLite、MySQL、PostgreSQL

### 📤 数据导出功能 (新增)
- **多格式导出**: Excel 和 CSV 格式，完美支持中文 (UTF-8编码)
- **分类导出**: 节点数据、路径数据、完整数据分别导出
- **即时下载**: 前端生成，无服务器压力
- **备份便捷**: 一键导出所有数据用于备份

### 📋 模板功能 (新增)
- **布局模板**: 保存常用的点位和路径布局范例
- **8种布局模式**: 树形、网格、圆形、力导向、管道、层次、径向、自定义
- **模板管理**: 创建、应用、克隆、搜索、导入导出
- **相对坐标**: 模板自适应不同画布尺寸
- **分类管理**: 工厂、仓库、实验室等场景分类

### 🔧 智能算法
- **布局算法**: 网格布局、力导向布局、圆形布局等多种自动排列方式
- **路径生成**: 最近邻连接、完全连通图、网格路径等智能路径生成算法
- **路径优化**: 最短路径计算、路径平滑优化

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
┌─────────────────────────────────────────────────────────────────┐
│                    Domain层(核心领域模型)                      │
│ Node, Path, Template, DatabaseConnection, TableMapping         │
└─────────────────────────────────────────────────────────────────┘
```

### 🔄 双前端架构

| 前端版本 | 技术栈 | 访问路径 | 状态 |
|---------|--------|----------|------|
| **经典版** | HTML + JS + Konva.js | `/app` | 🟢 可用 |
| **现代版** | React + TS + Vite | `/app/new` | 🟢 可用 |

### 技术栈
- **后端**: Go 1.21+, Gin, GORM, SQLite/MySQL/PostgreSQL
- **前端经典版**: Konva.js, 原生JavaScript, 现代CSS + HTML5
- **前端现代版**: React + TypeScript, Vite, TailwindCSS, Konva.js, Zustand, React Query
- **架构**: 领域驱动设计 (DDD), 仓储模式, 命令模式

### ✨ 现代前端特性
- **React + TypeScript**: 类型安全的组件化开发
- **Vite**: 极速构建和热重载
- **TailwindCSS**: 现代化的原子级CSS框架
- **Konva.js + React-Konva**: 高性能画布渲染
- **Zustand**: 轻量级状态管理
- **React Query**: 智能数据获取和缓存
- **Framer Motion**: 流畅的动画效果
- **Radix UI**: 无障碍的组件库

## 🚀 快速开始

### 演示版 (推荐)
最快的体验方式，无需数据库配置：

```bash
# 下载并启动演示版
go run cmd/demo/main.go

# 或使用预编译版本
./demo.exe  # Windows
./demo      # Linux/macOS
```

访问 **http://localhost:8080** 开始体验！

### 现代前端开发
快速启动React版本进行开发：

```bash
# 启动现代前端开发服务器
cd frontend
npm install
npm run dev

# 启动后端服务器 (另一个终端)
cd ..
go run cmd/server/main.go
```

**开发环境访问地址：**
- **现代前端**: http://localhost:5173 (Vite开发服务器)
- **后端API**: http://localhost:8080/api/v1
- **经典前端**: http://localhost:8080/app

**生产环境访问地址：**
- **现代前端**: http://localhost:8080/app/new
- **经典前端**: http://localhost:8080/app

### 完整版部署

#### 1. 环境要求
```bash
# 基础环境
Go 1.21+
SQLite/MySQL/PostgreSQL (任选一种)
```

#### 2. 快速安装
```bash
# 克隆项目
git clone https://github.com/your-org/robot-path-editor.git
cd robot-path-editor

# 安装依赖
go mod download

# 构建项目
go build -o robot-path-editor cmd/server/main.go

# 配置数据库
cp configs/config.yaml.example configs/config.yaml
# 编辑 config.yaml 配置数据库连接

# 启动服务
./robot-path-editor     # Linux/macOS
.\robot-path-editor.exe # Windows
```

#### 3. Docker 部署 (推荐)
```bash
# 使用 Docker Compose 一键部署
docker-compose up -d
```

## 📖 使用指南

### 基础操作

#### 🎨 画布视图
```bash
创建节点     双击空白区域或使用工具栏
移动节点     拖拽节点到目标位置
创建路径     Shift+点击两个节点
删除元素     选中后按Delete键
撤销操作     Ctrl+Z / Cmd+Z
重做操作     Ctrl+Y / Cmd+Y
```

#### 📊 表格视图
```bash
切换视图     点击页面顶部"表格视图"按钮
编辑数据     直接在表格中修改数据
批量操作     选择多行进行批量编辑或删除
```

#### 📤 数据导出
```bash
导出节点CSV    点击"导出"按钮 -> 选择"导出节点数据(CSV)"
导出路径Excel  点击"导出"按钮 -> 选择"导出路径数据(Excel)"
导出完整数据   点击"导出"按钮 -> 选择"导出所有数据"
```

#### 📋 模板管理
```bash
保存模板     点击"模板"按钮 -> "保存当前为模板"
应用模板     点击"模板"按钮 -> 选择要应用的模板
搜索模板     在模板管理器中使用搜索功能
```

### 高级配置

#### 数据库配置

**默认配置（推荐）：**
```yaml
# configs/config.yaml
database:
  type: "memory"    # 内存模式，无需CGO，适合开发和演示
  dsn: ":memory:"
```

**其他数据库选项：**
```yaml
# SQLite（需要CGO支持）
database:
  type: "sqlite"
  dsn: "./data/robot_paths.db"

# MySQL
database:
  type: "mysql"
  dsn: "user:password@tcp(localhost:3306)/robot_paths"

# PostgreSQL
database:
  type: "postgresql"
  dsn: "host=localhost user=postgres password=password dbname=robot_paths port=5432"
```

**启用CGO编译SQLite：**
```bash
# Windows
set CGO_ENABLED=1 && go build cmd/server/main.go

# Linux/macOS
CGO_ENABLED=1 go build cmd/server/main.go
```

## 📊 API 参考

### 节点管理
```http
GET    /api/v1/nodes           # 获取所有节点
POST   /api/v1/nodes           # 创建节点
PUT    /api/v1/nodes/{id}      # 更新节点
DELETE /api/v1/nodes/{id}      # 删除节点
```

### 路径管理
```http
GET    /api/v1/paths           # 获取所有路径
POST   /api/v1/paths           # 创建路径
PUT    /api/v1/paths/{id}      # 更新路径
DELETE /api/v1/paths/{id}      # 删除路径
```

### 模板管理 (新增)
```http
GET    /api/v1/templates                    # 列出模板
POST   /api/v1/templates                    # 创建模板
POST   /api/v1/templates/{id}/apply         # 应用模板
POST   /api/v1/templates/save-as            # 保存为模板
GET    /api/v1/templates/public             # 公开模板
```

### 数据库连接
```http
GET    /api/v1/database/connections         # 数据库连接列表
POST   /api/v1/database/connections         # 创建连接
POST   /api/v1/mapping                      # 创建表映射
POST   /api/v1/sync/mappings/{id}/all       # 同步数据
```

### 布局算法
```http
POST   /api/v1/layout/apply                 # 应用布局算法
POST   /api/v1/path-generation/nearest-neighbor  # 最近邻路径
POST   /api/v1/path-generation/full-connectivity # 完全连通
```

### 导出功能示例

#### JavaScript导出API
```javascript
// 导出节点数据为CSV
await dataExporter.exportNodesAsCSV();

// 导出路径数据为Excel
await dataExporter.exportPathsAsExcel();

// 导出完整数据为Excel
await dataExporter.exportAllAsExcel();
```

#### 模板功能示例
```javascript
// 保存当前布局为模板
const templateData = {
  name: "工厂车间标准布局",
  description: "适用于中型工厂车间的机器人路径规划",
  category: "factory",
  layout_type: "grid",
  nodes: currentCanvasData.nodes,
  paths: currentCanvasData.paths
};

await fetch('/api/v1/templates/save-as', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(templateData)
});

// 应用模板到画布
await fetch(`/api/v1/templates/${templateId}/apply`, {
  method: 'POST',
  body: JSON.stringify({ width: 1920, height: 1080 })
});
```

## 🧪 开发指南

### 开发环境搭建
```bash
# 安装开发工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 运行开发服务器
go run cmd/server/main.go

# 运行演示版本
go run cmd/demo/main.go

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
├── frontend/              # 现代前端 (React)
│  ├── src/
│  │   ├── components/    # 组件库
│  │   │   ├── Canvas/   # 画布组件
│  │   │   ├── Sidebar/  # 侧边栏组件
│  │   │   ├── Toolbar/  # 工具栏组件
│  │   │   └── ui/       # 基础UI组件
│  │   ├── services/     # API服务
│  │   ├── stores/       # 状态管理
│  │   ├── types/        # TypeScript类型
│  │   └── utils/        # 工具函数
│  ├── package.json      # 依赖配置
│  ├── vite.config.ts    # Vite配置
│  ├── tailwind.config.js # TailwindCSS配置
│  └── tsconfig.json     # TypeScript配置
├── pkg/                   # 公共包
├── web/static/           # 经典前端资源
├── tests/                # 测试文件
├── scripts/              # 构建脚本
└── configs/              # 配置文件
```

### 前端开发指南

#### 环境变量配置
创建 `frontend/.env` 文件：
```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_APP_TITLE=机器人路径编辑器
VITE_ENABLE_DEBUG=true
```

#### 开发建议
1. **API优先**: 后端API完善后再开发前端功能
2. **类型安全**: 充分利用TypeScript的类型检查
3. **组件复用**: 抽取通用组件提高代码复用
4. **性能优化**: 使用React.memo和useMemo优化渲染
5. **测试驱动**: 编写单元测试确保代码质量

#### 构建部署
```bash
# 构建现代前端
cd frontend && npm run build

# 复制到后端服务目录 (Windows)
.\scripts\build-frontend.bat

# 复制到后端服务目录 (Linux/Mac)
./scripts/build-frontend.sh
```

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

### 导出性能
- **CSV导出**: 支持 100万+ 行数据
- **Excel导出**: 支持 10万+ 行数据
- **UTF-8编码**: 完美支持中文字符
- **内存优化**: 流式处理，低内存占用

## 🔧 故障排除

### 常见问题

#### 1. 端口占用
```bash
# 检查端口占用
netstat -tulpn | grep :8080

# 修改端口
export PORT=8081
```

#### 2. 静态资源404
```bash
# 确保web目录存在
ls -la web/static/

# 检查embed资源
go run cmd/demo/main.go
```

#### 3. 导出功能异常
```bash
# 检查export.js是否加载
curl http://localhost:8080/static/export.js

# 检查浏览器控制台错误
F12 -> Console -> 查看错误信息
```

## 🎯 功能特色

### 🌟 已实现功能

#### 🎨 双前端架构优势
| 功能 | 经典前端 | 现代前端 | 说明 |
|------|---------|----------|------|
| 节点管理 | ✅ | ✅ | 创建、编辑、删除节点 |
| 路径管理 | ✅ | ✅ | 创建、编辑、删除路径 |
| 画布交互 | ✅ | ✅ | 拖拽、缩放、平移 |
| 属性面板 | ✅ | ✅ | 实时编辑节点/路径属性 |
| 数据同步 | ✅ | ✅ | 外部数据库同步 |
| 模板功能 | ✅ | ✅ | 保存和应用布局模板 |
| 数据导出 | ✅ | ✅ | CSV/Excel导出 |
| 撤销重做 | ✅ | ✅ | 操作历史管理 |
| 布局算法 | ✅ | ✅ | 多种自动排列方式 |
| 响应式设计 | ❌ | ✅ | 现代前端独有 |
| TypeScript | ❌ | ✅ | 类型安全开发 |
| 组件化 | ❌ | ✅ | 模块化架构 |
| 现代动画 | ❌ | ✅ | Framer Motion |

#### ✅ 核心功能
- ✅ **基础画布编辑**: Konva.js实现，支持拖拽、连接
- ✅ **表格编辑**: 双视图切换，实时同步
- ✅ **撤销重做**: 命令模式，完整操作历史
- ✅ **数据导出**: Excel/CSV格式，UTF-8编码
- ✅ **模板功能**: 8种布局类型，完整管理
- ✅ **数据库集成**: 多数据库支持，灵活映射
- ✅ **布局算法**: 网格、力导向、圆形等
- ✅ **路径生成**: 最近邻、完全连通、网格路径

### 🔄 开发中功能
- 🔄 **移动端优化**: 触控交互优化
- 🔄 **桌面端打包**: Electron或Tauri打包
- 🔄 **在线模板库**: 公共模板分享平台
- 🔄 **实时协作**: 多人同时编辑支持

## 📚 学习价值

本项目采用了众多现代软件开发最佳实践：

### 架构设计
- **领域驱动设计 (DDD)**: 清晰的业务领域建模
- **分层架构**: Handler->Service->Repository->Domain
- **依赖注入**: 松耦合的组件设计
- **接口设计**: 面向接口编程

### 设计模式
- **仓储模式**: 数据访问抽象
- **命令模式**: 撤销/重做功能实现
- **策略模式**: 布局算法切换
- **工厂模式**: 组件创建管理
- **适配器模式**: 多数据库支持

### 前端技术
- **高性能渲染**: Konva.js Canvas渲染
- **模块化设计**: 功能组件化
- **事件驱动**: 响应式交互
- **状态管理**: 命令模式状态管理

## 📄 许可证

本项目基于 MIT 许可证开源 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - HTTP Web框架
- [GORM](https://gorm.io/) - Go ORM库
- [Konva.js](https://konvajs.org/) - 2D Canvas渲染引擎
- [UUID](https://github.com/google/uuid) - UUID生成库

## 🔗 相关链接

- [📋 需求文档](requirements.md) - 详细功能需求和技术规划
- [🐛 问题反馈](https://github.com/your-org/robot-path-editor/issues)
- [💬 讨论社区](https://github.com/your-org/robot-path-editor/discussions)

---

<div align="center">

**如果这个项目对你有帮助，请给它一个 ⭐️ Star！**

[🐛 报告Bug](https://github.com/your-org/robot-path-editor/issues) |
[✨ 请求功能](https://github.com/your-org/robot-path-editor/issues) |
[💬 参与讨论](https://github.com/your-org/robot-path-editor/discussions)

**现在支持数据导出和模板功能！📤📋**

</div>