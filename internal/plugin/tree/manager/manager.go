package manager

import (
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/NetAuth/NetAuth/internal/plugin/tree/common"
	"github.com/NetAuth/NetAuth/internal/plugin/tree/consumer"
	"github.com/NetAuth/NetAuth/internal/tree"
)

type hookInserter func(string, string) error

// Manager is a mechanism to keep track of all plugins and handle the
// integration with the tree.
type Manager struct {
	plugins map[string]consumer.Ref
	logger  hclog.Logger
}

// New returns a new manager instance
func New() (Manager, error) {
	x := Manager{
		plugins: make(map[string]consumer.Ref),
		logger:  hclog.L().Named("plugin"),
	}

	return x, nil
}

// LoadPlugins loads all plugins either directly from a dynamic
// discovery, or from a statically defined list.
func (m *Manager) LoadPlugins() {
	if viper.GetBool("plugin.loadstatic") {
		for _, p := range viper.GetStringSlice("plugin.list") {
			m.logger.Trace("Loading new plugin", "plugin", p)
			im, err := consumer.New(p)
			if err != nil {
				m.logger.Warn("Error loading plugin", "error", err)
				continue
			}
			if err := im.Init(); err != nil {
				m.logger.Warn("Error initializing plugin", "error", err)
				continue
			}
			m.plugins[p] = im
		}
	}
}

// Shutdown calls shutdown in each plugin and should be called during
// server shutdown to prevent leaking processes.
func (m *Manager) Shutdown() {
	m.logger.Debug("Shutting down plugins")
	for n, p := range m.plugins {
		m.logger.Debug("Plugin shutdown", "plugin", n)
		p.Shutdown()
	}
}

// RegisterEntityHooks handles the generation and registration of
// hooks in the entity subsystem.
func (m *Manager) RegisterEntityHooks() {
	for i := range common.AutoEntityActions {
		action := common.AutoEntityActions[i]
		hc := func(r tree.RefContext) (tree.EntityHook, error) {
			return EntityHook{
				action: action,
				mref:   m,
			}, nil
		}
		m.logger.Trace("Registering EntityHookConstructor", "action", action)
		tree.RegisterEntityHookConstructor(string(action), hc)
	}
}

// ConfigureEntityChains is called with a reference to a hookInserter
// which will insert the named hook into the named chain.
func (m *Manager) ConfigureEntityChains(h hookInserter) {
	for _, pair := range entityChainConfig {
		parts := strings.Split(pair, ":")
		h(parts[1], parts[0])
	}
}

// InvokeEntityProcessing calls ProcessEntity in every plugin.
func (m *Manager) InvokeEntityProcessing(opts common.PluginOpts) (common.PluginResult, error) {
	var res = common.PluginResult{}
	var err error

	for p, r := range m.plugins {
		m.logger.Trace("Calling plugin", "plugin", p, "action", opts.Action)
		res, err = r.ProcessEntity(opts)
		if err != nil {
			return common.PluginResult{}, err
		}
		opts.Entity = &res.Entity
	}
	return res, nil
}

// RegisterGroupHooks handles the generation and registration of
// hooks in the group subsystem.
func (m *Manager) RegisterGroupHooks() {
	for i := range common.AutoGroupActions {
		action := common.AutoGroupActions[i]
		hc := func(r tree.RefContext) (tree.GroupHook, error) {
			return GroupHook{
				action: action,
				mref:   m,
			}, nil
		}
		m.logger.Trace("Registering GroupHookConstructor", "action", action)
		tree.RegisterGroupHookConstructor(string(action), hc)
	}
}

// ConfigureGroupChains is called with a reference to a hookInserter
// which will insert the named hook into the named chain.
func (m *Manager) ConfigureGroupChains(h hookInserter) {
	for _, pair := range groupChainConfig {
		parts := strings.Split(pair, ":")
		h(parts[1], parts[0])
	}
}

// InvokeGroupProcessing calls ProcessGroup in every plugin.
func (m *Manager) InvokeGroupProcessing(opts common.PluginOpts) (common.PluginResult, error) {
	var res = common.PluginResult{}
	var err error

	for p, r := range m.plugins {
		m.logger.Trace("Calling plugin", "plugin", p, "action", opts.Action)
		res, err = r.ProcessGroup(opts)
		if err != nil {
			return common.PluginResult{}, err
		}
		opts.Group = &res.Group
	}
	return res, nil
}
