# 🤖 机器人路径编辑器

**现代化的可视化机器人路径编辑器，支持画布编辑、数据库管理和模板功能**

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)

[🚀 快速开始](#-快速开始) • [📖 功能特性](#-功能特性) • [🛠️ 部署指南](docs/DEPLOYMENT.md) • [📊 API文档](docs/API.md)

## ✨ 功能特性

### 🎨 可视化编辑
- **交互式画布**: 基于 Konva.js 的高性能画布，支持拖拽、缩放、连接
- **双前端架构**: 经典版本 + React现代版本，满足不同需求
- **实时预览**: 所见即所得的编辑体验
- **撤销重做**: 完整的操作历史管理 (Ctrl+Z/Ctrl+Y)

### 🗄️ 数据管理
- **通用数据库编辑器**: 支持任意表结构的 CRUD 操作
- **多数据库支持**: SQLite、MySQL、PostgreSQL
- **灵活映射**: 可选择任意数据表作为点位表和路径表
- **实时同步**: 画布视图与表格视图数据双向同步

### 📋 智能功能
- **模板系统**: 8种布局模式，支持保存、应用、分享
- **布局算法**: 网格、力导向、圆形等自动排列算法
- **路径生成**: 最近邻连接、完全连通图等智能路径生成
- **数据导出**: Excel/CSV格式导出，完美支持中文

### 📱 跨平台支持
- **桌面端**: Windows、Linux、macOS
- **Web端**: 响应式设计，支持现代浏览器
- **移动端**: PWA支持，平板设备优化

## 🏗️ 技术架构

### 后端
- **语言**: Go 1.21+
- **框架**: Gin + GORM
- **架构**: 领域驱动设计 (DDD)
- **数据库**: SQLite/MySQL/PostgreSQL

### 前端
| 版本 | 技术栈 | 访问路径 | 状态 |
|------|--------|----------|------|
| **经典版** | HTML + JS + Konva.js | `/app` | ✅ 稳定 |
| **现代版** | React + TS + Vite | `/app/new` | ✅ 推荐 |

**现代版特性:**
- React + TypeScript
- Vite 极速构建
- TailwindCSS 现代UI
- Zustand 状态管理
- React Query 数据获取

## 🚀 快速开始

### 演示版 (最快体验)
```bash
# 下载并启动
go run cmd/demo/main.go

# 访问应用
open http://localhost:8080
```

### Docker 部署 (推荐)
```bash
# 克隆项目
git clone <your-repo-url>
cd robot-path-editor

# 启动应用
docker-compose up -d app

# 访问应用
open http://localhost:8080
```

### 开发环境
```bash
# 启动后端
go run cmd/server/main.go

# 启动前端 (可选，用于开发)
cd frontend
npm install && npm run dev

# 访问地址
# 现代前端: http://localhost:5173
# 经典前端: http://localhost:8080/app
```

## 📖 使用指南

### 基础操作
- **创建节点**: 双击空白区域或使用工具栏
- **移动节点**: 拖拽节点到目标位置
- **创建路径**: Shift+点击两个节点
- **删除元素**: 选中后按Delete键
- **撤销重做**: Ctrl+Z / Ctrl+Y

### 模板功能
- **保存模板**: 点击"模板"按钮 → "保存当前为模板"
- **应用模板**: 选择模板并应用到画布
- **模板分类**: 工厂、仓库、实验室等场景分类

### 数据导出
- **导出节点**: 点击"导出"按钮 → 选择CSV/Excel格式
- **导出路径**: 批量导出路径数据
- **完整导出**: 导出所有数据用于备份

## 🔧 配置说明

### 数据库配置
```yaml
# configs/config.yaml
database:
  type: "sqlite"    # sqlite/mysql/postgresql
  dsn: "./data/robot-path-editor.db"

server:
  host: "0.0.0.0"
  port: 8080

logger:
  level: "info"
```

### 环境变量
```bash
export DATABASE_TYPE=sqlite
export DATABASE_DSN=./data/robot-path-editor.db
export SERVER_PORT=8080
```

## 📚 文档

- [📋 详细部署指南](docs/DEPLOYMENT.md)
- [📊 API 接口文档](docs/API.md)
- [🐛 故障排除指南](docs/DEPLOYMENT.md#故障排除)

## 🛠️ 开发

### 构建项目
```bash
# Linux/macOS
./scripts/build.sh

# Windows
./scripts/build.bat
```

### 运行测试
```bash
go test ./... -v -cover
```

### 代码检查
```bash
golangci-lint run
```

## 📊 性能指标

- **响应时间**: < 100ms (95%ile)
- **节点支持**: 10,000+ 节点
- **路径支持**: 50,000+ 路径
- **渲染帧率**: 60 FPS

## 🤝 贡献

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目基于 MIT 许可证开源 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - HTTP Web框架
- [GORM](https://gorm.io/) - Go ORM库
- [Konva.js](https://konvajs.org/) - 2D Canvas渲染引擎
- [React](https://reactjs.org/) - 前端框架

---

<div align="center">

**如果这个项目对你有帮助，请给它一个 ⭐️ Star！**

[🐛 报告Bug](https://github.com/your-org/robot-path-editor/issues) |
[✨ 请求功能](https://github.com/your-org/robot-path-editor/issues) |
[💬 参与讨论](https://github.com/your-org/robot-path-editor/discussions)

</div>