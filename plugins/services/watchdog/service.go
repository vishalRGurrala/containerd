package watchdog

import (
	"fmt"

	"github.com/containerd/containerd/v2/plugins"
	"github.com/containerd/plugin"
	"github.com/containerd/plugin/registry"
)

const (
	pluginid = "777Watchdog"
)

func init() {
	fmt.Println("Starting " + pluginid)
	registry.Register(&plugin.Registration{
		Type:     plugins.WatchdogPlugin,
		ID:       "daemon-health",
		Requires: []plugin.Type{},
		InitFn: func(ic *plugin.InitContext) (interface{}, error) {

			return &service{}, nil
		},
	})
}

type service struct {
}
