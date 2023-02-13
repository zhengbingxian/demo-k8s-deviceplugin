package common

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/signal"
	"syscall"
)

// StartMonitor 启用文件监控和信号监控
func StartMonitor(folderPath string) (*fsnotify.Watcher, chan os.Signal, error) {
	watcher, err := newFSWatcher(folderPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create FS watcher: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT) // notify函数将感兴趣的信号(于第二个参数及以后)，转发到channel(第一个参数)。signal包不会为了向c发送信息而阻塞（就是说如果发送时c阻塞了，signal包会直接放弃）：调用者应该保证c有足够的缓存空间可以跟上期望的信号频率。对使用单一信号用于通知的通道，缓存为1就足够了。

	return watcher, sigChan, nil
}

// newFSWatcher 输入多个要监控的文件路径。 如果文件有变动，可以通过watcher.Events事件来捕获
func newFSWatcher(files ...string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		err = watcher.Add(f)
		if err != nil {
			watcher.Close()
			return nil, err
		}
	}

	return watcher, nil
}
