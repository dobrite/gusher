package main

import (
	"flag"
	"github.com/dobrite/gusher/go/gusher"
	"log"
	"net/http"
	"os"
	"strconv"
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

var sslFlag = flag.Bool("ssl", false, "enable ssl")
var portFlag = flag.Int("port", 3000, "port")

func main() {
	setupLogger()

	gmux := getMux()

	flag.Parse()

	port, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		port = *portFlag
	}

	keyFile := "server.key"
	certFile := "server.crt"

	log.Printf("server starting on port %d", port)

	if true {
		log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(port), certFile, keyFile, gmux))
	} else {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), gmux))
	}
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/index.html")
}
