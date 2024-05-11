package plugin

type Callback interface {
	Plugin
	CallbackRender(result interface{}, args []string) error
}
