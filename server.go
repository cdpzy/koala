package koala

type HandlerFunc func( Request, Response )

type Server interface{
    HandlerFunc( handlerFunc HandlerFunc )
    Serve( addr string ) error
}