// 测试新实现的数据库连接和表映射API
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

// 数据库连接请求结构
type CreateConnectionRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// 表映射请求结构
type CreateTableMappingRequest struct {
	ConnectionID string            `json:"connection_id"`
	TableName    string            `json:"table_name"`
	NodeMapping  *NodeTableMapping `json:"node_mapping,omitempty"`
	PathMapping  *PathTableMapping `json:"path_mapping,omitempty"`
}

type NodeTableMapping struct {
	IDField   string `json:"id_field"`
	NameField string `json:"name_field"`
	TypeField string `json:"type_field"`
	XField    string `json:"x_field"`
	YField    string `json:"y_field"`
	ZField    string `json:"z_field"`
}

type PathTableMapping struct {
	IDField        string `json:"id_field"`
	NameField      string `json:"name_field"`
	StartNodeField string `json:"start_node_field"`
	EndNodeField   string `json:"end_node_field"`
	WeightField    string `json:"weight_field"`
}

func main() {
	fmt.Println("🧪 开始测试机器人路径编辑器 API")
	fmt.Println("====================================================")

	// 等待服务器启动
	fmt.Println("⏳ 等待服务器启动...")
	time.Sleep(2 * time.Second)

	// 测试健康检查
	testHealthCheck()

	// 测试数据库连接管理
	testDatabaseConnections()

	// 测试表映射管理
	testTableMappings()

	// 测试数据同步（Mock版本）
	testDataSync()

	// 测试外部表验证
	testTableValidation()

	fmt.Println("\n✅ 所有API测试完成！")
}

func testHealthCheck() {
	fmt.Println("\n📡 测试健康检查...")
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		fmt.Printf("❌ 健康检查失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ 健康检查响应: %s\n", string(body))
}

func testDatabaseConnections() {
	fmt.Println("\n💾 测试数据库连接管理...")

	// 1. 创建数据库连接
	conn := CreateConnectionRequest{
		Name:     "测试MySQL连接",
		Type:     "mysql",
		Host:     "localhost",
		Port:     3306,
		Database: "robot_db",
		Username: "root",
		Password: "password",
	}

	connJSON, _ := json.Marshal(conn)
	resp, err := http.Post(baseURL+"/database/connections", "application/json", bytes.NewBuffer(connJSON))
	if err != nil {
		fmt.Printf("❌ 创建数据库连接失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ 创建数据库连接响应: %s\n", string(body))

	// 2. 列出数据库连接
	resp, err = http.Get(baseURL + "/database/connections")
	if err != nil {
		fmt.Printf("❌ 列出数据库连接失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("✅ 列出数据库连接响应: %s\n", string(body))

	// 3. 测试数据库连接
	resp, err = http.Post(baseURL+"/database/connections/test-id/test", "application/json", nil)
	if err != nil {
		fmt.Printf("❌ 测试数据库连接失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("✅ 测试数据库连接响应: %s\n", string(body))
}

func testTableMappings() {
	fmt.Println("\n🗺️ 测试表映射管理...")

	// 1. 创建节点表映射
	nodeMapping := CreateTableMappingRequest{
		ConnectionID: "test-connection",
		TableName:    "robot_nodes",
		NodeMapping: &NodeTableMapping{
			IDField:   "node_id",
			NameField: "node_name",
			TypeField: "node_type",
			XField:    "pos_x",
			YField:    "pos_y",
			ZField:    "pos_z",
		},
	}

	mappingJSON, _ := json.Marshal(nodeMapping)
	resp, err := http.Post(baseURL+"/mapping/tables", "application/json", bytes.NewBuffer(mappingJSON))
	if err != nil {
		fmt.Printf("❌ 创建节点表映射失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ 创建节点表映射响应: %s\n", string(body))

	// 2. 创建路径表映射
	pathMapping := CreateTableMappingRequest{
		ConnectionID: "test-connection",
		TableName:    "robot_paths",
		PathMapping: &PathTableMapping{
			IDField:        "path_id",
			NameField:      "path_name",
			StartNodeField: "start_node_id",
			EndNodeField:   "end_node_id",
			WeightField:    "weight",
		},
	}

	mappingJSON, _ = json.Marshal(pathMapping)
	resp, err = http.Post(baseURL+"/mapping/tables", "application/json", bytes.NewBuffer(mappingJSON))
	if err != nil {
		fmt.Printf("❌ 创建路径表映射失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("✅ 创建路径表映射响应: %s\n", string(body))

	// 3. 列出表映射
	resp, err = http.Get(baseURL + "/mapping/tables")
	if err != nil {
		fmt.Printf("❌ 列出表映射失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("✅ 列出表映射响应: %s\n", string(body))
}

func testDataSync() {
	fmt.Println("\n🔄 测试数据同步功能...")

	// 1. 同步节点数据
	resp, err := http.Post(baseURL+"/sync/mappings/test-mapping/nodes", "application/json", nil)
	if err != nil {
		fmt.Printf("❌ 同步节点数据失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ 同步节点数据响应: %s\n", string(body))

	// 2. 同步路径数据
	resp, err = http.Post(baseURL+"/sync/mappings/test-mapping/paths", "application/json", nil)
	if err != nil {
		fmt.Printf("❌ 同步路径数据失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("✅ 同步路径数据响应: %s\n", string(body))

	// 3. 全量数据同步
	resp, err = http.Post(baseURL+"/sync/mappings/test-mapping/all", "application/json", nil)
	if err != nil {
		fmt.Printf("❌ 全量数据同步失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("✅ 全量数据同步响应: %s\n", string(body))
}

func testTableValidation() {
	fmt.Println("\n✅ 测试外部表验证...")

	// 验证外部表结构
	url := baseURL + "/sync/validate-table?connection_id=test-connection&table_name=robot_nodes"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ 验证外部表失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ 验证外部表响应: %s\n", string(body))
}
