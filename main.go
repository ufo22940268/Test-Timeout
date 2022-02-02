package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {

	timeoutDuration := 10 * time.Second
	client := &http.Client{
		Timeout: timeoutDuration,

		//Uncomment following block to reproduce this error.
		//Transport: &http.Transport{
		//	DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		//		conn, err := net.DialTimeout(network, addr, timeoutDuration)
		//		if err != nil {
		//			return nil, err
		//		}
		//		err = conn.SetDeadline(time.Now().Add(timeoutDuration))
		//		if err != nil {
		//			return nil, err
		//		}
		//		return conn, nil
		//	},
		//	ResponseHeaderTimeout: timeoutDuration,
		//},
	}

	var wg sync.WaitGroup
	for i := 0; i < 300; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := request(client)
			if err != nil {
				fmt.Printf("error: %+v", err)
				return
			}
			fmt.Println("handle ", i)
		}(i)
	}
	wg.Wait()

	fmt.Println("finished")
}

func request(client *http.Client) error {
	resp, err := client.Get("https://cnbeta.com")
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body)[:10])
	return nil
}
