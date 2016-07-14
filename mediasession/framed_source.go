package mediasession

type FramedSource interface{
    MaxFrameSize() uint
    Next(buffTo []byte, maxSize uint, afterGettingFunc interface{}, onCloseFunc interface{})
}


type BaseFramedSource struct {}


func (baseFramedSource *BaseFramedSource) MaxFrameSize() uint {
    return 0
}