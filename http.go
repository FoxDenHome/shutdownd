package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

type shutdownHandler struct {
	logger   debug.Log
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *shutdownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		w.Write([]byte("POST only"))
		return
	}

	user, pw, ok := r.BasicAuth()
	if !ok || user != h.Username || pw != h.Password {
		w.Header().Add("WWW-Authenticate", "Basic realm=shutdownd")
		w.WriteHeader(401)
		w.Write([]byte("Please authenticate"))
		return
	}

	switch r.URL.Path {
	case "/shutdown":
		h.logger.Info(1, "Shutdown initiated")
		err := exec.Command("shutdown", "-s", "-f", "-t", "60").Run()
		if err != nil {
			h.logger.Error(1, fmt.Sprintf("Shutdown start error: %v", err))
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	case "/abort":
		h.logger.Info(1, "Shutdown aborted")
		err := exec.Command("shutdown", "-a").Run()
		if err != nil {
			h.logger.Error(1, fmt.Sprintf("Shutdown abort error: %v", err))
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

func (h *shutdownHandler) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	defer h.logger.Close()

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	cfilePath, err := getConfigPath()
	if err != nil {
		h.logger.Error(1, fmt.Sprintf("Could not locate config.json: %v", err))
		return
	}

	cfile, err := os.Open(cfilePath)
	if err != nil {
		h.logger.Error(1, fmt.Sprintf("Could not load config.json: %v", err))
		return
	}

	cfileDecoder := json.NewDecoder(cfile)
	err = cfileDecoder.Decode(h)
	if err != nil {
		h.logger.Error(1, fmt.Sprintf("Could not decode config.json: %v", err))
		return
	}

	if h.Username == "" || h.Password == "" {
		h.logger.Error(1, "Invalid auth configuration")
		return
	}

	server := http.Server{
		Addr:    ":6666",
		Handler: h,
	}

	go func() {
		for {
			c := <-r
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				server.Close()
				return
			}
		}
	}()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	h.logger.Info(1, "ShutdownD listening")

	err = server.ListenAndServe()
	if err != nil {
		h.logger.Error(1, fmt.Sprintf("Could not listen on HTTP: %v", err))
		return
	}
	return
}
