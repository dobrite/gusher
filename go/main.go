package main

import (
	"flag"
	"github.com/dobrite/gusher/go/gusher"
	"log"
	"net/http"
	"os"
)

func setupLogger() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func getMux() *http.ServeMux {
	gmux := gusher.NewServeMux("/app", "tester")
	gmux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	gmux.HandleFunc("/", IndexHandler)
	return gmux
}

type sslFlag struct {
	val string
}

func (s *sslFlag) String() string {
	return s.val
}

func (s *sslFlag) Set(val string) error {
	s.val = val
	return nil
}

func (s *sslFlag) IsBoolFlag() bool {
	return true
}

func setupFlags() {
}

func main() {
	setupLogger()
	setupFlags()

	gmux := getMux()

	ssl := &sslFlag{}

	flag.Var(ssl, "ssl", "enable ssl")

	flag.Parse()

	log.Println(ssl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	keyFile := "server.key"
	certFile := "server.crt"

	log.Println("server started on port " + port)
	//log.Fatal(http.ListenAndServe(":"+port, gmux))
	log.Fatal(http.ListenAndServeTLS(":"+port, certFile, keyFile, gmux))
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/index.html")
}
