package para

// global flag parameter
var (
	ResourceName string // 资源名，默认为zbx.com/mockgpu
	DebugMode    bool   // 是否为调试模式，默认false
	NodeName     string // 用于筛选当前node下的pending pod
)

const (
	MockGpuSocket = "/var/lib/kubelet/device-plugins/zbx-mockgpu.sock"
)
