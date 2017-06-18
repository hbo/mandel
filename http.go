package main


import "net/http"
//import "html"
import "log"
import "time"
import "fmt"
import "path" 


type imgHandler struct {
	img  TargetImage
	quit chan bool
}

func theError(w http.ResponseWriter, r *http.Request, text string,code int){
	base := path.Base(r.URL.Path)
	dir := path.Dir(r.URL.Path)
	err := fmt.Sprintf("Error %d:\n%s.%s: %s", code, dir, base, text)
	http.Error(w, err, 404)
}

func handleJSFile(w http.ResponseWriter, r *http.Request) {

	base := path.Base(r.URL.Path)
	
	switch  {
	case  base == "imgload.js"  :
		http.ServeFile(w, r, base)
	default:
		theError(w,r,"Hamwa nich", 404)
	}
	
}

func handleHTMLFile(w http.ResponseWriter, r *http.Request) {
	base := path.Base(r.URL.Path)
	
	switch  {
	case  base == "/"  :
		http.ServeFile(w, r, "bal.html")
	default:
		theError(w,r,"Hamwa nich", 404)
	}
}


func handleFile (w http.ResponseWriter, r *http.Request) {
	dir := path.Dir(r.URL.Path)
	switch dir {
	case "/js":
		handleJSFile(w,r)
	case "/":
		handleHTMLFile(w,r)
	default:
		theError(w,r,"Hamwa nich", 404)
	}	
}

func getImg (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/tmp/col.png")
}

func (i imgHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

//	compute prefix of r.URL.Path

	dir := path.Dir(r.URL.Path)
	
	switch  dir {
	case "/":
		handleFile(w, r)
	case "/quit":
		i.quit <- true
	case "/js":
		handleFile(w, r)
	case "/getimg":
		getImg(w, r)
	default:
		theError(w,r,"Hamwa nich", 404)
	}
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

