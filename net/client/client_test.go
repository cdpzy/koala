package client

import "testing"
import "time"
import "fmt"

func TestConcurrentAccess(t *testing.T) {
	c := NewClientManager()
	var i int64
	go func() {
		for {
			for i = 0; i < 10000; i++ {
				c.Register(c.NewAutoID(), &Client{})
			}

			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		for {
			for i = 0; i <= 10000; i++ {
				c.Unregister(i)
			}

			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		for {
			c.Iterator(func(id int64, c *Client) bool {
				fmt.Println("I:", id)
				return true
			})

			time.Sleep(time.Second * 1)
		}
	}()

	time.Sleep(time.Second * 60)
}
