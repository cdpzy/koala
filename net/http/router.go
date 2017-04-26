package http

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/doublemo/koala/net/http/tree"

	log "github.com/Sirupsen/logrus"
)

// Route 路由
type Route struct {
	Method     string     // 请求方法
	Path       string     // 路径 eg:/app/list/:page/:sort
	Action     string     // 方法
	Controller string     // 控制器 eg: app/controller/DefaultConter
	Params     url.Values // 参数
}

// Router 路由器
type Router struct {
	tree   *tree.Node //
	filter []Filter   // 流程过滤器
}

// NewRoute 创建路由
func NewRoute(method, path, action, controller string) *Route {
	route := &Route{
		Method:     strings.ToUpper(method),
		Path:       path,
		Action:     action,
		Controller: controller,
		Params:     make(url.Values),
	}

	if !strings.HasPrefix(route.Path, "/") {
		log.Errorln("Absolute URL required.")
		return nil
	}

	return route
}

// Register 注册路由
func (router *Router) Register(route *Route) error {
	return router.tree.Add(router.TreePath(route.Method, route.Path), route)
}

// Apply 路由匹配
func (router *Router) Apply(c *KoalaController) {
	if len(router.filter) < 1 {
		c.SetResult(c.RenderError(http.StatusInternalServerError, errors.New("500 SERVER ERROR: filter is nil")))
	}

	router.filter[0](router, c, router.filter[1:])
}

// Find 获取路由
func (router *Router) Find(method, path string) *Route {
	leaf, expansions := router.tree.Find(router.TreePath(method, path))
	if leaf == nil {
		return nil
	}

	route := leaf.Value.(*Route)
	if len(expansions) > 0 {
		route.Params = make(url.Values)
		for i, v := range expansions {
			route.Params[leaf.Wildcards[i]] = []string{v}
		}
	}

	return route
}

// TreePath 路由路径
func (router *Router) TreePath(method, path string) string {
	if method == "*" {
		method = ":METHOD"
	}

	return "/" + strings.ToUpper(method) + path
}

// NewRouter 创建路由
func NewRouter() *Router {
	return &Router{
		tree: tree.New(),
		filter: []Filter{
			PanicFilter,
			RouterFilter,
			ParamsFilter,
			InvokerFilter,
		},
	}
}
