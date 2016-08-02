package main

import (
    "fmt"
    "os"
    "bytes"
    "io"
    "strconv"
    //"encoding/binary"
)

const (
    PREFIX_SEI_NUT uint32 = 39
    SUFFIX_SEI_NUT = 40
)

// import(
//     // "github.com/doublemo/koala"
//     // "log"
//     // "reflect"
//     //"fmt"
// )

//func main() {
    // s := koala.NewRTSPServer()
    // log.Println(s.GetRequest().String())
    // log.Println(s.Serve(koala.SERVER_NET_PROTO_TCP, ":554"))
    // rtsp := koala.NewRTSPServer()
    // rtsp.HandlerFunc(Handler2)
    // go rtsp.Serve(":554")


    // http := koala.NewHTTPServer()
    // http.HandlerFunc(Handler)
    // log.Println(http.Serve(":9106"))
//}

// func Handler( req koala.Request, resp koala.Response ) {
//     method := koala.NewHandleMethod( req, resp )
//     v := reflect.ValueOf(method).MethodByName(req.GetMethod())
//     if !v.IsValid() {
//         resp.NotSupported( koala.AllowedMethod )
//         return
//     }

//     v.Call([]reflect.Value{})f
// }
   

// func Handler2( req koala.Request, resp koala.Response ) {
//     method := koala.NewHandleMethod( req, resp )
//     ch := req.GetInputChan()
//     for {
//         select {
//             case b,ok := <- ch:
//                  if !ok || b != nil {
//                      break
//                  }

//                  v := reflect.ValueOf(method).MethodByName(req.GetMethod())
//                  if !v.IsValid() {
//                      resp.NotSupported( koala.AllowedMethod )
//                     return
//                  }

//                 v.Call([]reflect.Value{})
//         }
//     }
// }




func main() {
    fid, err := os.Open("test.264")
    if err != nil {
        panic(err)
    }

    var (
      m uint32

    )
    for{
        to := make([]byte, 20000)
        _,err := fid.Read(to)
        if err == io.EOF {
            fmt.Println("EOF")
            break
        }
        
        //fmt.Printf("%v\n", to[0])
        reader := bytes.NewReader(to)

        for {
            h := make([]byte, 4)
            _, err = reader.Read(h)
            if err != nil {
                break
            }

            if toUint(h) != 0x00000001 && toUint(h) & 0xFFFFFF00 != 0x00000100 {
                continue
            }

            if toUint(h) & 0xFFFFFF00 == 0x00000100 {
                reader.UnreadByte()
            }
            
            

            h = make([]byte, 4)
             _, err = reader.Read(h)
            if err != nil {
                break
            }

            m = toUint(h)
            m = m >> 24
           // nalRefIdc   := (m & 0x60) >> 5
            nalUnitType := m & 0x1F
           // fmt.Println("IDC/TYPE:", nalRefIdc, "/", nalUnitType)

            if is264SPS(nalUnitType) {
                fmt.Println("SPS")
            } else if is264PPS(nalUnitType) {
                fmt.Println("PPS")
            }
      
        }
           
           //break
    }
}


func toUint( ptr []byte ) uint32 {
    //fmt.Println((ptr[0] << 24 ) | (ptr[1] << 16) | (ptr[2] << 8) | ptr[3])
    //return uint32((ptr[0] << 24 ) | (ptr[1] << 16) | (ptr[2] << 8) | ptr[3])

    return ( uint32(ptr[0]) << 24 ) | ( uint32(ptr[1]) << 16 ) | ( uint32(ptr[2]) << 8 ) | uint32(ptr[3])
}

func toUint16( ptr []byte ) uint16 {
    return uint16(ptr[0] << 8 | ptr[1])
}

func toUint24( ptr []byte ) uint32 {
    return uint32(ptr[0])<<16 | uint32(ptr[1])<<8 | uint32(ptr[2])
}

//二进制转十六进制
func btox(b string) string {
    base, _ := strconv.ParseInt(b, 2, 10)
    return strconv.FormatInt(base, 16)
}
 
//十六进制转二进制
func xtob(x string) string {
    base, _ := strconv.ParseInt(x, 16, 10)
    return strconv.FormatInt(base, 2)
}


func is264SPS( nal_unit_type uint32 ) bool {
    return nal_unit_type == 7
}

func is264PPS( nal_unit_type uint32 ) bool {
    return nal_unit_type == 8
}

func is264VCL( nal_unit_type uint32 ) bool {
    return nal_unit_type <= 5 && nal_unit_type > 0
}

func is264SEI( nal_unit_type uint32 ) bool {
    return nal_unit_type == 6
}

func is264EOF(nal_unit_type uint32) bool {
    return nal_unit_type == 10 || nal_unit_type == 11
}


// only 265
func is265VPS( nal_unit_type uint32 ) bool {
    return nal_unit_type == 32
}

func is265SPS( nal_unit_type uint32 ) bool {
    return nal_unit_type == 33
}

func is265PPS( nal_unit_type uint32 ) bool {
    return nal_unit_type == 34
}

func is265VCL( nal_unit_type uint32 ) bool {
    return nal_unit_type <= 31
}

func is265SEI( nal_unit_type uint32 ) bool {
    return nal_unit_type == PREFIX_SEI_NUT || nal_unit_type == SUFFIX_SEI_NUT
}

func is265EOF(nal_unit_type uint32) bool {
    return nal_unit_type == 36 || nal_unit_type == 37
}

func usuallyBeginsAccessUnit( code int,  nal_unit_type uint32) bool {
    if code == 264 {
        return (nal_unit_type >= 6 && nal_unit_type <= 9) || (nal_unit_type >= 14 && nal_unit_type <= 18)
    }

    return (nal_unit_type >= 32 && nal_unit_type <= 35) || 
           (nal_unit_type == 39)|| (nal_unit_type >= 41 && nal_unit_type <= 44) || 
           (nal_unit_type >= 48 && nal_unit_type <= 55)
}