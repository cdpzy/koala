package http

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"time"

	log "github.com/Sirupsen/logrus"

	"fmt"
)

// KoalaController 控制器
type KoalaController struct {
	Name          string             // 控制器名称
	MethodName    string             // 调用方法名称
	Req           *Request           // 输入
	Resp          *Response          // 输出
	Res           Result             // 输出资源
	Controllers   *ControllerManager //
	Type          *ControllerType    //
	MethodType    *MethodType        //
	Action        string             // 动作
	Params        *Params            // 所有参数
	AppController interface{}        // 当前执行的Controller
}

// SetCookie 设置Cookie
func (controller *KoalaController) SetCookie(cookie *http.Cookie) {
	http.SetCookie(controller.Resp.Out, cookie)
}

// Render 渲染
func (controller *KoalaController) Render() Result {
	return nil
}

// Request 获取Request
func (controller *KoalaController) Request() *Request {
	return controller.Req
}

// Response 获取Response
func (controller *KoalaController) Response() *Response {
	return controller.Resp
}

// Result ..
func (controller *KoalaController) Result() Result {
	return controller.Res
}

// RenderError 错误信息渲染
func (controller *KoalaController) RenderError(status int, err error) Result {
	controller.setStatusIfNil(status)
	return ErrorResult{Error: err}
}

// RenderText 文本信息渲染
func (controller *KoalaController) RenderText(text string, objs ...interface{}) Result {
	finalText := text
	if len(objs) > 0 {
		finalText = fmt.Sprintf(text, objs...)
	}
	return &RenderTextResult{finalText}
}

// setStatusIfNil 检查状态
func (controller *KoalaController) setStatusIfNil(status int) {
	if controller.Resp.Status == 0 {
		controller.Resp.Status = status
	}
}

// SetResult 设置
func (controller *KoalaController) SetResult(r Result) {
	controller.Res = r
}

// SetAction 设置控制器操作
func (controller *KoalaController) SetAction(controllerName, methodName string) error {
	controller.Type = controller.Controllers.Get(controllerName)
	if controller.Type == nil {
		return errors.New("failed to find controller " + controllerName)
	}

	controller.MethodType = controller.Type.Method(methodName)
	if controller.MethodType == nil {
		return errors.New("failed to find method " + methodName)
	}

	controller.Name = controller.Type.Type.Name()
	controller.MethodName = controller.MethodType.Name
	controller.Action = controller.Name + "." + controller.MethodName
	appControllerPtr := reflect.New(controller.Type.Type)
	appController := appControllerPtr.Elem()
	cValue := reflect.ValueOf(controller)
	for _, index := range controller.Type.ControllerIndexes {
		appController.FieldByIndex(index).Set(cValue)
	}

	controller.AppController = appControllerPtr.Interface()
	return nil
}

// NotFound 404
func (controller *KoalaController) NotFound(msg string, objs ...interface{}) Result {
	finalText := msg
	if len(objs) > 0 {
		finalText = fmt.Sprintf(msg, objs...)
	}

	return controller.RenderError(http.StatusNotFound, errors.New(finalText))
}

// Forbidden 403
func (controller *KoalaController) Forbidden(msg string, objs ...interface{}) Result {
	finalText := msg
	if len(objs) > 0 {
		finalText = fmt.Sprintf(msg, objs...)
	}

	return controller.RenderError(http.StatusForbidden, errors.New(finalText))
}

// RenderFile returns a file, either displayed inline or downloaded
func (controller *KoalaController) RenderFile(file *os.File, delivery ContentDisposition) Result {
	controller.setStatusIfNil(http.StatusOK)

	var (
		modtime       = time.Now()
		fileInfo, err = file.Stat()
	)
	if err != nil {
		log.Error("RenderFile error:", err)
		return nil
	}

	if fileInfo != nil {
		modtime = fileInfo.ModTime()
	}

	return controller.RenderBinary(file, filepath.Base(file.Name()), delivery, modtime)
}

// RenderBinary bytes
func (controller *KoalaController) RenderBinary(memfile io.Reader, filename string, delivery ContentDisposition, modtime time.Time) Result {
	controller.setStatusIfNil(http.StatusOK)

	return &BinaryResult{
		Reader:   memfile,
		Name:     filename,
		Delivery: delivery,
		Length:   -1, // http.ServeContent gets the length itself unless memfile is a stream.
		ModTime:  modtime,
	}
}

// RenderJon JSON渲染
func (controller *KoalaController) RenderJSON(o interface{}, pretty bool) Result {
	controller.setStatusIfNil(http.StatusOK)
	return RenderJSONResult{o, "", pretty}
}

// RenderJSONP JSONP
func (controller *KoalaController) RenderJSONP(callback string, o interface{}, pretty bool) Result {
	controller.setStatusIfNil(http.StatusOK)
	return RenderJSONResult{o, callback, pretty}
}

// RenderJSONText JSONText
func (controller *KoalaController) RenderJSONText(o string) Result {
	controller.setStatusIfNil(http.StatusOK)
	return RenderJSONTextResult{o, ""}
}

// RenderJSONPText JSONText
func (controller *KoalaController) RenderJSONPText(callback string, o string) Result {
	controller.setStatusIfNil(http.StatusOK)
	return RenderJSONTextResult{o, callback}
}

// NewKoalaController 创建控制器
func NewKoalaController(req *Request, resp *Response) *KoalaController {
	return &KoalaController{Req: req, Resp: resp, Params: &Params{}}
}
