package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"
)

type headersList []string

func (h *headersList) String() string {
	return fmt.Sprint(*h)
}

func (h *headersList) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func main() {
	var method string
	var headers headersList = headersList{"User-Agent: gocurl/1.0"}
	var data string
	var tlsKeyPath string
	var waitTime string
	var url string
	var err error

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-X":
			i++
			method = os.Args[i]
		case "-H":
			i++
			headers.Set(os.Args[i])
		case "-d":
			i++
			data = os.Args[i]
		case "-p", "-path":
			i++
			tlsKeyPath = os.Args[i]
		case "-w", "-wait":
			i++
			waitTime = os.Args[i]
		default:
			url = os.Args[i]
		}
	}

	if url == "" {
		fmt.Println("You must specify a URL")
		return
	}

	if method == "" {
		method = "GET"
	}

	waitTimeBeforeRequest := 0
	if waitTime != "" {
		waitTimeBeforeRequest, err = strconv.Atoi(waitTime)
		if err != nil {
			fmt.Println("You must provide a number for wait time, actual value: ", waitTime)
			return
		}
	}

	client := &http.Client{}

	if tlsKeyPath != "" {
		// 创建一个文件，如果存在就覆盖， 打开方式为：读写，权限为：0666
		file, err := os.OpenFile(tlsKeyPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Printf("create tls key file failed, err: %s", err)
			return
		}
		defer file.Close()
		tlsConfig := &tls.Config{
			KeyLogWriter: file,
		}

		client.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	var reqData *bytes.Buffer = bytes.NewBuffer(([]byte)(""))

	if method != "GET" && data != "" {
		reqData = bytes.NewBuffer([]byte(data))
	}

	req, err := http.NewRequest(method, url, reqData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, header := range headers {
		headerParts := strings.Split(header, ":")
		if len(headerParts) != 2 {
			fmt.Println("Invalid header format. Use Key:Value format.")
			return
		}

		req.Header.Set(strings.TrimSpace(headerParts[0]), strings.TrimSpace(headerParts[1]))
	}
	reqInfo, _ := httputil.DumpRequest(req, true)
	fmt.Printf("Request:\n%s\n", reqInfo)

	if waitTimeBeforeRequest > 0 {
		fmt.Printf("Wait %d seconds before request, to set the tlsKeyLog\n", waitTimeBeforeRequest)
		time.Sleep(time.Duration(waitTimeBeforeRequest) * time.Second)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %s\n", resp.Status)
	for k, v := range resp.Header {
		fmt.Printf("%s: %s\n", k, v)
	}

	//respBody, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("Response Body: %s\n", respBody)

}
