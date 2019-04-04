package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

type intHandle struct{}

func (h intHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, f := range l {
		fmt.Fprintf(w, "int: %s\n", f.Name)
	}
}

func main() {
	err := http.ListenAndServe(":7777", intHandle{})
	log.Fatal(err)
}
