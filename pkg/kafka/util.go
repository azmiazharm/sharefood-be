// Package kafka
package kafka

import (
	"crypto/sha256"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"hash"
	"io/ioutil"

	"sharefood/pkg/logger"

	"github.com/xdg/scram"
)

var SHA256 scram.HashGeneratorFcn = func() hash.Hash { return sha256.New() }
var SHA512 scram.HashGeneratorFcn = func() hash.Hash { return sha512.New() }

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

func createTlsConfig(c TLS) (t *tls.Config) {
	t = &tls.Config{
		InsecureSkipVerify: c.SkipVerify,
	}
	if c.CertFile != "" && c.KeyFile != "" && c.CaFile != "" {
		cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
		if err != nil {
			logger.Fatal(err)
		}

		caCert, err := ioutil.ReadFile(c.CaFile)
		if err != nil {
			logger.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: c.SkipVerify,
		}
	}
	return t
}
