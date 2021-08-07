package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gopsql/goconf"
	"golang.org/x/crypto/ssh"
)

type (
	SSHPrivateKey []byte

	Configs struct {
		Address string "Address to listen to"

		Remote string "Remote url of the git repository"
		Local  string "Local path of the git repository"

		SSHPrivateKey SSHPrivateKey "Base64 encoded SSH private key"

		SSHPrivateKeyPassword string "Password of the SSH private key"

		UserName  string "Author name of a new git commit"
		UserEmail string "Author email of a new git commit"
	}
)

func (s *SSHPrivateKey) SetString(input string) (err error) {
	*s = []byte(strings.TrimSpace(input))
	return nil
}

func (s SSHPrivateKey) String() string {
	return "\n" + string(s)
}

func hasConfigs() bool {
	_, err := os.Stat(configFileLocation)
	return err == nil
}

func getConfigs() *Configs {
	c := Configs{}
	content, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		return &c
	}
	goconf.Unmarshal(content, &c)
	return &c
}

func updateConfigs(configs *Configs) {
	if len(configs.Address) == 0 {
		configs.Address = "127.0.0.1:8080"
	}
	if len(configs.SSHPrivateKey) == 0 {
		privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		keyPem := pem.EncodeToMemory(&pem.Block{
			Type:    "RSA PRIVATE KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
		})
		configs.SSHPrivateKey = keyPem
		pub, _ := ssh.NewPublicKey(&privateKey.PublicKey)
		log.Println("new public key:", strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pub))))
	}
	content, _ := goconf.Marshal(*configs)
	ioutil.WriteFile(configFileLocation, content, 0600)
}
