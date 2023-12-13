package main

import (
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zhouxiaofang/tecoai-device-plugin/pkg/server"
	"gopkg.in/fsnotify.v1"
)

func main() {
	log.Info("tecoai device plugin starting")
	tecoaiSrv := server.NewTecoaiServer()
	go tecoaiSrv.Run()

	// 向 kubelet 注册
	if err := tecoaiSrv.RegisterToKubelet(); err != nil { //步骤一：在kubelet注册自定义设备的unix socket、API version、ResourceName
		log.Fatalf("register to kubelet error: %v", err)
	} else {
		log.Infoln("register to kubelet successfully")
	}

	// 监听 kubelet.sock，一旦创建则重启
	devicePluginSocket := filepath.Join(server.DevicePluginPath, server.KubeletSocket)
	log.Info("device plugin socket name:", devicePluginSocket)
	watcher, err := fsnotify.NewWatcher() //步骤二：创建监听器
	if err != nil {
		log.Error("Failed to created FS watcher.")
		os.Exit(1)
	}
	defer watcher.Close()
	err = watcher.Add(server.DevicePluginPath) //步骤三：创建监听器需要监听的自定义设备的socket路径//步骤三：创建监听器需要监听的自定义设备的socket路径
	if err != nil {
		log.Error("watch kubelet error")
		return
	}
	log.Info("watching kubelet.sock")
	for {
		select {
		case event := <-watcher.Events:
			log.Infof("watch kubelet events: %s, event name: %s, isCreate: %v", event.Op.String(), event.Name, event.Op&fsnotify.Create == fsnotify.Create)
			if event.Name == devicePluginSocket && event.Op&fsnotify.Create == fsnotify.Create {
				time.Sleep(time.Second)
				log.Fatalf("inotify: %s created, restarting.", devicePluginSocket)
			}
		case err := <-watcher.Errors:
			log.Fatalf("inotify: %s", err)
		}
	}
}
