package main

import (
	"log"
	"os"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/teris-io/shortid"
	gossh "golang.org/x/crypto/ssh"
)

const privateKeyPath = "keys/test_id_rsa"

func main() {
	generateMockSSHKeys()

	sshPort := ":2222"

	respCh := make(chan string)

	handler := &SSHHandler{
		respCh: respCh,
	}

	go func() {
		time.Sleep(time.Second * 3)
		id, _ := shortid.Generate()
		respCh <- "http://hookpipe.com/" + id

		time.Sleep(time.Second * 10)

		for {
			respCh <- "received data from hook"
			time.Sleep(time.Second * 3)
		}
	}()

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

	log.Fatal(server.ListenAndServe())
}

type SSHHandler struct {
	respCh chan string
}

func (h *SSHHandler) handleSSHSession(session ssh.Session) {
	forwardUrl := session.RawCommand()
	_ = forwardUrl
	resp := <-h.respCh
	session.Write([]byte(resp + "\n"))

	for data := range h.respCh {
		session.Write([]byte(data + "\n"))
	}
}
