package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func main() {
	req, err := http.NewRequest(http.MethodGet, "http://admin:admin@localhost:8761/eureka/apps", nil)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		panic(err)
	}

	clt := http.Client{}
	resp, err := clt.Do(req)
	if err != nil {
		panic(err)
	}

	s, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", s)
}
