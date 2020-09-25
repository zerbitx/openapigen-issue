package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"github.com/zerbitx/openapigen-issue/openapi"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	done := make(chan struct{})
	callCount := 0

	// just echo back the body that was sent
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		if callCount++; callCount == 2 {
			defer close(done)
		}

		body, err := ioutil.ReadAll(req.Body)
		
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		fmt.Println("Received: ", string(body))
	})

	addr := "localhost:8000"
	go http.ListenAndServe(addr, mux)

	openapi.NewAPIClient(&openapi.Configuration{
		Host:   addr,
		Scheme: "http",
	}).GoApi.OmitEmpty(context.Background(), openapi.ZeroTypes{
		EmptyString: "",
		FalseBool:   false,
		ZeroInt:     0,
	})

	var buf bytes.Buffer
	buf.Write([]byte(`{"emptyString": "", "falseBool": false, zeroInt: 0}`))
	req, _ := http.NewRequest(http.MethodPost,
		"http://"+addr+"/omit/empty",
		&buf)

	(&http.Client{}).Do(req)

	<-done
}
