package device

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"net"
	"os"
	"zbx.io/mockgpu/pkg/para"
)

type MyPluginServer struct {
	server       *grpc.Server
	CacheDevices []*ExtendDevice
}

func (mps *MyPluginServer) NewServer() *MyPluginServer {
	server := grpc.NewServer([]grpc.ServerOption{}...)
	return &MyPluginServer{server: server, CacheDevices: []*ExtendDevice{}}
}

func (mps *MyPluginServer) Start() error {
	sock, err := net.Listen("unix", para.MockGpuSocket)
	if err != nil {
		return err
	}
	pluginapi.RegisterDevicePluginServer(mps.server, mps)
	go func() {
		mps.server.Serve(sock)
	}()
	return nil
}

// Clean 重启或清理时使用，以对unix套接字进行移除。
func (mps *MyPluginServer) Clean() error {
	if err := os.Remove(para.MockGpuSocket); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (mps *MyPluginServer) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	options := &pluginapi.DevicePluginOptions{
		GetPreferredAllocationAvailable: false,
	}
	return options, nil
}

func (mps *MyPluginServer) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	_ = s.Send(&pluginapi.ListAndWatchResponse{Devices: ConvertDeviceType(mps.CacheDevices)}) // 当有client向server请求设备列表时，需要返回设备列表。 而内部保存的是扩展device类型，所以需要转换下
	return nil
}

func (mps *MyPluginServer) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return &pluginapi.PreferredAllocationResponse{}, nil
}

// Allocate 在重启创建过程中，kubelet作为client发起请求，此时deviceplugin可以运行一些设备的特定操作，告诉kubelet进行设置容器的必须的环境变量, 容器必须要挂在哪些文件等。
func (mps *MyPluginServer) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	klog.Infoln("Allocate", reqs.ContainerRequests)
	if len(reqs.ContainerRequests) > 1 {
		return &pluginapi.AllocateResponse{}, errors.New("multiple Container Requests not supported")
	}
	resp := pluginapi.AllocateResponse{}

	cResp := pluginapi.ContainerAllocateResponse{}
	cResp.Envs = make(map[string]string)
	cResp.Envs["ak"] = "12312414"
	cResp.Envs["sk"] = "123124"
	// 这里假设：使用mockgpu的容器，必须要挂在/var/logs/zbx.so文件到容器path下才行。
	cResp.Mounts = append(cResp.Mounts,
		&pluginapi.Mount{HostPath: "/var/logs/zbx.so",
			ContainerPath: "/root/zbx.so",
			ReadOnly:      true})

	resp.ContainerResponses = append(resp.ContainerResponses, &cResp)
	return &resp, nil
}

// PreStartContainer is unimplemented for this plugin
func (mps *MyPluginServer) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}
