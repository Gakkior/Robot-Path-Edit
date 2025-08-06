// Package services 布局服务实现
//
// 设计参考：
// - Graphviz的布局算法
// - D3.js的力导向布局
// - Cytoscape的网络布局
//
// 特点：
// 1. 多种布局算法：力导向、层次化、圆形、网格
// 2. 自适应布局：根据节点数量和关系自动调整
// 3. 性能优化：大图布局算法优化
package services

import (
	"context"
	"math"
	"math/rand"
	"time"

	"robot-path-editor/internal/domain"
)

// LayoutService 布局服务接口
type LayoutService interface {
	// 力导向布局
	ApplyForceDirectedLayout(ctx context.Context, nodes []domain.Node, paths []domain.Path, config ForceDirectedConfig) ([]domain.Node, error)

	// 层次化布局
	ApplyHierarchicalLayout(ctx context.Context, nodes []domain.Node, paths []domain.Path, config HierarchicalConfig) ([]domain.Node, error)

	// 圆形布局
	ApplyCircularLayout(ctx context.Context, nodes []domain.Node, config CircularConfig) ([]domain.Node, error)

	// 网格布局
	ApplyGridLayout(ctx context.Context, nodes []domain.Node, config GridConfig) ([]domain.Node, error)
}

// ForceDirectedConfig 力导向布局配置
type ForceDirectedConfig struct {
	Width      float64 `json:"width" default:"1000"`
	Height     float64 `json:"height" default:"800"`
	Iterations int     `json:"iterations" default:"50"`
	SpringK    float64 `json:"spring_k" default:"0.1"`
	RepelK     float64 `json:"repel_k" default:"1000"`
	Damping    float64 `json:"damping" default:"0.9"`
}

// HierarchicalConfig 层次化布局配置
type HierarchicalConfig struct {
	Width       float64 `json:"width" default:"1000"`
	Height      float64 `json:"height" default:"800"`
	LayerHeight float64 `json:"layer_height" default:"100"`
	NodeSpacing float64 `json:"node_spacing" default:"80"`
}

// CircularConfig 圆形布局配置
type CircularConfig struct {
	CenterX float64 `json:"center_x" default:"500"`
	CenterY float64 `json:"center_y" default:"400"`
	Radius  float64 `json:"radius" default:"300"`
}

// GridConfig 网格布局配置
type GridConfig struct {
	Width      float64 `json:"width" default:"1000"`
	Height     float64 `json:"height" default:"800"`
	Columns    int     `json:"columns" default:"5"`
	NodeSpaceX float64 `json:"node_space_x" default:"100"`
	NodeSpaceY float64 `json:"node_space_y" default:"100"`
}

// layoutService 布局服务实现
type layoutService struct {
	rand *rand.Rand
}

// NewLayoutService 创建新的布局服务实例
func NewLayoutService() LayoutService {
	return &layoutService{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ApplyForceDirectedLayout 应用力导向布局算法
func (s *layoutService) ApplyForceDirectedLayout(ctx context.Context, nodes []domain.Node, paths []domain.Path, config ForceDirectedConfig) ([]domain.Node, error) {
	if len(nodes) == 0 {
		return nodes, nil
	}

	// 设置默认配置
	if config.Width <= 0 {
		config.Width = 1000
	}
	if config.Height <= 0 {
		config.Height = 800
	}
	if config.Iterations <= 0 {
		config.Iterations = 50
	}
	if config.SpringK <= 0 {
		config.SpringK = 0.1
	}
	if config.RepelK <= 0 {
		config.RepelK = 1000
	}
	if config.Damping <= 0 {
		config.Damping = 0.9
	}

	// 构建邻接表
	adjacency := make(map[string][]string)
	for _, path := range paths {
		adjacency[string(path.StartNodeID)] = append(adjacency[string(path.StartNodeID)], string(path.EndNodeID))
		adjacency[string(path.EndNodeID)] = append(adjacency[string(path.EndNodeID)], string(path.StartNodeID))
	}

	// 初始化参数
	width, height := 1000.0, 800.0
	k := math.Sqrt((width * height) / float64(len(nodes))) // 理想距离

	// 随机初始位置 (如果节点位置为空)
	updatedNodes := make([]domain.Node, len(nodes))
	for i, node := range nodes {
		updatedNode := node
		if node.Position.X == 0 && node.Position.Y == 0 {
			updatedNode.Position.X = s.rand.Float64() * width
			updatedNode.Position.Y = s.rand.Float64() * height
		}
		updatedNodes[i] = updatedNode
	}

	// 力导向算法迭代
	for iter := 0; iter < config.Iterations; iter++ {
		// 计算斥力
		for i := range updatedNodes {
			fx, fy := 0.0, 0.0

			for j := range updatedNodes {
				if i != j {
					dx := updatedNodes[i].Position.X - updatedNodes[j].Position.X
					dy := updatedNodes[i].Position.Y - updatedNodes[j].Position.Y
					distance := math.Sqrt(dx*dx + dy*dy)

					if distance > 0 {
						repelForce := config.RepelK / (distance * distance)
						fx += (dx / distance) * repelForce
						fy += (dy / distance) * repelForce
					}
				}
			}

			// 计算引力（仅对相邻节点）
			nodeID := string(updatedNodes[i].ID)
			if neighbors, exists := adjacency[nodeID]; exists {
				for _, neighborID := range neighbors {
					// 找到邻居节点的索引
					for j := range updatedNodes {
						if string(updatedNodes[j].ID) == neighborID {
							dx := updatedNodes[j].Position.X - updatedNodes[i].Position.X
							dy := updatedNodes[j].Position.Y - updatedNodes[i].Position.Y
							distance := math.Sqrt(dx*dx + dy*dy)

							if distance > 0 {
								attractForce := config.SpringK * (distance - k)
								fx += (dx / distance) * attractForce
								fy += (dy / distance) * attractForce
							}
							break
						}
					}
				}
			}

			// 应用力和阻尼
			updatedNodes[i].Position.X += fx * config.Damping
			updatedNodes[i].Position.Y += fy * config.Damping

			// 边界检查
			if updatedNodes[i].Position.X < 0 {
				updatedNodes[i].Position.X = 0
			}
			if updatedNodes[i].Position.X > width {
				updatedNodes[i].Position.X = width
			}
			if updatedNodes[i].Position.Y < 0 {
				updatedNodes[i].Position.Y = 0
			}
			if updatedNodes[i].Position.Y > height {
				updatedNodes[i].Position.Y = height
			}
		}

		// 检查上下文是否被取消
		select {
		case <-ctx.Done():
			return updatedNodes, ctx.Err()
		default:
		}
	}

	return updatedNodes, nil
}

// ApplyHierarchicalLayout 应用层次化布局算法
func (s *layoutService) ApplyHierarchicalLayout(ctx context.Context, nodes []domain.Node, paths []domain.Path, config HierarchicalConfig) ([]domain.Node, error) {
	if len(nodes) == 0 {
		return nodes, nil
	}

	// 设置默认配置
	if config.Width <= 0 {
		config.Width = 1000
	}
	if config.Height <= 0 {
		config.Height = 800
	}
	if config.LayerHeight <= 0 {
		config.LayerHeight = 100
	}
	if config.NodeSpacing <= 0 {
		config.NodeSpacing = 80
	}

	// 简化的层次化布局：按节点类型分层
	layers := make(map[domain.NodeType][]domain.Node)
	for _, node := range nodes {
		layers[node.Type] = append(layers[node.Type], node)
	}

	updatedNodes := make([]domain.Node, 0, len(nodes))
	layerIndex := 0

	for _, layerNodes := range layers {
		y := float64(layerIndex) * config.LayerHeight
		nodeSpacing := config.Width / float64(len(layerNodes)+1)

		for i, node := range layerNodes {
			updatedNode := node
			updatedNode.Position.X = float64(i+1) * nodeSpacing
			updatedNode.Position.Y = y
			updatedNodes = append(updatedNodes, updatedNode)
		}
		layerIndex++
	}

	return updatedNodes, nil
}

// ApplyCircularLayout 应用圆形布局算法
func (s *layoutService) ApplyCircularLayout(ctx context.Context, nodes []domain.Node, config CircularConfig) ([]domain.Node, error) {
	if len(nodes) == 0 {
		return nodes, nil
	}

	// 设置默认配置
	if config.CenterX <= 0 {
		config.CenterX = 500
	}
	if config.CenterY <= 0 {
		config.CenterY = 400
	}
	if config.Radius <= 0 {
		config.Radius = 300
	}

	updatedNodes := make([]domain.Node, len(nodes))
	angleStep := 2 * math.Pi / float64(len(nodes))

	for i, node := range nodes {
		angle := float64(i) * angleStep
		updatedNode := node
		updatedNode.Position.X = config.CenterX + config.Radius*math.Cos(angle)
		updatedNode.Position.Y = config.CenterY + config.Radius*math.Sin(angle)
		updatedNodes[i] = updatedNode
	}

	return updatedNodes, nil
}

// ApplyGridLayout 应用网格布局算法
func (s *layoutService) ApplyGridLayout(ctx context.Context, nodes []domain.Node, config GridConfig) ([]domain.Node, error) {
	if len(nodes) == 0 {
		return nodes, nil
	}

	// 设置默认配置
	if config.Width <= 0 {
		config.Width = 1000
	}
	if config.Height <= 0 {
		config.Height = 800
	}
	if config.Columns <= 0 {
		config.Columns = 5
	}
	if config.NodeSpaceX <= 0 {
		config.NodeSpaceX = 100
	}
	if config.NodeSpaceY <= 0 {
		config.NodeSpaceY = 100
	}

	updatedNodes := make([]domain.Node, len(nodes))

	for i, node := range nodes {
		row := i / config.Columns
		col := i % config.Columns

		updatedNode := node
		updatedNode.Position.X = float64(col) * config.NodeSpaceX
		updatedNode.Position.Y = float64(row) * config.NodeSpaceY
		updatedNodes[i] = updatedNode
	}

	return updatedNodes, nil
}
