package device

import (
	"fmt"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

// ExtendDevice 对k8s device对象的扩展结构。mockgpu device等特定device需要额外一些特殊字段进行存储
type ExtendDevice struct {
	pluginapi.Device
	Memory uint64 // 每个mockgpu的显存大小
}

func MockDevices() []*ExtendDevice {
	var extDevices []*ExtendDevice
	for i := 0; i < 2; i++ {
		d := ExtendDevice{}
		d.ID = fmt.Sprintf("%d", i)
		d.Memory = uint64(8000)
		extDevices = append(extDevices, &d)
	}
	return extDevices
}

// ConvertDeviceType 当有client向server请求设备列表时，需要返回设备列表。 由于结构体内部保存的是extenddevice类型，需转换为k8s device原生类型
func ConvertDeviceType(extdevs []*ExtendDevice) []*pluginapi.Device {
	var res []*pluginapi.Device
	for _, dev := range extdevs {
		res = append(res, &pluginapi.Device{
			ID:       dev.ID,
			Health:   dev.Health,
			Topology: nil, // 原生k8s device没有memory等属性，此处不需要上报
		})
	}
	return res
}
