package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func generateMockSSHKeys() {
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatalf("Failed to generate private key: %v", err)
		}

		privateKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})

		err = os.WriteFile(privateKeyPath, privateKeyPEM, 0600)
		if err != nil {
			log.Fatalf("Failed to save private key to file: %v", err)
		}

		log.Printf("Generated mock SSH private key at %s\n", privateKeyPath)
	}
}
