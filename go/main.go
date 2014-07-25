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

var portFlag = flag.Int("port", 3000, "server port")
var sslFlag = flag.Bool("ssl", false, "enable ssl")
var sslKeyFileFlag = flag.String("ssl key", "", "server ssl key file")
var sslCertFileFlag = flag.String("ssl cert", "", "server ssl cert file")

func main() {
	flag.Parse()

	setupLogger()

	gmux := getMux()

	port, err := strconv.Atoi(os.Getenv("GUSHER_PORT"))

	if err != nil {
		port = *portFlag
	}

	ssl, err := strconv.ParseBool(os.Getenv("GUSHER_SSL"))

	if err != nil {
		ssl = *sslFlag
	}

	var keyFile string
	var certFile string

	if ssl {
		keyFile = os.Getenv("GUSHER_SSL_KEY_FILE")
		certFile = os.Getenv("GUSHER_SSL_CERT_FILE")
		if keyFile == "" {
			keyFile = *sslKeyFileFlag
		}
		if certFile == "" {
			certFile = *sslCertFileFlag
		}
	}

	if ssl && keyFile != "" && certFile != "" {
		log.Printf("ssl server starting on port %d", port)
		log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(port), certFile, keyFile, gmux))
	} else {
		log.Printf("server starting on port %d", port)
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), gmux))
	}
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/index.html")
}
