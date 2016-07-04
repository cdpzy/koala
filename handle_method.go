package koala

import (
    "log"
)

type HandleMethod struct{
    r Request
    w Response
}

func (handleMethod *HandleMethod) OPTIONS() {
    log.Println(handleMethod.r.String())
}


func NewHandleMethod( r Request, w Response ) *HandleMethod {
    handle := new(HandleMethod)
    handle.r = r
    handle.w = w
    return handle
}