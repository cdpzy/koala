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
	name              string                 //
	callback          CallBack               //
	isStopPropagation bool                   //
	data              interface{}            //
	params            map[string]interface{} //
}

func (e *DefaultEvent) Name() string {
	return e.name
}

func (e *DefaultEvent) CallBack() CallBack {
	return e.callback
}

func (e *DefaultEvent) StopPropagation() {
	e.isStopPropagation = true
}

func (e *DefaultEvent) IsStopPropagation() bool {
	return e.isStopPropagation
}

func (e *DefaultEvent) SetData(data interface{}) {
	e.data = data
}

func (e *DefaultEvent) GetData() interface{} {
	return e.data
}

func (e *DefaultEvent) SetParam(k string, v interface{}) {
	e.params[k] = v
}

func (e *DefaultEvent) GetParam(k string) interface{} {
	return e.params[k]
}

func (e *DefaultEvent) SetParams(p map[string]interface{}) {
	e.params = p
}

func (e *DefaultEvent) GetParams() map[string]interface{} {
	return e.params
}
