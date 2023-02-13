package device

import (
	"golang.org/x/net/context"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"path"
	"time"
	"zbx.io/mockgpu/pkg/common"
	"zbx.io/mockgpu/pkg/para"
)

type MyPluginClient struct {
}

func (mpc *MyPluginClient) Register() error {
	conn, err := common.Dial(pluginapi.KubeletSocket, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pluginapi.NewRegistrationClient(conn)
	reqt := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(para.MockGpuSocket), // 告诉kubelet，deviceplugin作为server端监听的文件名称。这里省略/var/lib/kubelet/device-plugins/路径
		ResourceName: para.ResourceName,             // 告诉kubelet，当前设备插件管理的设备资源名，将来会显示在k8s中
		Options: &pluginapi.DevicePluginOptions{
			GetPreferredAllocationAvailable: false,
		},
	}
	_, err = client.Register(context.Background(), reqt)
	if err != nil {
		return err
	}
	return nil
}
