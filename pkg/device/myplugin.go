package device

// MyPlugin
// Device Plugin应该有两部分功能： 一部分，是作为server，实现kubelet约定的deviceplugin的所有接口。 另一部分，是作为client，向kubelet发起注册自身。
// 故设计client和server两个结构，各自处理各自的行为。
type MyPlugin struct {
	Client *MyPluginClient // 作为client，向kubelet发起注册用
	Server *MyPluginServer // 将实现kubelet约定的server的所有方法
}
