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
			for i = 0; i < 10000; i++ {
				c.Register(c.NewAutoID(), &Client{})
			}

			time.Sleep(time.Second * 1)
		}
	}()

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
				c.Unregister(fmt.Sprint(i))
			}

			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		for {
			c.Iterator(func(id string, c *Client) bool {
				a := c.GetFlag()
				fmt.Println("I:", id, a)
				return true
			})

			time.Sleep(time.Second * 1)
		}
	}()

	time.Sleep(time.Second * 20)
}

func TestFlag(t *testing.T) {
	var i int64
	c := &Client{}
	c.ID = "111"
	go func() {
		for {
			for i = 0; i < 10000; i++ {
				c.SetFlag(FlagClientAuthorized)
			}

			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		for {
			for i = 0; i < 10000; i++ {
				c.SetFlag(FlagClientEncrypt)
			}

			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		for {
			for i = 0; i < 10000; i++ {
				c.SetFlag(FlagClientKickedOut)
			}

			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		for {
			for i = 0; i < 10000; i++ {
				a := c.GetFlag()
				fmt.Println("I:", a)
			}

			time.Sleep(time.Second * 1)
		}
	}()

	time.Sleep(time.Second * 20)
}
