package koala

type Response interface{
    Recv()
    Write( []byte ) error
    String() string
}

type BaseResponse struct {}

func (baseResponse *BaseResponse) Recv() {}

func (baseResponse *BaseResponse) String() string {
    return ""
}