package koala


// status code
const ( 
    RESPONSE_STATUS_CODE_CONTINUE     = 100  // Continue
    RESPONSE_STATUS_CODE_OK           = 200  // OK
    RESPONSE_STATUS_CODE_CREATED      = 201  // Created
    RESPONSE_STATUS_CODE_LOWSS        = 250  // Low on Storage Space
    RESPONSE_STATUS_CODE_MULTIPLECHOICES = 300 // Multiple Choices
    RESPONSE_STATUS_CODE_MOVEDPERMANENTLY= 301 // Moved Permanently
    RESPONSE_STATUS_CODE_MOVEDTEMPORARILY= 302 // Moved Temporarily
    RESPONSE_STATUS_CODE_SEEOTHER     = 303 // See Other
    RESPONSE_STATUS_CODE_NOTMODIFIED  = 304 // Not Modified
    RESPONSE_STATUS_CODE_USEPROXY     = 305 // Use Proxy
    RESPONSE_STATUS_CODE_BADREQUEST   = 400 // Bad Request
    RESPONSE_STATUS_CODE_UNAUTHORIZED = 401 // Unauthorized
    RESPONSE_STATUS_CODE_PAYMENTREQUIRED = 402 // Payment Required
    RESPONSE_STATUS_CODE_FORBIDDEN    = 403 // Forbidden
    RESPONSE_STATUS_CODE_NOTFOUND     = 404 // Not Found
    RESPONSE_STATUS_CODE_NOTALLOWED   = 405 // Method Not Allowed
    RESPONSE_STATUS_CODE_NOTACCEPTABLE= 406 // Not Acceptable
    RESPONSE_STATUS_CODE_PROXYAUTHREQUIRED = 407 // Proxy Authentication Required
    RESPONSE_STATUS_CODE_REQUESTTIMEOUT    = 408 // Request Time-out
    RESPONSE_STATUS_CODE_GONE              = 410 // Gone
    RESPONSE_STATUS_CODE_LENGTHREQUIRED    = 411 // Length Required
    RESPONSE_STATUS_CODE_PRECONDITIONFAILED= 412 // Precondition Failed
    RESPONSE_STATUS_CODE_REQENTITYTL       = 413 // Request Entity Too Large
    RESPONSE_STATUS_CODE_REQURLTL          = 414 // Request-URI Too Large
    RESPONSE_STATUS_CODE_UNSUPPORTEDMEDIA  = 415 // Unsupported Media Type
    RESPONSE_STATUS_CODE_PARAMNOTUNDERSTOOD= 451 // Parameter Not Understood
    RESPONSE_STATUS_CODE_CONFERENCENOTFOUND= 452 // Conference Not Found
    RESPONSE_STATUS_CODE_NOTENOUGHBANDWIDTH= 453 // Not Enough Bandwidth
    RESPONSE_STATUS_CODE_SESSIONNOTFOUND   = 454 // Session Not Found
    RESPONSE_STATUS_CODE_METHODNOTVALIDSTATE=455 // Method Not Valid in This State
    RESPONSE_STATUS_CODE_HEADERNOTVALIDRES = 456 // Header Field Not Valid for Resource
    RESPONSE_STATUS_CODE_INVALIDRANGE      = 457 // Invalid Range
    RESPONSE_STATUS_CODE_PARAMREADONLY     = 458 // Parameter Is Read-Only
    RESPONSE_STATUS_CODE_AGGREGATENOTALLOWED = 459 // Aggregate operation not allowed
    RESPONSE_STATUS_CODE_ONLYAGGREGATEALLOWED= 460 // Only aggregate operation allowed
    RESPONSE_STATUS_CODE_UNSUPPORTEDTRANSOIRT= 461 // Unsupported transport
    RESPONSE_STATUS_CODE_DESTINATIONUNREACHABLE = 462 // Destination unreachable
    RESPONSE_STATUS_CODE_KEYMANGEMENTFAILURE    = 463 // Key management Failure
    RESPONSE_STATUS_CODE_INTERNALSERVERERROR    = 500 // Internal Server Error
    RESPONSE_STATUS_CODE_NOTIMPLEMENTED         = 501 // Not Implemented
    RESPONSE_STATUS_CODE_BADGATEWAY             = 502 // Bad Gateway
    RESPONSE_STATUS_CODE_SERVICEUNAVAILABLE     = 503 // Service Unavailable
    RESPONSE_STATUS_CODE_GATEWAYTIMEOUT         = 504 // Gateway Time-out
    RESPONSE_STATUS_CODE_RTSPVERNOTSUPPORTED    = 505 // RTSP Version not supported
    RESPONSE_STATUS_CODE_OPTIONNOTSUPPORTED     = 551 // Option not supported
)

type Response interface{
    Recv()
    Write( []byte ) error
    NotFound()
    String() string
}

type BaseResponse struct {}

func (baseResponse *BaseResponse) Recv() {}

func (baseResponse *BaseResponse) String() string {
    return ""
}