package watchdog

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"runtime"

	"github.com/containerd/containerd/v2/plugins"
	"github.com/containerd/plugin"
	"github.com/containerd/plugin/registry"
	"github.com/coreos/go-systemd/v22/daemon"
)

func init() {
	registry.Register(&plugin.Registration{
		Type:     plugins.WatchdogPlugin,
		ID:       "software-watchdog",
		Requires: []plugin.Type{},
		InitFn: func(ic *plugin.InitContext) (interface{}, error) {
			if runtime.GOOS == "windows" {
				return nil, fmt.Errorf("host windows does not support watchdog: %w", plugin.ErrSkipPlugin)
			}
			watchdogUsec := os.Getenv("WATCHDOG_USEC")
			fmt.Println("WATCHDOG_USEC:", watchdogUsec)

			if watchdogUsec != "" {
				return nil, fmt.Errorf("no watchdog interval is configured: %w", plugin.ErrSkipPlugin)
			} 
			// start a go routine that notifies watchdog
			watchdogInterval, err := strconv.Atoi(watchdogUsec)
			if err != nil {
				fmt.Println("Error converting WATCHDOG_USEC to integer:", err)
			}

			fmt.Printf("WATCHDOG_USEC is set to %d microseconds.\n", watchdogInterval)

			notificationInterval := time.Duration(watchdogInterval/2) * time.Microsecond
			// Start a Go routine to periodically notify systemd
			notifySystemd(notificationInterval)
			return &service{}, nil
		},
	})
}

type service struct {
}

func notifySystemd(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			// Notify systemd that the service is still alive
			ack, err := daemon.SdNotify(false, daemon.SdNotifyWatchdog)
			if err != nil {
				fmt.Println("WATCHDOG ERRROR - ", err)
			}
			fmt.Println("Sent watchdog notification -", ack)
		}
	}()
}
