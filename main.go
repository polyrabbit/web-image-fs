package main

import (
	"path/filepath"

	"github.com/polyrabbit/imagefs/config"
	"github.com/polyrabbit/imagefs/fs"
	"github.com/polyrabbit/imagefs/webpage"
	"github.com/sirupsen/logrus"
)

func main() {
	if !config.Execute() {
		return
	}
	mountPoint, err := filepath.Abs(config.MountPoint)
	if err != nil {
		logrus.WithError(err).WithField("mountPoint", mountPoint).Fatal("Failed to get abs file path")
		return
	}
	httpClient := webpage.MustNewHTTPClient(config.BaseURL, config.Timeout)
	server := fs.MustMount(mountPoint, httpClient)
	go server.ListenForUnmount()
	logrus.Infof("Mounted to %q, use ctrl+c to terminate.", mountPoint)
	server.Wait()
}
