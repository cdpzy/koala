package http

import (
	"reflect"
	"strings"
)

// ControllerType 控制器类型
type ControllerType struct {
	Type              reflect.Type
	Methods           []*MethodType
	ControllerIndexes [][]int // FieldByIndex to all embedded *Controllers
}

// MethodType 方法类型
type MethodType struct {
	Name           string
	Args           []*MethodArg
	RenderArgNames map[int][]string
	lowerName      string
}

// MethodArg 方法参数
type MethodArg struct {
	Name string
	Type reflect.Type
}

// ControllerManager 控制器管理
type ControllerManager struct {
	controllers map[string]*ControllerType
}

// Method searches for a given exported method (case insensitive)
func (controllerType *ControllerType) Method(name string) *MethodType {
	lowerName := strings.ToLower(name)
	for _, method := range controllerType.Methods {
		if method.lowerName == lowerName {
			return method
		}
	}
	return nil
}

// Register 注册
func (controllerManager *ControllerManager) Register(c interface{}, methods []*MethodType) {
	t := reflect.TypeOf(c)
	elem := t.Elem()

	for _, m := range methods {
		m.lowerName = strings.ToLower(m.Name)
		for _, arg := range m.Args {
			arg.Type = arg.Type.Elem()
		}
	}

	controllerManager.controllers[strings.ToLower(elem.PkgPath()+"/"+elem.Name())] = &ControllerType{
		Type:              elem,
		Methods:           methods,
		ControllerIndexes: controllerManager.findControllerIndexs(elem),
	}
}

// Get 获取路由器
func (controllerManager *ControllerManager) Get(controllerName string) *ControllerType {
	if c, ok := controllerManager.controllers[strings.ToLower(controllerName)]; ok {
		return c
	}
	return nil
}

func (controllerManager *ControllerManager) findControllerIndexs(appControllerType reflect.Type) (indexes [][]int) {
	type nodeType struct {
		val   reflect.Value
		index []int
	}

	appControllerPtr := reflect.New(appControllerType)
	queue := []nodeType{{appControllerPtr, []int{}}}
	for len(queue) > 0 {
		// Get the next value and de-reference it if necessary.
		var (
			node     = queue[0]
			elem     = node.val
			elemType = elem.Type()
		)
		if elemType.Kind() == reflect.Ptr {
			elem = elem.Elem()
			elemType = elem.Type()
		}
		queue = queue[1:]

		// Look at all the struct fields.
		for i := 0; i < elem.NumField(); i++ {
			// If this is not an anonymous field, skip it.
			structField := elemType.Field(i)
			if !structField.Anonymous {
				continue
			}

			fieldValue := elem.Field(i)
			fieldType := structField.Type

			// If it's a Controller, record the field indexes to get here.
			if fieldType == reflect.TypeOf(&KoalaController{}) {
				indexes = append(indexes, append(node.index, i))
				continue
			}

			queue = append(queue,
				nodeType{fieldValue, append(append([]int{}, node.index...), i)})
		}
	}
	return
}

// NewControllerManager /./
func NewControllerManager() *ControllerManager {
	return &ControllerManager{controllers: make(map[string]*ControllerType)}
}
