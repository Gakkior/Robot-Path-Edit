# 📁 项目结构

## 总体架构

```
robot-path-editor/
├── cmd/                    # 应用入口
│   ├── server/            # 主服务器
│   └── demo/              # 演示版本
├── internal/              # 内部包 (Go后端)
│   ├── app/              # 应用程序组装器
│   ├── config/           # 配置管理
│   ├── database/         # 数据库适配器
│   ├── domain/           # 领域模型
│   ├── handlers/         # HTTP处理器
│   ├── repositories/     # 数据仓储层
│   └── services/         # 业务服务层
├── frontend/              # 现代前端 (React)
│   ├── src/
│   │   ├── components/   # React组件
│   │   ├── services/     # API服务
│   │   ├── stores/       # 状态管理
│   │   ├── types/        # TypeScript类型
│   │   └── utils/        # 工具函数
│   ├── package.json      # 前端依赖
│   └── vite.config.ts    # Vite配置
├── pkg/                   # 公共包
│   ├── logger/           # 日志组件
│   └── middleware/       # 中间件
├── web/                   # 经典前端资源
│   └── static/           # 静态文件
├── scripts/               # 构建和部署脚本
│   ├── build.sh          # Linux/macOS构建脚本
│   ├── build.bat         # Windows构建脚本
│   ├── dev-start.sh      # Linux开发启动脚本
│   ├── dev-start.bat     # Windows开发启动脚本
│   ├── build-frontend.sh # 前端构建脚本
│   └── init.sql          # 数据库初始化脚本
├── docs/                  # 项目文档
│   ├── DEPLOYMENT.md     # 部署指南
│   ├── API.md            # API文档
│   └── PROJECT_STRUCTURE.md # 项目结构说明
├── configs/               # 配置文件
│   ├── config.yaml       # 主配置文件
│   └── config.yaml.example # 配置示例
├── data/                  # 数据存储目录
├── tests/                 # 测试文件
│   ├── unit/             # 单元测试
│   └── integration/      # 集成测试
├── build/                 # 构建输出目录
├── docker-compose.yml     # Docker Compose配置
├── Dockerfile            # Docker镜像构建文件
├── go.mod               # Go模块定义
├── go.sum               # Go依赖锁定
├── README.md            # 项目主文档
└── CHANGELOG.md         # 变更日志
```

## 核心目录说明

### `/cmd` - 应用入口
- `server/main.go`: 完整版服务器，支持数据库配置
- `demo/main.go`: 演示版本，内置SQLite，无需配置

### `/internal` - 后端核心代码
采用领域驱动设计 (DDD) 架构：

```
Handler层 → Service层 → Repository层 → Domain层
    ↓          ↓            ↓           ↓
HTTP接口   业务逻辑      数据访问      核心模型
```

- **Domain**: 核心业务模型 (Node, Path, Template等)
- **Services**: 业务逻辑层，处理复杂业务规则
- **Repositories**: 数据访问层，抽象不同数据源
- **Handlers**: HTTP处理器，处理Web请求

### `/frontend` - 现代前端
React + TypeScript + Vite 技术栈：

- **components/**: 可复用的React组件
  - `Canvas/`: 画布相关组件
  - `Sidebar/`: 侧边栏组件
  - `Toolbar/`: 工具栏组件
  - `ui/`: 基础UI组件 (基于Radix UI)
- **services/**: API服务和数据获取
- **stores/**: Zustand状态管理
- **types/**: TypeScript类型定义
- **utils/**: 工具函数

### `/web` - 经典前端
原生JavaScript + HTML + CSS：

- `static/`: 静态资源文件
- `app.js`: 主应用逻辑
- `canvas.js`: 画布操作
- `table.js`: 表格功能

### `/scripts` - 构建脚本
- `build.sh/build.bat`: 应用程序构建
- `dev-start.sh/dev-start.bat`: 开发环境启动
- `build-frontend.sh`: 前端资源构建

### `/docs` - 项目文档
- `DEPLOYMENT.md`: 详细部署指南
- `API.md`: API接口文档
- `PROJECT_STRUCTURE.md`: 本文档

## 数据流架构

### 请求处理流程
```
HTTP请求 → Gin路由 → Handler → Service → Repository → Database
                                    ↓
HTTP响应 ← JSON序列化 ← Domain模型 ← 数据查询结果
```

### 前端架构
```
用户交互 → React组件 → Zustand状态 → API服务 → 后端接口
                         ↓
画布更新 ← 状态变更 ← React Query缓存 ← API响应
```

## 文件命名规范

### Go代码
- 文件名：`snake_case.go`
- 包名：`lowercase`
- 类型名：`PascalCase`
- 函数名：`camelCase` (公开) / `camelCase` (私有)

### 前端代码
- 组件文件：`PascalCase.tsx`
- 工具文件：`camelCase.ts`
- 样式文件：`kebab-case.css`

### 配置文件
- 配置：`kebab-case.yaml`
- 脚本：`kebab-case.sh/.bat`
- 文档：`UPPERCASE.md`

## 端口分配

| 服务 | 端口 | 说明 |
|------|------|------|
| 后端API | 8080 | 主服务器端口 |
| 前端开发 | 5173 | Vite开发服务器 |
| MySQL | 3306 | 数据库端口 |
| Redis | 6379 | 缓存端口 (可选) |

## 环境区分

### 开发环境
- 前端：Vite开发服务器 (5173)
- 后端：`go run` 直接运行
- 数据库：本地SQLite

### 生产环境
- 前端：静态文件嵌入后端
- 后端：编译后的二进制文件
- 数据库：MySQL/PostgreSQL

### Docker环境
- 应用：容器化部署
- 数据库：Docker Compose管理
- 网络：内部容器网络