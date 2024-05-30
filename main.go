package main

import (
	"log"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	sshPort := ":2222"
	server := &ssh.Server{
		Addr:    sshPort,
		Handler: handleSSHSession,
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

	log.Fatal(server.ListenAndServe())
}

func handleSSHSession(session ssh.Session) {}
