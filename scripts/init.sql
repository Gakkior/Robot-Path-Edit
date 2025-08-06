-- 机器人路径编辑器数据库初始化脚本
-- 为MySQL创建基本的表结构和示例数据

-- 创建数据库(如果不存在)
CREATE DATABASE IF NOT EXISTS robot_paths CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE robot_paths;

-- 创建节点表(示例)
CREATE TABLE IF NOT EXISTS robot_nodes (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) DEFAULT 'point',
    status VARCHAR(20) DEFAULT 'active',
    x DOUBLE DEFAULT 0,
    y DOUBLE DEFAULT 0,
    z DOUBLE DEFAULT 0,
    color VARCHAR(20) DEFAULT '#007bff',
    size DOUBLE DEFAULT 10.0,
    shape VARCHAR(20) DEFAULT 'circle',
    properties JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    version INT DEFAULT 1
) ENGINE=InnoDB;

-- 创建路径表(示例)  
CREATE TABLE IF NOT EXISTS robot_paths (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) DEFAULT 'normal',
    status VARCHAR(20) DEFAULT 'active',
    start_node_id VARCHAR(36) NOT NULL,
    end_node_id VARCHAR(36) NOT NULL,
    weight DOUBLE DEFAULT 1.0,
    length DOUBLE DEFAULT 0,
    direction VARCHAR(20) DEFAULT 'bidirectional',
    curve_type VARCHAR(20) DEFAULT 'linear',
    color VARCHAR(20) DEFAULT '#6c757d',
    width DOUBLE DEFAULT 2.0,
    style VARCHAR(20) DEFAULT 'solid',
    properties JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    version INT DEFAULT 1,
    INDEX idx_start_node (start_node_id),
    INDEX idx_end_node (end_node_id)
) ENGINE=InnoDB;

-- 插入示例节点数据
INSERT IGNORE INTO robot_nodes (id, name, type, x, y, z, color, size, shape) VALUES
('node-001', '起始点', 'point', 100, 100, 0, '#28a745', 12.0, 'circle'),
('node-002', '工作站A', 'station', 300, 150, 0, '#007bff', 15.0, 'square'),
('node-003', '工作站B', 'station', 500, 200, 0, '#007bff', 15.0, 'square'),
('node-004', '充电点', 'charging', 150, 300, 0, '#ffc107', 12.0, 'triangle'),
('node-005', '终点', 'point', 400, 350, 0, '#dc3545', 12.0, 'circle');

-- 插入示例路径数据
INSERT IGNORE INTO robot_paths (id, name, start_node_id, end_node_id, weight, direction) VALUES
('path-001', '起始到工作站A', 'node-001', 'node-002', 1.0, 'bidirectional'),
('path-002', '工作站A到工作站B', 'node-002', 'node-003', 1.2, 'bidirectional'), 
('path-003', '起始到充电点', 'node-001', 'node-004', 1.5, 'bidirectional'),
('path-004', '工作站B到终点', 'node-003', 'node-005', 1.0, 'bidirectional'),
('path-005', '充电点到终点', 'node-004', 'node-005', 1.3, 'bidirectional');

-- 创建数据库连接表
CREATE TABLE IF NOT EXISTS database_connections (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    host VARCHAR(255),
    port INT,
    database_name VARCHAR(100),
    username VARCHAR(100),
    password VARCHAR(255),
    properties JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 创建表映射配置表
CREATE TABLE IF NOT EXISTS table_mappings (
    id VARCHAR(36) PRIMARY KEY,
    connection_id VARCHAR(36) NOT NULL,
    table_name VARCHAR(100) NOT NULL,
    node_mapping JSON,
    path_mapping JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_connection (connection_id)
) ENGINE=InnoDB;

-- 创建模板表
CREATE TABLE IF NOT EXISTS templates (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    category VARCHAR(50) DEFAULT 'custom',
    tags JSON,
    layout_type VARCHAR(30) NOT NULL,
    layout_config JSON,
    template_data JSON NOT NULL,
    preview JSON,
    usage_count INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    version INT DEFAULT 1,
    labels JSON,
    annotations JSON,
    INDEX idx_category (category),
    INDEX idx_layout_type (layout_type),
    INDEX idx_status (status),
    INDEX idx_public (is_public)
) ENGINE=InnoDB;

-- 插入示例模板数据
INSERT IGNORE INTO templates (id, name, description, category, layout_type, template_data) VALUES
('tpl-001', '工厂车间标准布局', '适用于小型工厂车间的标准机器人路径布局', 'factory', 'grid', 
 '{"nodes":[{"template_id":"node_1","name":"进料点","type":"station","relative_position":{"x":0.1,"y":0.2,"z":0}},{"template_id":"node_2","name":"加工区","type":"station","relative_position":{"x":0.5,"y":0.3,"z":0}},{"template_id":"node_3","name":"出料点","type":"station","relative_position":{"x":0.9,"y":0.7,"z":0}}],"paths":[{"template_id":"path_1","name":"进料路径","start_node_temp_id":"node_1","end_node_temp_id":"node_2"},{"template_id":"path_2","name":"出料路径","start_node_temp_id":"node_2","end_node_temp_id":"node_3"}],"canvas_config":{"width":1920,"height":1080,"zoom":1.0}}'),

('tpl-002', '仓库物流布局', '适用于智能仓库的AGV路径规划', 'warehouse', 'tree',
 '{"nodes":[{"template_id":"node_1","name":"入库口","type":"point","relative_position":{"x":0.1,"y":0.5,"z":0}},{"template_id":"node_2","name":"货架A","type":"storage","relative_position":{"x":0.3,"y":0.3,"z":0}},{"template_id":"node_3","name":"货架B","type":"storage","relative_position":{"x":0.3,"y":0.7,"z":0}},{"template_id":"node_4","name":"出库口","type":"point","relative_position":{"x":0.9,"y":0.5,"z":0}}],"paths":[{"template_id":"path_1","name":"入库路径A","start_node_temp_id":"node_1","end_node_temp_id":"node_2"},{"template_id":"path_2","name":"入库路径B","start_node_temp_id":"node_1","end_node_temp_id":"node_3"},{"template_id":"path_3","name":"出库路径A","start_node_temp_id":"node_2","end_node_temp_id":"node_4"},{"template_id":"path_4","name":"出库路径B","start_node_temp_id":"node_3","end_node_temp_id":"node_4"}],"canvas_config":{"width":1920,"height":1080,"zoom":1.0}}');

-- 创建用户表(为将来的多用户功能准备)
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_nodes_type ON robot_nodes(type);
CREATE INDEX IF NOT EXISTS idx_nodes_status ON robot_nodes(status);
CREATE INDEX IF NOT EXISTS idx_paths_type ON robot_paths(type);
CREATE INDEX IF NOT EXISTS idx_paths_status ON robot_paths(status);

-- 显示创建的表
SHOW TABLES;

-- 显示表结构
DESCRIBE robot_nodes;
DESCRIBE robot_paths;