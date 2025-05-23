package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "test handler")
}

func main(){
	http.HandleFunc("/", handler)
	fmt.Println("svr start 80")
	http.ListenAndServe(":80", nil)
}