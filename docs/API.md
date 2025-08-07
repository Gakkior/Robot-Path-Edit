# 📊 API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **Content-Type**: `application/json`
- **响应格式**: JSON

## 系统接口

### 健康检查
```http
GET /health
```

**响应示例:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## 节点管理

### 获取所有节点
```http
GET /nodes
```

### 创建节点
```http
POST /nodes
Content-Type: application/json

{
  "name": "节点1",
  "type": "normal",
  "position": {
    "x": 100,
    "y": 200,
    "z": 0
  }
}
```

### 更新节点
```http
PUT /nodes/{id}
Content-Type: application/json

{
  "name": "更新的节点名",
  "position": {
    "x": 150,
    "y": 250
  }
}
```

### 删除节点
```http
DELETE /nodes/{id}
```

## 路径管理

### 获取所有路径
```http
GET /paths
```

### 创建路径
```http
POST /paths
Content-Type: application/json

{
  "name": "路径1",
  "from": "node-1",
  "to": "node-2",
  "type": "normal",
  "weight": 1.0
}
```

### 更新路径
```http
PUT /paths/{id}
Content-Type: application/json

{
  "weight": 2.0,
  "type": "bidirectional"
}
```

### 删除路径
```http
DELETE /paths/{id}
```

## 模板管理

### 获取模板列表
```http
GET /templates
```

### 保存为模板
```http
POST /templates/save-as
Content-Type: application/json

{
  "name": "工厂布局模板",
  "description": "标准工厂车间布局",
  "category": "factory",
  "layout_type": "grid"
}
```

### 应用模板
```http
POST /templates/{id}/apply
Content-Type: application/json

{
  "width": 1920,
  "height": 1080
}
```

## 布局算法

### 应用布局算法
```http
POST /layout/apply
Content-Type: application/json

{
  "algorithm": "force-directed"
}
```

**支持的算法:**
- `force-directed`: 力导向布局
- `hierarchical`: 层次布局
- `circular`: 圆形布局
- `grid`: 网格布局

## 路径生成

### 生成最近邻路径
```http
POST /path-generation/nearest-neighbor
```

### 生成完全连通图
```http
POST /path-generation/full-connectivity
```

## 数据库连接

### 获取连接列表
```http
GET /database/connections
```

### 创建数据库连接
```http
POST /database/connections
Content-Type: application/json

{
  "name": "生产数据库",
  "type": "mysql",
  "host": "localhost",
  "port": 3306,
  "database": "robot_paths",
  "username": "robot",
  "password": "password"
}
```

### 测试数据库连接
```http
POST /database/connections/{id}/test
```

## 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "参数错误",
  "error": "字段验证失败"
}
```

## 使用示例

### JavaScript
```javascript
// 获取节点列表
const response = await fetch('/api/v1/nodes');
const data = await response.json();

// 创建节点
const newNode = await fetch('/api/v1/nodes', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    name: '新节点',
    type: 'normal',
    position: { x: 100, y: 200 }
  })
});
```

### curl
```bash
# 获取节点列表
curl http://localhost:8080/api/v1/nodes

# 创建节点
curl -X POST http://localhost:8080/api/v1/nodes \
  -H "Content-Type: application/json" \
  -d '{"name":"节点1","type":"normal","position":{"x":100,"y":200}}'
```