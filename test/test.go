package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("HELLO")
    })
	r.HandleFunc("/{rest:.*}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(mux.Vars(req)["rest"])
		fmt.Println(req.URL.Path)
	})
	req, _ := http.NewRequest("GET", "/foo/some/more/things", nil)
	r.ServeHTTP(&httptest.ResponseRecorder{}, req)
	req, _ = http.NewRequest("GET", "/hello/some/more/things", nil)
	r.ServeHTTP(&httptest.ResponseRecorder{}, req)
	req, _ = http.NewRequest("GET", "/hello", nil)
	r.ServeHTTP(&httptest.ResponseRecorder{}, req)
	req, _ = http.NewRequest("GET", "/", nil)
	r.ServeHTTP(&httptest.ResponseRecorder{}, req)
}
