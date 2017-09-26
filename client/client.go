package main

import (
	"fmt"
	"time"
	"log"

	"github.com/gyf19/yar-go/yar"
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

func main() {
	client, err := yar.Dial("tcp", "localhost:12345", "json")
	checkErr(err)

	var args = &ArgsIn{
		Urls: []string{"http://www.google.com/", "http://www.facebook.com/", "http://www.terra.com.br/"},
	}

	var reply = &ArgsOut{}

	err = client.Call("Fetch.MultipleRequest", args, reply)
	checkErr(err)

	if reply.Err != nil {
		log.Fatalf("Replay: %v \n", reply.Err)
		return
	}

	for _, rp := range reply.Responses {
		fmt.Printf(
			"[%v] elapsed time for request [%s] with [%d]\n",
			rp.Time,
			rp.Url,
			rp.Bytes)
	}

	fmt.Printf("[%v] elapsed time.\n", reply.Time)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Ops: %v\n", err)
	}
}
