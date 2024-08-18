package watchdog

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/containerd/containerd/v2/plugins"
	"github.com/containerd/plugin"
	"github.com/containerd/plugin/registry"
	"github.com/coreos/go-systemd/v22/daemon"
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
			watchdogUsec := os.Getenv("WATCHDOG_USEC")

			fmt.Println("WATCHDOG_USEC:", watchdogUsec)

			if watchdogUsec == "" {
				fmt.Println("WATCHDOG_USEC environment variable is not set.")
			} else {
				// start a go routine that notifies watchdog
				watchdogInterval, err := strconv.Atoi(watchdogUsec)
				if err != nil {
					fmt.Println("Error converting WATCHDOG_USEC to integer:", err)

				}

				fmt.Printf("WATCHDOG_USEC is set to %d microseconds.\n", watchdogInterval)

				notificationInterval := time.Duration(watchdogInterval/2) * time.Microsecond
				// Start a Go routine to periodically notify systemd
				notifyDaemon(notificationInterval)
			}

			return &service{}, nil
		},
	})
}

type service struct {
}

func notifyDaemon(interval time.Duration) {
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
