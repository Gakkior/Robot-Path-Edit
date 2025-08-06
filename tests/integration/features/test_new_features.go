// 测试新实现的导出和模板功能
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

func main() {
	fmt.Println("🧪 测试导出和模板功能")
	fmt.Println("====================================")

	// 等待服务器启动
	fmt.Println("⏳ 等待服务器启动...")
	time.Sleep(2 * time.Second)

	// 测试模板功能
	testTemplateFeatures()

	// 测试前端导出功能（通过访问静态文件）
	testExportJSLoad()

	fmt.Println("\n✅ 所有新功能测试完成！")
}

func testTemplateFeatures() {
	fmt.Println("\n📋 测试模板功能...")

	// 1. 创建模板
	template := map[string]interface{}{
		"name":        "测试工厂布局",
		"description": "用于测试的工厂车间布局模板",
		"category":    "factory",
		"layout_type": "grid",
		"tags":        []string{"测试", "工厂", "网格"},
		"is_public":   false,
		"template_data": map[string]interface{}{
			"nodes": []map[string]interface{}{
				{
					"template_id": "node_1",
					"name":        "起始点",
					"type":        "point",
					"relative_position": map[string]float64{
						"x": 0.1,
						"y": 0.2,
						"z": 0.0,
					},
					"style": map[string]interface{}{
						"color": "#007bff",
						"size":  10.0,
						"shape": "circle",
					},
				},
				{
					"template_id": "node_2",
					"name":        "终点",
					"type":        "point",
					"relative_position": map[string]float64{
						"x": 0.8,
						"y": 0.7,
						"z": 0.0,
					},
					"style": map[string]interface{}{
						"color": "#28a745",
						"size":  10.0,
						"shape": "circle",
					},
				},
			},
			"paths": []map[string]interface{}{
				{
					"template_id":        "path_1",
					"name":               "主路径",
					"type":               "normal",
					"start_node_temp_id": "node_1",
					"end_node_temp_id":   "node_2",
					"direction":          "bidirectional",
					"curve_type":         "linear",
					"style": map[string]interface{}{
						"color": "#6c757d",
						"width": 2.0,
						"style": "solid",
					},
				},
			},
			"canvas_config": map[string]interface{}{
				"width":        1920,
				"height":       1080,
				"zoom":         1.0,
				"grid_enabled": true,
				"grid_size":    20,
			},
		},
	}

	// 创建模板
	templateJSON, _ := json.Marshal(template)
	resp, err := http.Post(baseURL+"/templates", "application/json", bytes.NewBuffer(templateJSON))
	if err != nil {
		fmt.Printf("❌ 创建模板失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ 创建模板响应: %s\n", string(body))

	// 2. 列出模板
	resp, err = http.Get(baseURL + "/templates")
	if err != nil {
		fmt.Printf("❌ 列出模板失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("✅ 列出模板响应: %s\n", string(body))

	// 3. 获取公开模板
	resp, err = http.Get(baseURL + "/templates/public")
	if err != nil {
		fmt.Printf("❌ 获取公开模板失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("✅ 公开模板响应: %s\n", string(body))

	// 4. 搜索模板
	resp, err = http.Get(baseURL + "/templates/search?q=工厂")
	if err != nil {
		fmt.Printf("❌ 搜索模板失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("✅ 搜索模板响应: %s\n", string(body))

	// 5. 按分类获取模板
	resp, err = http.Get(baseURL + "/templates/category/factory")
	if err != nil {
		fmt.Printf("❌ 按分类获取模板失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("✅ 分类模板响应: %s\n", string(body))

	// 6. 获取模板统计
	resp, err = http.Get(baseURL + "/templates/stats")
	if err != nil {
		fmt.Printf("❌ 获取模板统计失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("✅ 模板统计响应: %s\n", string(body))
}

func testExportJSLoad() {
	fmt.Println("\n📤 测试导出功能JavaScript文件...")

	// 测试导出JS文件是否可访问
	resp, err := http.Get("http://localhost:8080/static/export.js")
	if err != nil {
		fmt.Printf("❌ 导出JS文件访问失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 1000 { // 检查文件是否有合理的大小
			fmt.Printf("✅ 导出JS文件加载成功，大小: %d 字节\n", len(body))

			// 检查是否包含关键功能
			content := string(body)
			features := []string{
				"DataExporter",
				"exportNodesAsCSV",
				"exportPathsAsCSV",
				"exportNodesAsExcel",
				"exportPathsAsExcel",
				"UTF-8",
			}

			for _, feature := range features {
				if contains(content, feature) {
					fmt.Printf("✅ 包含功能: %s\n", feature)
				} else {
					fmt.Printf("❌ 缺少功能: %s\n", feature)
				}
			}
		} else {
			fmt.Printf("❌ 导出JS文件大小异常: %d 字节\n", len(body))
		}
	} else {
		fmt.Printf("❌ 导出JS文件HTTP状态异常: %d\n", resp.StatusCode)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) && (s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsAt(s, substr))))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
