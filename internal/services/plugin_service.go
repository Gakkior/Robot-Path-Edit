// Package services 插件系统服务实现
package services

import (
	"context"
	"fmt"
	"plugin"
	"reflect"
	"sync"

	"robot-path-editor/internal/domain"
)

// Plugin 插件接口
type Plugin interface {
	Name() string
	Version() string
	Description() string
	Initialize(ctx context.Context, config map[string]interface{}) error
	Shutdown(ctx context.Context) error
}

// LayoutPlugin 布局插件接口
type LayoutPlugin interface {
	Plugin
	ApplyLayout(nodes []domain.Node, paths []domain.Path, config map[string]interface{}) ([]domain.Node, error)
}

// PathGenerationPlugin 路径生成插件接口
type PathGenerationPlugin interface {
	Plugin
	GeneratePaths(nodes []domain.Node, config map[string]interface{}) ([]domain.Path, error)
}

// DataProcessorPlugin 数据处理插件接口
type DataProcessorPlugin interface {
	Plugin
	ProcessNodes(nodes []domain.Node, config map[string]interface{}) ([]domain.Node, error)
	ProcessPaths(paths []domain.Path, config map[string]interface{}) ([]domain.Path, error)
}

// EventHandler 事件处理器类�?
type EventHandler func(event Event) error

// Event 事件结构
type Event struct {
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// PluginRegistry 插件注册�?
type PluginRegistry struct {
	mu                    sync.RWMutex
	layoutPlugins         map[string]LayoutPlugin
	pathGenerationPlugins map[string]PathGenerationPlugin
	dataProcessorPlugins  map[string]DataProcessorPlugin
	eventHandlers         map[string][]EventHandler
	loadedPlugins         map[string]Plugin
}

// PluginService 插件服务接口
type PluginService interface {
	// 插件生命周期
	LoadPlugin(ctx context.Context, pluginPath string) error
	UnloadPlugin(ctx context.Context, pluginName string) error
	ListPlugins() []PluginInfo
	GetPluginStatus(pluginName string) (PluginStatus, error)

	// 布局插件
	RegisterLayoutPlugin(plugin LayoutPlugin) error
	UnregisterLayoutPlugin(pluginName string) error
	ApplyLayoutPlugin(ctx context.Context, pluginName string, nodes []domain.Node, paths []domain.Path, config map[string]interface{}) ([]domain.Node, error)
	ListLayoutPlugins() []string

	// 路径生成插件
	RegisterPathGenerationPlugin(plugin PathGenerationPlugin) error
	UnregisterPathGenerationPlugin(pluginName string) error
	GeneratePathsWithPlugin(ctx context.Context, pluginName string, nodes []domain.Node, config map[string]interface{}) ([]domain.Path, error)
	ListPathGenerationPlugins() []string

	// 数据处理插件
	RegisterDataProcessorPlugin(plugin DataProcessorPlugin) error
	UnregisterDataProcessorPlugin(pluginName string) error
	ProcessDataWithPlugin(ctx context.Context, pluginName string, nodes []domain.Node, paths []domain.Path, config map[string]interface{}) ([]domain.Node, []domain.Path, error)
	ListDataProcessorPlugins() []string

	// 事件系统
	RegisterEventHandler(eventType string, handler EventHandler) error
	UnregisterEventHandler(eventType string, handlerID string) error
	EmitEvent(event Event) error
	SubscribeToEvents(eventTypes []string) (<-chan Event, error)
}

// PluginInfo 插件信息
type PluginInfo struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Status      PluginStatus           `json:"status"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// PluginStatus 插件状�?
type PluginStatus string

const (
	PluginStatusLoaded   PluginStatus = "loaded"
	PluginStatusActive   PluginStatus = "active"
	PluginStatusError    PluginStatus = "error"
	PluginStatusDisabled PluginStatus = "disabled"
)

// pluginService 插件服务实现
type pluginService struct {
	registry     *PluginRegistry
	eventChannel chan Event
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewPluginService 创建插件服务
func NewPluginService() PluginService {
	ctx, cancel := context.WithCancel(context.Background())
	service := &pluginService{
		registry: &PluginRegistry{
			layoutPlugins:         make(map[string]LayoutPlugin),
			pathGenerationPlugins: make(map[string]PathGenerationPlugin),
			dataProcessorPlugins:  make(map[string]DataProcessorPlugin),
			eventHandlers:         make(map[string][]EventHandler),
			loadedPlugins:         make(map[string]Plugin),
		},
		eventChannel: make(chan Event, 100),
		ctx:          ctx,
		cancel:       cancel,
	}

	// 启动事件处理协程
	go service.eventProcessor()

	return service
}

// LoadPlugin 加载插件 (支持Go plugin系统)
func (s *pluginService) LoadPlugin(ctx context.Context, pluginPath string) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	// 加载Go插件
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("加载插件失败: %v", err)
	}

	// 查找插件入口�?
	symbol, err := p.Lookup("NewPlugin")
	if err != nil {
		return fmt.Errorf("未找到插件入口点 'NewPlugin': %v", err)
	}

	// 检查入口点类型
	newPluginFunc, ok := symbol.(func() Plugin)
	if !ok {
		return fmt.Errorf("插件入口点类型错误，期望: func() Plugin")
	}

	// 创建插件实例
	pluginInstance := newPluginFunc()

	// 初始化插�?
	if err := pluginInstance.Initialize(ctx, nil); err != nil {
		return fmt.Errorf("插件初始化失�? %v", err)
	}

	// 根据插件类型注册
	pluginName := pluginInstance.Name()
	s.registry.loadedPlugins[pluginName] = pluginInstance

	// 检查插件类型并注册到相应的注册�?
	if layoutPlugin, ok := pluginInstance.(LayoutPlugin); ok {
		s.registry.layoutPlugins[pluginName] = layoutPlugin
	}
	if pathPlugin, ok := pluginInstance.(PathGenerationPlugin); ok {
		s.registry.pathGenerationPlugins[pluginName] = pathPlugin
	}
	if dataPlugin, ok := pluginInstance.(DataProcessorPlugin); ok {
		s.registry.dataProcessorPlugins[pluginName] = dataPlugin
	}

	return nil
}

// UnloadPlugin 卸载插件
func (s *pluginService) UnloadPlugin(ctx context.Context, pluginName string) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	plugin, exists := s.registry.loadedPlugins[pluginName]
	if !exists {
		return fmt.Errorf("插件 %s 未加�?, pluginName)
	}

	// 关闭插件
	if err := plugin.Shutdown(ctx); err != nil {
		return fmt.Errorf("插件关闭失败: %v", err)
	}

	// 从所有注册表中移�?
	delete(s.registry.loadedPlugins, pluginName)
	delete(s.registry.layoutPlugins, pluginName)
	delete(s.registry.pathGenerationPlugins, pluginName)
	delete(s.registry.dataProcessorPlugins, pluginName)

	return nil
}

// ListPlugins 列出所有插�?
func (s *pluginService) ListPlugins() []PluginInfo {
	s.registry.mu.RLock()
	defer s.registry.mu.RUnlock()

	var plugins []PluginInfo
	for _, plugin := range s.registry.loadedPlugins {
		pluginType := s.getPluginType(plugin)
		plugins = append(plugins, PluginInfo{
			Name:        plugin.Name(),
			Version:     plugin.Version(),
			Description: plugin.Description(),
			Type:        pluginType,
			Status:      PluginStatusActive, // 简化状态管�?
		})
	}

	return plugins
}

// GetPluginStatus 获取插件状�?
func (s *pluginService) GetPluginStatus(pluginName string) (PluginStatus, error) {
	s.registry.mu.RLock()
	defer s.registry.mu.RUnlock()

	if _, exists := s.registry.loadedPlugins[pluginName]; exists {
		return PluginStatusActive, nil
	}
	return PluginStatusDisabled, fmt.Errorf("插件 %s 未找�?, pluginName)
}

// RegisterLayoutPlugin 注册布局插件
func (s *pluginService) RegisterLayoutPlugin(plugin LayoutPlugin) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	s.registry.layoutPlugins[plugin.Name()] = plugin
	s.registry.loadedPlugins[plugin.Name()] = plugin
	return nil
}

// UnregisterLayoutPlugin 注销布局插件
func (s *pluginService) UnregisterLayoutPlugin(pluginName string) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	delete(s.registry.layoutPlugins, pluginName)
	return nil
}

// ApplyLayoutPlugin 应用布局插件
func (s *pluginService) ApplyLayoutPlugin(ctx context.Context, pluginName string, nodes []domain.Node, paths []domain.Path, config map[string]interface{}) ([]domain.Node, error) {
	s.registry.mu.RLock()
	plugin, exists := s.registry.layoutPlugins[pluginName]
	s.registry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("布局插件 %s 未找�?, pluginName)
	}

	return plugin.ApplyLayout(nodes, paths, config)
}

// ListLayoutPlugins 列出布局插件
func (s *pluginService) ListLayoutPlugins() []string {
	s.registry.mu.RLock()
	defer s.registry.mu.RUnlock()

	var plugins []string
	for name := range s.registry.layoutPlugins {
		plugins = append(plugins, name)
	}
	return plugins
}

// RegisterPathGenerationPlugin 注册路径生成插件
func (s *pluginService) RegisterPathGenerationPlugin(plugin PathGenerationPlugin) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	s.registry.pathGenerationPlugins[plugin.Name()] = plugin
	s.registry.loadedPlugins[plugin.Name()] = plugin
	return nil
}

// UnregisterPathGenerationPlugin 注销路径生成插件
func (s *pluginService) UnregisterPathGenerationPlugin(pluginName string) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	delete(s.registry.pathGenerationPlugins, pluginName)
	return nil
}

// GeneratePathsWithPlugin 使用插件生成路径
func (s *pluginService) GeneratePathsWithPlugin(ctx context.Context, pluginName string, nodes []domain.Node, config map[string]interface{}) ([]domain.Path, error) {
	s.registry.mu.RLock()
	plugin, exists := s.registry.pathGenerationPlugins[pluginName]
	s.registry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("路径生成插件 %s 未找�?, pluginName)
	}

	return plugin.GeneratePaths(nodes, config)
}

// ListPathGenerationPlugins 列出路径生成插件
func (s *pluginService) ListPathGenerationPlugins() []string {
	s.registry.mu.RLock()
	defer s.registry.mu.RUnlock()

	var plugins []string
	for name := range s.registry.pathGenerationPlugins {
		plugins = append(plugins, name)
	}
	return plugins
}

// RegisterDataProcessorPlugin 注册数据处理插件
func (s *pluginService) RegisterDataProcessorPlugin(plugin DataProcessorPlugin) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	s.registry.dataProcessorPlugins[plugin.Name()] = plugin
	s.registry.loadedPlugins[plugin.Name()] = plugin
	return nil
}

// UnregisterDataProcessorPlugin 注销数据处理插件
func (s *pluginService) UnregisterDataProcessorPlugin(pluginName string) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	delete(s.registry.dataProcessorPlugins, pluginName)
	return nil
}

// ProcessDataWithPlugin 使用插件处理数据
func (s *pluginService) ProcessDataWithPlugin(ctx context.Context, pluginName string, nodes []domain.Node, paths []domain.Path, config map[string]interface{}) ([]domain.Node, []domain.Path, error) {
	s.registry.mu.RLock()
	plugin, exists := s.registry.dataProcessorPlugins[pluginName]
	s.registry.mu.RUnlock()

	if !exists {
		return nil, nil, fmt.Errorf("数据处理插件 %s 未找�?, pluginName)
	}

	processedNodes, err := plugin.ProcessNodes(nodes, config)
	if err != nil {
		return nil, nil, fmt.Errorf("节点处理失败: %v", err)
	}

	processedPaths, err := plugin.ProcessPaths(paths, config)
	if err != nil {
		return nil, nil, fmt.Errorf("路径处理失败: %v", err)
	}

	return processedNodes, processedPaths, nil
}

// ListDataProcessorPlugins 列出数据处理插件
func (s *pluginService) ListDataProcessorPlugins() []string {
	s.registry.mu.RLock()
	defer s.registry.mu.RUnlock()

	var plugins []string
	for name := range s.registry.dataProcessorPlugins {
		plugins = append(plugins, name)
	}
	return plugins
}

// RegisterEventHandler 注册事件处理�?
func (s *pluginService) RegisterEventHandler(eventType string, handler EventHandler) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	s.registry.eventHandlers[eventType] = append(s.registry.eventHandlers[eventType], handler)
	return nil
}

// UnregisterEventHandler 注销事件处理�?(简化实�?
func (s *pluginService) UnregisterEventHandler(eventType string, handlerID string) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	// 简化实现：清空该事件类型的所有处理器
	delete(s.registry.eventHandlers, eventType)
	return nil
}

// EmitEvent 发出事件
func (s *pluginService) EmitEvent(event Event) error {
	select {
	case s.eventChannel <- event:
		return nil
	default:
		return fmt.Errorf("事件队列已满")
	}
}

// SubscribeToEvents 订阅事件 (简化实�?
func (s *pluginService) SubscribeToEvents(eventTypes []string) (<-chan Event, error) {
	// 简化实现：返回主事件通道
	// 在生产环境中，应该为每个订阅者创建专门的通道并过滤事件类�?
	return s.eventChannel, nil
}

// 私有方法

// eventProcessor 事件处理�?
func (s *pluginService) eventProcessor() {
	for {
		select {
		case event := <-s.eventChannel:
			s.handleEvent(event)
		case <-s.ctx.Done():
			return
		}
	}
}

// handleEvent 处理事件
func (s *pluginService) handleEvent(event Event) {
	s.registry.mu.RLock()
	handlers, exists := s.registry.eventHandlers[event.Type]
	s.registry.mu.RUnlock()

	if !exists {
		return
	}

	// 并发处理所有处理器
	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h(event); err != nil {
				// 在生产环境中应该记录日志
				fmt.Printf("事件处理失败: %v\n", err)
			}
		}(handler)
	}
}

// getPluginType 获取插件类型
func (s *pluginService) getPluginType(plugin Plugin) string {
	pluginType := reflect.TypeOf(plugin)
	if pluginType.Implements(reflect.TypeOf((*LayoutPlugin)(nil)).Elem()) {
		return "layout"
	}
	if pluginType.Implements(reflect.TypeOf((*PathGenerationPlugin)(nil)).Elem()) {
		return "path_generation"
	}
	if pluginType.Implements(reflect.TypeOf((*DataProcessorPlugin)(nil)).Elem()) {
		return "data_processor"
	}
	return "unknown"
}

// Shutdown 关闭插件服务
func (s *pluginService) Shutdown(ctx context.Context) error {
	s.cancel()

	// 关闭所有插�?
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	for name, plugin := range s.registry.loadedPlugins {
		if err := plugin.Shutdown(ctx); err != nil {
			fmt.Printf("插件 %s 关闭失败: %v\n", name, err)
		}
	}

	close(s.eventChannel)
	return nil
}
