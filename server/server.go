package main

import (
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"
	"net/http"
	"io/ioutil"

	"github.com/gyf19/yar-go/yar"
	"log"
)

type Fetch struct{}

type ArgsIn struct {
	Urls []string
}

type ArgsOut struct {
	Err       error
	Time      time.Duration
	Responses []Response
}

type Response struct {
	Url   string
	Time  time.Duration
	Bytes int
	Err   error
}

func (f *Fetch) MultipleRequest(in *ArgsIn, out *ArgsOut) error {
	if len(in.Urls) == 0 {
		out.Err = fmt.Errorf("you need to expesify the url one or more")
		return nil
	}

	chn := make(chan *Response, len(in.Urls))

	start := time.Now()
	for _, url := range in.Urls {
		go Request(url, chn)
	}

	for range in.Urls {
		res := <-chn
		out.Responses = append(out.Responses, *res)
	}

	out.Time = time.Since(start)
	return nil
}

func Request(url string, chanOut chan *Response) error {
	output := &Response{Url: url}
	defer func(output *Response) {
		chanOut <- output
	}(output)

	log.Printf("Requested to [%s]\n", url)
	start := time.Now()
	resp, err := http.Get(output.Url)
	secs := time.Since(start)

	output.Time = secs
	if err != nil {
		output.Err = err
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		output.Err = err
		return nil
	}

	output.Bytes = len(body)
	return nil
}

var worker = runtime.NumCPU()

func main() {
	runtime.GOMAXPROCS(worker)

	var server = yar.NewServer()
	server.Register(new(Fetch))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		listener, err := net.Listen("tcp", ":12345")
		if err != nil {
			fmt.Println(err)
			return
		}
		wg.Done()
		server.Accept(listener)
	}()

	log.Println("Running...")
	wg.Wait()
	fmt.Scanln()
}
