package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"github.com/zerbitx/openapigen-issue/openapi"
	"net/http"
	"sync"
)

func main() {
	mux := http.NewServeMux()

	// just print out the body that was sent
	var wg sync.WaitGroup
	wg.Add(2)
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
	  defer wg.Done()

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
	}).GoApi.
		// Explicitly set values, are serialized away here
		OmitEmpty(context.Background(), openapi.ZeroTypes{
		EmptyString: "",
		FalseBool:   false,
		ZeroInt:     0,
	})

	var buf bytes.Buffer
	buf.Write([]byte(`{"emptyString": "", "falseBool": false, zeroInt: 0}`))
	// But it doesn't have to be that way....
	req, _ := http.NewRequest(http.MethodPost,
		"http://"+addr+"/omit/empty",
		&buf)

	(&http.Client{}).Do(req)

	wg.Wait()
}
