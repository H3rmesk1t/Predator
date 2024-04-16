package xtls

import (
	"crypto/tls"
	"io/ioutil"
	"software.sslmate.com/src/go-pkcs12"
)

func parsePKCS12FromFile(c PKCS12Config) (*tls.Certificate, error) {
	data, err := ioutil.ReadFile(c.Path)
	if err != nil {
		return nil, err
	}

	privateKey, certificate, _, err := pkcs12.DecodeChain(data, c.Password)
	if err != nil {
		return nil, err
	}
	return &tls.Certificate{
		Certificate: [][]byte{certificate.Raw},
		PrivateKey:  privateKey,
		Leaf:        certificate,
	}, nil
}

func NewTLSConfig(options *ClientOptions) (*tls.Config, error) {
	var err error
	var cert *tls.Certificate

	if options.PKCS12.Path != "" && options.PKCS12.Password != "" {
		cert, err = parsePKCS12FromFile(options.PKCS12)
		if err != nil {
			return nil, err
		}
	}

	tlsClientConfig := &tls.Config{
		InsecureSkipVerify: options.TLSSkipVerify,
		MinVersion:         options.TLSMinVersion,
		MaxVersion:         options.TLSMaxVersion,
	}

	if cert != nil {
		tlsClientConfig.Certificates = []tls.Certificate{*cert}
	}
	return tlsClientConfig, nil
}
