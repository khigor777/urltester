package module

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type UrlStatusCode struct {
	fileName string
	timeoutSec int
}

func New(fileName string, timeoutSec int) *UrlStatusCode  {
	return &UrlStatusCode{
		fileName:fileName,
		timeoutSec:timeoutSec,
	}
}

func (u *UrlStatusCode)Scan()  {
	semaphore := make(chan struct{}, 200)
	scaner, _ := readFile(u.fileName)
	startAll := time.Now()
	for scaner.Scan(){
		semaphore <- struct{}{}
		go func() {
			start := time.Now()
			status, err := sendRequest(scaner.Text())
			fmt.Println(status, err, scaner.Text(), time.Since(start))
			<-semaphore
		}()
	}
	fmt.Println(time.Since(startAll))
}

func readFile(name string) (*bufio.Scanner, error) {
	b, err  := ioutil.ReadFile(name)
	return bufio.NewScanner(bytes.NewReader(b)), err
}

func sendRequest(url string)  (int, error){
	url = "http://" + url
	cl := http.Client{
		Timeout: time.Second * 20,
	}
	resp, err := cl.Get(url)

	if err != nil{
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

