package mediasession

import (
    "testing"
    "runtime"
    "time"
)


func TestSubSessionReferencecounter( t *testing.T ) {
    sess := &ServerMediaSubSession{}

    for i := 1 ; i <= 50 ; i++ {
        if i % 2 == 0 {
            go func(){
                sess.DecrementReferencecounter()
                runtime.Gosched()
            }()
        } else {
            go func(){
                sess.IncrementReferencecounter()
                runtime.Gosched()
            }()
        }
        
    }

    time.Sleep(time.Second)
    t.Log("counter:", sess.GetReferencecounter())
}