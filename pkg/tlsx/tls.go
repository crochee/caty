package tlsx

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

	"github.com/crochee/lirity"
)

type Config struct {
	Ca   lirity.FileOrContent `json:"ca" yaml:"ca"`
	Cert lirity.FileOrContent `json:"cert" yaml:"cert"`
	Key  lirity.FileOrContent `json:"key" yaml:"key"`
}

// TLSConfig output tls
func TLSConfig(clientAuth tls.ClientAuthType, cfg Config) (*tls.Config, error) {
	caPEMBlock, err := cfg.Ca.Read()
	if err != nil {
		return nil, err
	}
	var certPEMBlock []byte
	if certPEMBlock, err = cfg.Cert.Read(); err != nil {
		return nil, err
	}
	var keyPEMBlock []byte
	if keyPEMBlock, err = cfg.Key.Read(); err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEMBlock) {
		return nil, errors.New("failed to parse root certificate")
	}
	var certificate tls.Certificate
	if certificate, err = tls.X509KeyPair(certPEMBlock, keyPEMBlock); err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		ClientAuth:         clientAuth, // 服务端认证客户端
		ClientCAs:          pool,       // 服务端认证客户端
		CipherSuites:       []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		MinVersion:         tls.VersionTLS12,
		RootCAs:            pool, // 客户端认证服务端
		InsecureSkipVerify: false,
	}, nil
}
