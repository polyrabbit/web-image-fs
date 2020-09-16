package fs

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/polyrabbit/imagefs/webpage"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/polyrabbit/etcdfs/config"
	"github.com/sirupsen/logrus"
)

type Server struct {
	*fuse.Server
	mountPoint string
}

func MustMount(mountPoint string, client *webpage.HTTPClient, rootDom webpage.DomNode) *Server {
	opts := &fs.Options{
		MountOptions: fuse.MountOptions{
			Options: config.MountOptions,
			Debug:   false,
			FsName:  "imagefs",
		},
	}
	server, err := fs.Mount(mountPoint, NewRoot(client, rootDom), opts)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to mount")
		return nil
	}
	return &Server{
		Server:     server,
		mountPoint: mountPoint,
	}
}

func (s *Server) ListenForUnmount() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	sig := <-c
	logrus.Infof("Got %s signal, unmounting %q...", sig, s.mountPoint)
	err := s.Unmount()
	if err != nil {
		logrus.WithError(err).Errorf("Failed to unmount, try %q manually.", "umount "+s.mountPoint)
	}
	<-c // Double ctrl+c
	logrus.Warn("Force exiting...")
	os.Exit(1)
}
