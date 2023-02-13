package common

import (
	goflag "flag"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
	"net"
	"os"
	"time"
	"zbx.io/mockgpu/pkg/para"
)

func Dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	c, err := grpc.Dial(unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func NewGoFlagSet() *goflag.FlagSet {
	fs := goflag.NewFlagSet(os.Args[0], goflag.ExitOnError)
	fs.BoolVar(&para.DebugMode, "debug", false, "debug mode")
	klog.InitFlags(fs)
	return fs
}
