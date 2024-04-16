package xtls

import "crypto/tls"

type PKCS12Config struct {
	Path     string
	Password string
}

type ClientOptions struct {
	PKCS12        PKCS12Config `mapstructure:"pkcs12" json:"pkcs12" yaml:"pkcs12"`
	TLSSkipVerify bool         `mapstructure:"tls_skip_verify" json:"tls_skip_verify" yaml:"tls_skip_verify"`
	TLSMinVersion uint16       `mapstructure:"tls_min_version" json:"tls_min_version" yaml:"tls_min_version"`
	TLSMaxVersion uint16       `mapstructure:"tls_max_version" json:"tls_max_version" yaml:"tls_max_version"`
}

func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		TLSSkipVerify: true,
		TLSMinVersion: tls.VersionSSL30,
		TLSMaxVersion: tls.VersionTLS13,
	}
}
