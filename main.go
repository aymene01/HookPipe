package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gliderlabs/ssh"
	"github.com/teris-io/shortid"
	gossh "golang.org/x/crypto/ssh"
)

const privateKeyPath = "keys/test_id_rsa"

var clients sync.Map

type HTTPHandler struct{}

func (h *HTTPHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ch, ok := clients.Load(id)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id not found"))
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "error reading request body")
		log.Printf("error: %v", err)
		return
	}
	defer r.Body.Close()

	ch.(chan string) <- string(b)
}

func startHTTPServer() error {
	router := http.NewServeMux()
	handler := &HTTPHandler{}

	router.HandleFunc("/{id}/*", handler.handleWebhook)

	return http.ListenAndServe(":5000", router)
}

func startSSHServer() error {
	generateMockSSHKeys()

	sshPort := ":2222"
	handler := &SSHHandler{}

	server := &ssh.Server{
		Addr:    sshPort,
		Handler: handler.handleSSHSession,
		ServerConfigCallback: func(ctx ssh.Context) *gossh.ServerConfig {
			cfg := &gossh.ServerConfig{
				ServerVersion: "SSH-2.0-sendit",
			}

			cfg.Ciphers = []string{"aes128-gcm@openssh.com"}

			return cfg
		},
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		},
	}

	b, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	signer, err := gossh.ParsePrivateKey(b)
	if err != nil {
		log.Fatal(err)
	}

	server.AddHostKey(signer)

	return server.ListenAndServe()
}

func main() {
	go startSSHServer()
	startHTTPServer()
}

type SSHHandler struct{}

func (h *SSHHandler) handleSSHSession(session ssh.Session) {
	id := shortid.MustGenerate()
	webhookURL := "http://hookpipe.com/" + id
	session.Write([]byte(webhookURL + "\n"))

	respChan := make(chan string)

	clients.Store(id, respChan)

	for data := range respChan {
		session.Write([]byte(data + "\n"))
	}
}
