package main


import "net/http"
import "html"
import "log"
import "time"
import "fmt"


type imgHandler struct {
	img  TargetImage
	quit chan bool
}


func (i imgHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	if r.URL.Path == "/quit" {
		i.quit <- true
	}
}


func httpserver (img TargetImage, quit chan bool){
	var i imgHandler

	i.img = img
	i.quit = quit
	
	s := &http.Server{
		Addr:           ":8180",
		Handler:        i,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
	
}

