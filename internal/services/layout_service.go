// Package services 布局服务实现
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
	ArrangeNodes(ctx context.Context, algorithm string) (map[string]domain.Position, error)
	ApplyGridLayout(nodes []domain.Node, spacing float64) []domain.Node
	ApplyForceDirectedLayout(nodes []domain.Node, paths []domain.Path, iterations int) []domain.Node
	ApplyCircularLayout(nodes []domain.Node, radius, centerX, centerY float64) []domain.Node
}

type layoutService struct {
	nodeService NodeService
	pathService PathService
	rand        *rand.Rand
}

func NewLayoutService(nodeService NodeService, pathService PathService) LayoutService {
	return &layoutService{
		nodeService: nodeService,
		pathService: pathService,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *layoutService) ArrangeNodes(ctx context.Context, algorithm string) (map[string]domain.Position, error) {
	// 简单实现，后续可以扩展
	return make(map[string]domain.Position), nil
}

// ApplyGridLayout 网格布局
func (s *layoutService) ApplyGridLayout(nodes []domain.Node, spacing float64) []domain.Node {
	if len(nodes) == 0 {
		return nodes
	}

	// 计算网格尺寸
	cols := int(math.Ceil(math.Sqrt(float64(len(nodes)))))

	updatedNodes := make([]domain.Node, len(nodes))
	for i, node := range nodes {
		row := i / cols
		col := i % cols

		updatedNode := node
		updatedNode.Position.X = float64(col)*spacing + 100
		updatedNode.Position.Y = float64(row)*spacing + 100
		updatedNodes[i] = updatedNode
	}

	return updatedNodes
}

// ApplyForceDirectedLayout 力导向布局 (简化版Fruchterman-Reingold算法)
func (s *layoutService) ApplyForceDirectedLayout(nodes []domain.Node, paths []domain.Path, iterations int) []domain.Node {
	if len(nodes) == 0 {
		return nodes
	}

	// 构建邻接关系
	adjacency := make(map[string][]string)
	for _, path := range paths {
		adjacency[string(path.StartNodeID)] = append(adjacency[string(path.StartNodeID)], string(path.EndNodeID))
		adjacency[string(path.EndNodeID)] = append(adjacency[string(path.EndNodeID)], string(path.StartNodeID))
	}

	// 初始化参�?
	width, height := 1000.0, 800.0
	k := math.Sqrt((width * height) / float64(len(nodes))) // 理想距离

	// 随机初始位置 (如果节点位置�?)
	updatedNodes := make([]domain.Node, len(nodes))
	for i, node := range nodes {
		updatedNode := node
		if node.Position.X == 0 && node.Position.Y == 0 {
			updatedNode.Position.X = s.rand.Float64() * width
			updatedNode.Position.Y = s.rand.Float64() * height
		}
		updatedNodes[i] = updatedNode
	}

	// 迭代计算�?
	for iter := 0; iter < iterations; iter++ {
		// 计算每个节点的受�?
		forces := make(map[string]struct{ fx, fy float64 })

		for i := range updatedNodes {
			forces[string(updatedNodes[i].ID)] = struct{ fx, fy float64 }{0, 0}
		}

		// 计算排斥�?(所有节点对之间)
		for i := 0; i < len(updatedNodes); i++ {
			for j := i + 1; j < len(updatedNodes); j++ {
				node1, node2 := &updatedNodes[i], &updatedNodes[j]
				dx := node1.Position.X - node2.Position.X
				dy := node1.Position.Y - node2.Position.Y
				distance := math.Sqrt(dx*dx + dy*dy)

				if distance < 0.01 {
					distance = 0.01 // 避免除零
				}

				// 库仑排斥�?
				repulsiveForce := k * k / distance
				fx := repulsiveForce * dx / distance
				fy := repulsiveForce * dy / distance

				force1 := forces[string(node1.ID)]
				force1.fx += fx
				force1.fy += fy
				forces[string(node1.ID)] = force1

				force2 := forces[string(node2.ID)]
				force2.fx -= fx
				force2.fy -= fy
				forces[string(node2.ID)] = force2
			}
		}

		// 计算吸引�?(连接的节点之�?
		for _, path := range paths {
			var node1, node2 *domain.Node
			for i := range updatedNodes {
				if updatedNodes[i].ID == path.StartNodeID {
					node1 = &updatedNodes[i]
				}
				if updatedNodes[i].ID == path.EndNodeID {
					node2 = &updatedNodes[i]
				}
			}

			if node1 != nil && node2 != nil {
				dx := node2.Position.X - node1.Position.X
				dy := node2.Position.Y - node1.Position.Y
				distance := math.Sqrt(dx*dx + dy*dy)

				if distance > 0.01 {
					// 胡克引力
					attractiveForce := distance * distance / k
					fx := attractiveForce * dx / distance
					fy := attractiveForce * dy / distance

					force1 := forces[string(node1.ID)]
					force1.fx += fx
					force1.fy += fy
					forces[string(node1.ID)] = force1

					force2 := forces[string(node2.ID)]
					force2.fx -= fx
					force2.fy -= fy
					forces[string(node2.ID)] = force2
				}
			}
		}

		// 应用力并更新位置
		temperature := 10.0 * (1.0 - float64(iter)/float64(iterations)) // 温度递减
		for i := range updatedNodes {
			force := forces[string(updatedNodes[i].ID)]

			// 限制最大移动距�?
			displacement := math.Min(math.Sqrt(force.fx*force.fx+force.fy*force.fy), temperature)
			if displacement > 0.01 {
				updatedNodes[i].Position.X += force.fx / displacement * temperature
				updatedNodes[i].Position.Y += force.fy / displacement * temperature
			}

			// 保持在画布范围内
			updatedNodes[i].Position.X = math.Max(50, math.Min(width-50, updatedNodes[i].Position.X))
			updatedNodes[i].Position.Y = math.Max(50, math.Min(height-50, updatedNodes[i].Position.Y))
		}
	}

	return updatedNodes
}

// ApplyCircularLayout 圆形布局
func (s *layoutService) ApplyCircularLayout(nodes []domain.Node, radius, centerX, centerY float64) []domain.Node {
	if len(nodes) == 0 {
		return nodes
	}

	updatedNodes := make([]domain.Node, len(nodes))
	angleStep := 2 * math.Pi / float64(len(nodes))

	for i, node := range nodes {
		angle := float64(i) * angleStep
		updatedNode := node
		updatedNode.Position.X = centerX + radius*math.Cos(angle)
		updatedNode.Position.Y = centerY + radius*math.Sin(angle)
		updatedNodes[i] = updatedNode
	}

	return updatedNodes
}
