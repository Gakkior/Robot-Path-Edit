// Package services 鎻掍欢绯荤粺鏈嶅姟瀹炵幇
package services

import (
	"context"
	"fmt"
	"plugin"
	"reflect"
	"sync"

	"robot-path-editor/internal/domain"
)

// Plugin 鎻掍欢鎺ュ彛
type Plugin interface {
	Name() string
	Version() string
	Description() string
	Initialize(ctx context.Context, config map[string]interface{}) error
	Shutdown(ctx context.Context) error
}

// LayoutPlugin 甯冨眬鎻掍欢鎺ュ彛
type LayoutPlugin interface {
	Plugin
	ApplyLayout(nodes []domain.Node, paths []domain.Path, config map[string]interface{}) ([]domain.Node, error)
}

// PathGenerationPlugin 璺緞鐢熸垚鎻掍欢鎺ュ彛
type PathGenerationPlugin interface {
	Plugin
	GeneratePaths(nodes []domain.Node, config map[string]interface{}) ([]domain.Path, error)
}

// DataProcessorPlugin 鏁版嵁澶勭悊鎻掍欢鎺ュ彛
type DataProcessorPlugin interface {
	Plugin
	ProcessNodes(nodes []domain.Node, config map[string]interface{}) ([]domain.Node, error)
	ProcessPaths(paths []domain.Path, config map[string]interface{}) ([]domain.Path, error)
}

// EventHandler 浜嬩欢澶勭悊鍣ㄧ被鍨?
type EventHandler func(event Event) error

// Event 浜嬩欢缁撴瀯
type Event struct {
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// PluginRegistry 鎻掍欢娉ㄥ唽琛?
type PluginRegistry struct {
	mu                    sync.RWMutex
	layoutPlugins         map[string]LayoutPlugin
	pathGenerationPlugins map[string]PathGenerationPlugin
	dataProcessorPlugins  map[string]DataProcessorPlugin
	eventHandlers         map[string][]EventHandler
	loadedPlugins         map[string]Plugin
}

// PluginService 鎻掍欢鏈嶅姟鎺ュ彛
type PluginService interface {
	// 鎻掍欢鐢熷懡鍛ㄦ湡
	LoadPlugin(ctx context.Context, pluginPath string) error
	UnloadPlugin(ctx context.Context, pluginName string) error
	ListPlugins() []PluginInfo
	GetPluginStatus(pluginName string) (PluginStatus, error)

	// 甯冨眬鎻掍欢
	RegisterLayoutPlugin(plugin LayoutPlugin) error
	UnregisterLayoutPlugin(pluginName string) error
	ApplyLayoutPlugin(ctx context.Context, pluginName string, nodes []domain.Node, paths []domain.Path, config map[string]interface{}) ([]domain.Node, error)
	ListLayoutPlugins() []string

	// 璺緞鐢熸垚鎻掍欢
	RegisterPathGenerationPlugin(plugin PathGenerationPlugin) error
	UnregisterPathGenerationPlugin(pluginName string) error
	GeneratePathsWithPlugin(ctx context.Context, pluginName string, nodes []domain.Node, config map[string]interface{}) ([]domain.Path, error)
	ListPathGenerationPlugins() []string

	// 鏁版嵁澶勭悊鎻掍欢
	RegisterDataProcessorPlugin(plugin DataProcessorPlugin) error
	UnregisterDataProcessorPlugin(pluginName string) error
	ProcessDataWithPlugin(ctx context.Context, pluginName string, nodes []domain.Node, paths []domain.Path, config map[string]interface{}) ([]domain.Node, []domain.Path, error)
	ListDataProcessorPlugins() []string

	// 浜嬩欢绯荤粺
	RegisterEventHandler(eventType string, handler EventHandler) error
	UnregisterEventHandler(eventType string, handlerID string) error
	EmitEvent(event Event) error
	SubscribeToEvents(eventTypes []string) (<-chan Event, error)
}

// PluginInfo 鎻掍欢淇℃伅
type PluginInfo struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Status      PluginStatus           `json:"status"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// PluginStatus 鎻掍欢鐘舵€?
type PluginStatus string

const (
	PluginStatusLoaded   PluginStatus = "loaded"
	PluginStatusActive   PluginStatus = "active"
	PluginStatusError    PluginStatus = "error"
	PluginStatusDisabled PluginStatus = "disabled"
)

// pluginService 鎻掍欢鏈嶅姟瀹炵幇
type pluginService struct {
	registry     *PluginRegistry
	eventChannel chan Event
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewPluginService 鍒涘缓鎻掍欢鏈嶅姟
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

	// 鍚姩浜嬩欢澶勭悊鍗忕▼
	go service.eventProcessor()

	return service
}

// LoadPlugin 鍔犺浇鎻掍欢 (鏀寔Go plugin绯荤粺)
func (s *pluginService) LoadPlugin(ctx context.Context, pluginPath string) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	// 鍔犺浇Go鎻掍欢
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("鍔犺浇鎻掍欢澶辫触: %v", err)
	}

	// 鏌ユ壘鎻掍欢鍏ュ彛鐐?
	symbol, err := p.Lookup("NewPlugin")
	if err != nil {
		return fmt.Errorf("鏈壘鍒版彃浠跺叆鍙ｇ偣 'NewPlugin': %v", err)
	}

	// 妫€鏌ュ叆鍙ｇ偣绫诲瀷
	newPluginFunc, ok := symbol.(func() Plugin)
	if !ok {
		return fmt.Errorf("鎻掍欢鍏ュ彛鐐圭被鍨嬮敊璇紝鏈熸湜: func() Plugin")
	}

	// 鍒涘缓鎻掍欢瀹炰緥
	pluginInstance := newPluginFunc()

	// 鍒濆鍖栨彃浠?
	if err := pluginInstance.Initialize(ctx, nil); err != nil {
		return fmt.Errorf("鎻掍欢鍒濆鍖栧け璐? %v", err)
	}

	// 鏍规嵁鎻掍欢绫诲瀷娉ㄥ唽
	pluginName := pluginInstance.Name()
	s.registry.loadedPlugins[pluginName] = pluginInstance

	// 妫€鏌ユ彃浠剁被鍨嬪苟娉ㄥ唽鍒扮浉搴旂殑娉ㄥ唽琛?
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

// UnloadPlugin 鍗歌浇鎻掍欢
func (s *pluginService) UnloadPlugin(ctx context.Context, pluginName string) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	plugin, exists := s.registry.loadedPlugins[pluginName]
	if !exists {
		return fmt.Errorf("鎻掍欢 %s 鏈姞杞?, pluginName)
	}

	// 鍏抽棴鎻掍欢
	if err := plugin.Shutdown(ctx); err != nil {
		return fmt.Errorf("鎻掍欢鍏抽棴澶辫触: %v", err)
	}

	// 浠庢墍鏈夋敞鍐岃〃涓Щ闄?
	delete(s.registry.loadedPlugins, pluginName)
	delete(s.registry.layoutPlugins, pluginName)
	delete(s.registry.pathGenerationPlugins, pluginName)
	delete(s.registry.dataProcessorPlugins, pluginName)

	return nil
}

// ListPlugins 鍒楀嚭鎵€鏈夋彃浠?
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
			Status:      PluginStatusActive, // 绠€鍖栫姸鎬佺鐞?
		})
	}

	return plugins
}

// GetPluginStatus 鑾峰彇鎻掍欢鐘舵€?
func (s *pluginService) GetPluginStatus(pluginName string) (PluginStatus, error) {
	s.registry.mu.RLock()
	defer s.registry.mu.RUnlock()

	if _, exists := s.registry.loadedPlugins[pluginName]; exists {
		return PluginStatusActive, nil
	}
	return PluginStatusDisabled, fmt.Errorf("鎻掍欢 %s 鏈壘鍒?, pluginName)
}

// RegisterLayoutPlugin 娉ㄥ唽甯冨眬鎻掍欢
func (s *pluginService) RegisterLayoutPlugin(plugin LayoutPlugin) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	s.registry.layoutPlugins[plugin.Name()] = plugin
	s.registry.loadedPlugins[plugin.Name()] = plugin
	return nil
}

// UnregisterLayoutPlugin 娉ㄩ攢甯冨眬鎻掍欢
func (s *pluginService) UnregisterLayoutPlugin(pluginName string) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	delete(s.registry.layoutPlugins, pluginName)
	return nil
}

// ApplyLayoutPlugin 搴旂敤甯冨眬鎻掍欢
func (s *pluginService) ApplyLayoutPlugin(ctx context.Context, pluginName string, nodes []domain.Node, paths []domain.Path, config map[string]interface{}) ([]domain.Node, error) {
	s.registry.mu.RLock()
	plugin, exists := s.registry.layoutPlugins[pluginName]
	s.registry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("甯冨眬鎻掍欢 %s 鏈壘鍒?, pluginName)
	}

	return plugin.ApplyLayout(nodes, paths, config)
}

// ListLayoutPlugins 鍒楀嚭甯冨眬鎻掍欢
func (s *pluginService) ListLayoutPlugins() []string {
	s.registry.mu.RLock()
	defer s.registry.mu.RUnlock()

	var plugins []string
	for name := range s.registry.layoutPlugins {
		plugins = append(plugins, name)
	}
	return plugins
}

// RegisterPathGenerationPlugin 娉ㄥ唽璺緞鐢熸垚鎻掍欢
func (s *pluginService) RegisterPathGenerationPlugin(plugin PathGenerationPlugin) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	s.registry.pathGenerationPlugins[plugin.Name()] = plugin
	s.registry.loadedPlugins[plugin.Name()] = plugin
	return nil
}

// UnregisterPathGenerationPlugin 娉ㄩ攢璺緞鐢熸垚鎻掍欢
func (s *pluginService) UnregisterPathGenerationPlugin(pluginName string) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	delete(s.registry.pathGenerationPlugins, pluginName)
	return nil
}

// GeneratePathsWithPlugin 浣跨敤鎻掍欢鐢熸垚璺緞
func (s *pluginService) GeneratePathsWithPlugin(ctx context.Context, pluginName string, nodes []domain.Node, config map[string]interface{}) ([]domain.Path, error) {
	s.registry.mu.RLock()
	plugin, exists := s.registry.pathGenerationPlugins[pluginName]
	s.registry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("璺緞鐢熸垚鎻掍欢 %s 鏈壘鍒?, pluginName)
	}

	return plugin.GeneratePaths(nodes, config)
}

// ListPathGenerationPlugins 鍒楀嚭璺緞鐢熸垚鎻掍欢
func (s *pluginService) ListPathGenerationPlugins() []string {
	s.registry.mu.RLock()
	defer s.registry.mu.RUnlock()

	var plugins []string
	for name := range s.registry.pathGenerationPlugins {
		plugins = append(plugins, name)
	}
	return plugins
}

// RegisterDataProcessorPlugin 娉ㄥ唽鏁版嵁澶勭悊鎻掍欢
func (s *pluginService) RegisterDataProcessorPlugin(plugin DataProcessorPlugin) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	s.registry.dataProcessorPlugins[plugin.Name()] = plugin
	s.registry.loadedPlugins[plugin.Name()] = plugin
	return nil
}

// UnregisterDataProcessorPlugin 娉ㄩ攢鏁版嵁澶勭悊鎻掍欢
func (s *pluginService) UnregisterDataProcessorPlugin(pluginName string) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	delete(s.registry.dataProcessorPlugins, pluginName)
	return nil
}

// ProcessDataWithPlugin 浣跨敤鎻掍欢澶勭悊鏁版嵁
func (s *pluginService) ProcessDataWithPlugin(ctx context.Context, pluginName string, nodes []domain.Node, paths []domain.Path, config map[string]interface{}) ([]domain.Node, []domain.Path, error) {
	s.registry.mu.RLock()
	plugin, exists := s.registry.dataProcessorPlugins[pluginName]
	s.registry.mu.RUnlock()

	if !exists {
		return nil, nil, fmt.Errorf("鏁版嵁澶勭悊鎻掍欢 %s 鏈壘鍒?, pluginName)
	}

	processedNodes, err := plugin.ProcessNodes(nodes, config)
	if err != nil {
		return nil, nil, fmt.Errorf("鑺傜偣澶勭悊澶辫触: %v", err)
	}

	processedPaths, err := plugin.ProcessPaths(paths, config)
	if err != nil {
		return nil, nil, fmt.Errorf("璺緞澶勭悊澶辫触: %v", err)
	}

	return processedNodes, processedPaths, nil
}

// ListDataProcessorPlugins 鍒楀嚭鏁版嵁澶勭悊鎻掍欢
func (s *pluginService) ListDataProcessorPlugins() []string {
	s.registry.mu.RLock()
	defer s.registry.mu.RUnlock()

	var plugins []string
	for name := range s.registry.dataProcessorPlugins {
		plugins = append(plugins, name)
	}
	return plugins
}

// RegisterEventHandler 娉ㄥ唽浜嬩欢澶勭悊鍣?
func (s *pluginService) RegisterEventHandler(eventType string, handler EventHandler) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	s.registry.eventHandlers[eventType] = append(s.registry.eventHandlers[eventType], handler)
	return nil
}

// UnregisterEventHandler 娉ㄩ攢浜嬩欢澶勭悊鍣?(绠€鍖栧疄鐜?
func (s *pluginService) UnregisterEventHandler(eventType string, handlerID string) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	// 绠€鍖栧疄鐜帮細娓呯┖璇ヤ簨浠剁被鍨嬬殑鎵€鏈夊鐞嗗櫒
	delete(s.registry.eventHandlers, eventType)
	return nil
}

// EmitEvent 鍙戝嚭浜嬩欢
func (s *pluginService) EmitEvent(event Event) error {
	select {
	case s.eventChannel <- event:
		return nil
	default:
		return fmt.Errorf("浜嬩欢闃熷垪宸叉弧")
	}
}

// SubscribeToEvents 璁㈤槄浜嬩欢 (绠€鍖栧疄鐜?
func (s *pluginService) SubscribeToEvents(eventTypes []string) (<-chan Event, error) {
	// 绠€鍖栧疄鐜帮細杩斿洖涓讳簨浠堕€氶亾
	// 鍦ㄧ敓浜х幆澧冧腑锛屽簲璇ヤ负姣忎釜璁㈤槄鑰呭垱寤轰笓闂ㄧ殑閫氶亾骞惰繃婊や簨浠剁被鍨?
	return s.eventChannel, nil
}

// 绉佹湁鏂规硶

// eventProcessor 浜嬩欢澶勭悊鍣?
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

// handleEvent 澶勭悊浜嬩欢
func (s *pluginService) handleEvent(event Event) {
	s.registry.mu.RLock()
	handlers, exists := s.registry.eventHandlers[event.Type]
	s.registry.mu.RUnlock()

	if !exists {
		return
	}

	// 骞跺彂澶勭悊鎵€鏈夊鐞嗗櫒
	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h(event); err != nil {
				// 鍦ㄧ敓浜х幆澧冧腑搴旇璁板綍鏃ュ織
				fmt.Printf("浜嬩欢澶勭悊澶辫触: %v\n", err)
			}
		}(handler)
	}
}

// getPluginType 鑾峰彇鎻掍欢绫诲瀷
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

// Shutdown 鍏抽棴鎻掍欢鏈嶅姟
func (s *pluginService) Shutdown(ctx context.Context) error {
	s.cancel()

	// 鍏抽棴鎵€鏈夋彃浠?
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	for name, plugin := range s.registry.loadedPlugins {
		if err := plugin.Shutdown(ctx); err != nil {
			fmt.Printf("鎻掍欢 %s 鍏抽棴澶辫触: %v\n", name, err)
		}
	}

	close(s.eventChannel)
	return nil
}
