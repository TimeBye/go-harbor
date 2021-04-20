/*
Copyright 2020 The go-harbor Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
*/

package rest

import (
	"context"
	"encoding/base64"
	"fmt"
	flowcontrol2 "github.com/TimeBye/go-harbor/pkg/rest/util/flowcontrol"
	"net"
	"net/http"
	"time"
)

// TLSClientConfig contains settings to enable transport layer security
// +k8s:deepcopy-gen=true
type TLSClientConfig struct {
	// Server should be accessed without verifying the TLS certificate. For testing only.
	Insecure bool
	// ServerName is passed to the server for SNI and is used in the client to check server
	// ceritificates against. If ServerName is empty, the hostname used to contact the
	// server is used.
	ServerName string

	// Server requires TLS client certificate authentication
	CertFile string
	// Server requires TLS client certificate authentication
	KeyFile string
	// Trusted root certificates for server
	CAFile string

	// CertData holds PEM-encoded bytes (typically read from a client certificate file).
	// CertData takes precedence over CertFile
	CertData []byte
	// KeyData holds PEM-encoded bytes (typically read from a client certificate key file).
	// KeyData takes precedence over KeyFile
	KeyData []byte
	// CAData holds PEM-encoded bytes (typically read from a root certificates bundle).
	// CAData takes precedence over CAFile
	CAData []byte

	// NextProtos is a list of supported application level protocols, in order of preference.
	// Used to populate tls.Config.NextProtos.
	// To indicate to the server http/1.1 is preferred over http/2, set to ["http/1.1", "h2"] (though the server is free to ignore that preference).
	// To use only http/1.1, set to ["http/1.1"].
	NextProtos []string
}

const (
	DefaultQPS            float32 = 5.0
	DefaultBurst          int     = 10
	DefaultTimeOut        int     = 10
	DefaultVersionApiPath string  = "/api/v2.0"
)

// Config holds the common attributes that can be passed to a Kubernetes client on
// initialization.
type Config struct {
	// Host must be a host string, a host:port pair, or a URL to the base of the apiserver.
	// If a URL is given then the (optional) Path of that URL represents a prefix that must
	// be appended to all request URIs used to access the apiserver. This allows a frontend
	// proxy to easily relocate all of the apiserver endpoints.
	Host string
	// APIPath is a sub-path that points to an API root.
	APIPath string

	// ContentConfig contains settings that affect how objects are transformed when
	// sent to the server.
	ContentConfig

	// Server requires Basic authentication
	Username string
	Password string

	// Server requires Bearer authentication. This client will not attempt to use
	// refresh tokens for an OAuth2 flow.
	// TODO: demonstrate an OAuth2 compatible client.
	BearerToken string

	// Path to a file containing a BearerToken.
	// If set, the contents are periodically read.
	// The last successfully read value takes precedence over BearerToken.
	BearerTokenFile string

	/*	// Impersonate is the configuration that RESTClient will use for impersonation.
		Impersonate ImpersonationConfig

		// Server requires plugin-specified authentication.
		AuthProvider *clientcmdapi.AuthProviderConfig

		// Callback to persist config for AuthProvider.
		AuthConfigPersister AuthProviderConfigPersister

		// Exec-based authentication provider.
		ExecProvider *clientcmdapi.ExecConfig*/

	// TLSClientConfig contains settings to enable transport layer security
	TLSClientConfig

	// UserAgent is an optional field that specifies the caller of this request.
	UserAgent string

	// DisableCompression bypasses automatic GZip compression requests to the
	// server.
	DisableCompression bool

	// Transport may be used for custom HTTP behavior. This attribute may not
	// be specified with the TLS client certificate options. Use WrapTransport
	// to provide additional per-server middleware behavior.
	Transport http.RoundTripper
	// WrapTransport will be invoked for custom HTTP behavior after the underlying
	// transport is initialized (either the transport created from TLSClientConfig,
	// Transport, or http.DefaultTransport). The config may layer other RoundTrippers
	// on top of the returned RoundTripper.
	//
	// A future release will change this field to an array. Use config.Wrap()
	// instead of setting this value directly.
	//	WrapTransport transport.WrapperFunc

	// QPS indicates the maximum QPS to the master from this client.
	// If it's zero, the created RESTClient will use DefaultQPS: 5
	QPS float32

	// Maximum burst for throttle.
	// If it's zero, the created RESTClient will use DefaultBurst: 10.
	Burst int

	// Rate limiter for limiting connections to the master from this client. If present overwrites QPS/Burst
	RateLimiter flowcontrol2.RateLimiter

	// The maximum length of time to wait before giving up on a server request. A value of zero means no timeout.
	Timeout time.Duration

	// Dial specifies the dial function for creating unencrypted TCP connections.
	Dial func(ctx context.Context, network, address string) (net.Conn, error)

	// Version forces a specific version to be used (if registered)
	// Do we need this?
	// Version string
}

// RESTClientFor returns a RESTClient that satisfies the requested attributes on a client Config
// object. Note that a RESTClient may require fields that are optional when initializing a Client.
// A RESTClient created by this method is generic - it expects to operate on an API that follows
// the Kubernetes conventions, but may not be the Kubernetes API.
func RESTClientFor(config *Config) (*RESTClient, error) {
	qps := config.QPS
	if config.QPS == 0.0 {
		qps = DefaultQPS
	}
	burst := config.Burst
	if config.Burst == 0 {
		burst = DefaultBurst
	}

	if config.Timeout == 0 {
		burst = DefaultTimeOut
	}

	baseURL, err := DefaultServerURL(config.APIPath)
	if err != nil {
		return nil, err
	}
	transport := TransportFor()

	var httpClient *http.Client
	if transport != http.DefaultTransport {
		httpClient = &http.Client{Transport: transport}
		if config.Timeout > 0 {
			httpClient.Timeout = config.Timeout
		}
	}
	headers := map[string]string{}
	if config.Username != "" && config.Password != "" {
		pwd := fmt.Sprintf("%s:%s", config.Username, config.Password)
		headers["authorization"] = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(pwd)))
	}
	return NewRESTClient(baseURL, DefaultVersionApiPath, config.ContentConfig, headers, qps, burst, config.RateLimiter, httpClient)
}

func NewDefaultConfig(host string, username string, password string) *Config {
	return &Config{
		APIPath:  host,
		Username: username,
		Password: password,
	}
}
