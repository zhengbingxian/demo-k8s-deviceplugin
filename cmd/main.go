package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"syscall"
	"zbx.io/mockgpu/pkg/common"
	"zbx.io/mockgpu/pkg/device"
	"zbx.io/mockgpu/pkg/para"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mockgpu",
		Short: "mock-gpu程序，for k8s，以device-plugin方式实现",
		Run: func(cmd *cobra.Command, args []string) {
			if err := start(); err != nil {
				klog.Fatal(err)
			}
		},
	}
)

func init() {
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().StringVar(&para.ResourceName, "resource-name", "zbx.com/mockgpu", "resource name")
	rootCmd.PersistentFlags().AddGoFlagSet(common.NewGoFlagSet())
}

func start() error {
	watcher, sigChan, err := common.StartMonitor("/var/lib/kubelet/device-plugins/")
	if err != nil {
		return err
	}
	var plugin = device.MyPlugin{
		Client: &device.MyPluginClient{},
		Server: &device.MyPluginServer{},
	}
	plugin.Server = plugin.Server.NewServer()
	plugin.Server.CacheDevices = device.MockDevices()

restart:
	// as server
	err = plugin.Server.Clean()
	if err != nil {
		klog.Warning("plugin server clean failed")
	}
	plugin.Server.Start()
	if err != nil {
		klog.Warning("plugin server start failed")
	}
	// as client
	plugin.Client.Register()
	if err != nil {
		klog.Warning("plugin client register failed")
	}
events:
	for {
		select {
		case event := <-watcher.Events:
			if event.Name == pluginapi.KubeletSocket && event.Op&fsnotify.Create == fsnotify.Create {
				klog.Infof("inotify: %s created, restarting.", pluginapi.KubeletSocket)
				goto restart
			}
		case err := <-watcher.Errors:
			klog.Infof("inotify: %s", err)
		case s := <-sigChan:
			switch s {
			case syscall.SIGHUP:
				klog.Info("Received SIGHUP, restarting.")
				goto restart
			default:
				klog.Infof("Received signal %v, shutting down.", s)
				break events
			}
		}
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
