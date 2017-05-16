package event

// CallBack 事件回调
type CallBack func(Event)

// Event 事件
type Event interface {
	Name() string                      // 事件名称
	CallBack() CallBack                // 回调函数
	StopPropagation()                  // 禁止事件传播
	IsStopPropagation() bool           //
	SetData(data interface{})          //
	GetData() interface{}              //
	SetParam(string, interface{})      //
	GetParam(string) interface{}       //
	SetParams(map[string]interface{})  //
	GetParams() map[string]interface{} //
}

// DefaultEvent 默认事件
type DefaultEvent struct {
	EventName         string                 //
	EventCallback     CallBack               //
	isStopPropagation bool                   //
	EventData         interface{}            //
	EventParams       map[string]interface{} //
}

func (e *DefaultEvent) Name() string {
	return e.EventName
}

func (e *DefaultEvent) CallBack() CallBack {
	return e.EventCallback
}

func (e *DefaultEvent) StopPropagation() {
	e.isStopPropagation = true
}

func (e *DefaultEvent) IsStopPropagation() bool {
	return e.isStopPropagation
}

func (e *DefaultEvent) SetData(data interface{}) {
	e.EventData = data
}

func (e *DefaultEvent) GetData() interface{} {
	return e.EventData
}

func (e *DefaultEvent) SetParam(k string, v interface{}) {
	e.EventParams[k] = v
}

func (e *DefaultEvent) GetParam(k string) interface{} {
	return e.EventParams[k]
}

func (e *DefaultEvent) SetParams(p map[string]interface{}) {
	e.EventParams = p
}

func (e *DefaultEvent) GetParams() map[string]interface{} {
	return e.EventParams
}
