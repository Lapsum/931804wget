package main

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Counter struct {
	Count int
}

func (c *Counter) Write(buf []byte) (n int, err error) {
	c.Count += len(buf)
	return len(buf), nil
}

func main() {
	if len(os.Args) <= 1 {
		println("You must give URL")
		os.Exit(1)
	}

	URL := os.Args[1]
	if resp, err := http.Get(URL); err != nil {
		println(err.Error())
		os.Exit(2)

	} else if resp.Status != "200 OK"{
		defer resp.Body.Close()

		println(resp.Status)
		os.Exit(3)

	} else {
		defer resp.Body.Close()

		parts := strings.Split(URL, "/")
		name := parts[len(parts) - 1]

		if file, err := os.Create(name); err != nil {
			println(err.Error())
			os.Exit(4)
		} else {
			defer file.Close()

			counter := Counter{ 0 }

			c1 := make(chan bool)
			c2 := make(chan bool)

			go func() {
				for true {
					time.Sleep(time.Second)
					select {
					case <-c1:
						c2 <- true
						return
					default:
						println(counter.Count, " bytes received")
					}
				}
			}()

			teeReader := io.TeeReader(resp.Body, &counter)
			if _, err := io.Copy(file, teeReader); err != nil {
				println(err.Error())
				os.Exit(5)
			} else {
				println("File is downloaded")
			}

			c1 <- true
			<- c2
		}
	}
}
