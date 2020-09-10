package config

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	MountPoint   string
	MountOptions []string
	Timeout      time.Duration
	BaseURL      string
	Verbose      bool
	EnablePprof  bool
)

const (
	defaultTimeout      = 10 * time.Second
	defaultPprofAddress = "localhost:9327"
)

var (
	rootCmd = &cobra.Command{
		Use:   fmt.Sprintf("%s [mount-point] [url]", os.Args[0]),
		Short: "Mount web images to local file system - find help/update at https://github.com/polyrabbit/web-image-fs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return cmd.Help()
			}
			MountPoint = args[0]
			BaseURL = args[1]
			return nil
		},
	}
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "15:04:05", FullTimestamp: true})
	rootCmd.Flags().DurationVar(&Timeout, "http-timeout", defaultTimeout, "http request timeout")
	rootCmd.Flags().BoolVar(&EnablePprof, "enable-pprof", false, fmt.Sprintf("enable runtime profiling data via HTTP server. Address is at %q", "http://"+defaultPprofAddress+"/debug/pprof"))
	rootCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.Flags().StringSliceVar(&MountOptions, "mount-options", []string{"nonempty"}, "options are passed as -o string to fusermount")

	rootCmd.Flags().SortFlags = false
	rootCmd.SilenceErrors = true
}

func Execute() bool {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorln(err)
		return false
	}
	if len(MountPoint) == 0 {
		return false
	}

	if Verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if EnablePprof {
		go func() {
			if err := http.ListenAndServe(defaultPprofAddress, nil); err != nil {
				logrus.WithError(err).Error("Failed to serve pprof")
			}
		}()
	}
	return true
}
