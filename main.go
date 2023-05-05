package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
)

var validUsername = os.Getenv("SHUTDOWND_USERNAME")
var validPassword = os.Getenv("SHUTDOWND_PASSWORD")

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		w.Write([]byte("POST only"))
		return
	}

	user, pw, ok := r.BasicAuth()
	if !ok || user != validUsername || pw != validPassword {
		w.Header().Add("WWW-Authenticate", "Basic realm=shutdownd")
		w.WriteHeader(401)
		w.Write([]byte("Please authenticate"))
		return
	}

	switch r.URL.Path {
	case "/shutdown":
		log.Printf("Shutdown initiated")
		err := exec.Command("shutdown", "-s", "-f", "-t", "60").Run()
		if err != nil {
			log.Printf("Shutdown error: %v", err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	default:
		w.WriteHeader(404)
		w.Write([]byte("Path not mapped"))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	if validUsername == "" || validPassword == "" {
		panic("Must set SHUTDOWND_USERNAME and SHUTDOWND_PASSWORD")
	}
	log.Printf("ShutdownD listening")
	http.ListenAndServe(":6666", http.HandlerFunc(httpHandler))
}
