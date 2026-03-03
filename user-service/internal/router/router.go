package router

type RouterGroup struct {
	AuthRouter AuthRouter
	UserRouter UserRouter
}

var RouterGroupApp = new(RouterGroup)
