// Package services 甯冨眬鏈嶅姟瀹炵幇
package services

import (
	"context"
	"math"
	"math/rand"
	"time"

	"robot-path-editor/internal/domain"
)

// LayoutService 甯冨眬鏈嶅姟鎺ュ彛
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
	// 绠€鍗曞疄鐜帮紝鍚庣画鍙互鎵╁睍
	return make(map[string]domain.Position), nil
}

// ApplyGridLayout 缃戞牸甯冨眬
func (s *layoutService) ApplyGridLayout(nodes []domain.Node, spacing float64) []domain.Node {
	if len(nodes) == 0 {
		return nodes
	}

	// 璁＄畻缃戞牸灏哄
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

// ApplyForceDirectedLayout 鍔涘鍚戝竷灞€ (绠€鍖栫増Fruchterman-Reingold绠楁硶)
func (s *layoutService) ApplyForceDirectedLayout(nodes []domain.Node, paths []domain.Path, iterations int) []domain.Node {
	if len(nodes) == 0 {
		return nodes
	}

	// 鏋勫缓閭绘帴鍏崇郴
	adjacency := make(map[string][]string)
	for _, path := range paths {
		adjacency[string(path.StartNodeID)] = append(adjacency[string(path.StartNodeID)], string(path.EndNodeID))
		adjacency[string(path.EndNodeID)] = append(adjacency[string(path.EndNodeID)], string(path.StartNodeID))
	}

	// 鍒濆鍖栧弬鏁?
	width, height := 1000.0, 800.0
	k := math.Sqrt((width * height) / float64(len(nodes))) // 鐞嗘兂璺濈

	// 闅忔満鍒濆浣嶇疆 (濡傛灉鑺傜偣浣嶇疆涓?)
	updatedNodes := make([]domain.Node, len(nodes))
	for i, node := range nodes {
		updatedNode := node
		if node.Position.X == 0 && node.Position.Y == 0 {
			updatedNode.Position.X = s.rand.Float64() * width
			updatedNode.Position.Y = s.rand.Float64() * height
		}
		updatedNodes[i] = updatedNode
	}

	// 杩唬璁＄畻鍔?
	for iter := 0; iter < iterations; iter++ {
		// 璁＄畻姣忎釜鑺傜偣鐨勫彈鍔?
		forces := make(map[string]struct{ fx, fy float64 })

		for i := range updatedNodes {
			forces[string(updatedNodes[i].ID)] = struct{ fx, fy float64 }{0, 0}
		}

		// 璁＄畻鎺掓枼鍔?(鎵€鏈夎妭鐐瑰涔嬮棿)
		for i := 0; i < len(updatedNodes); i++ {
			for j := i + 1; j < len(updatedNodes); j++ {
				node1, node2 := &updatedNodes[i], &updatedNodes[j]
				dx := node1.Position.X - node2.Position.X
				dy := node1.Position.Y - node2.Position.Y
				distance := math.Sqrt(dx*dx + dy*dy)

				if distance < 0.01 {
					distance = 0.01 // 閬垮厤闄ら浂
				}

				// 搴撲粦鎺掓枼鍔?
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

		// 璁＄畻鍚稿紩鍔?(杩炴帴鐨勮妭鐐逛箣闂?
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
					// 鑳″厠寮曞姏
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

		// 搴旂敤鍔涘苟鏇存柊浣嶇疆
		temperature := 10.0 * (1.0 - float64(iter)/float64(iterations)) // 娓╁害閫掑噺
		for i := range updatedNodes {
			force := forces[string(updatedNodes[i].ID)]

			// 闄愬埗鏈€澶хЩ鍔ㄨ窛绂?
			displacement := math.Min(math.Sqrt(force.fx*force.fx+force.fy*force.fy), temperature)
			if displacement > 0.01 {
				updatedNodes[i].Position.X += force.fx / displacement * temperature
				updatedNodes[i].Position.Y += force.fy / displacement * temperature
			}

			// 淇濇寔鍦ㄧ敾甯冭寖鍥村唴
			updatedNodes[i].Position.X = math.Max(50, math.Min(width-50, updatedNodes[i].Position.X))
			updatedNodes[i].Position.Y = math.Max(50, math.Min(height-50, updatedNodes[i].Position.Y))
		}
	}

	return updatedNodes
}

// ApplyCircularLayout 鍦嗗舰甯冨眬
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
