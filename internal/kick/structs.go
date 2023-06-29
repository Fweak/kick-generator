package kick

import (
	"log"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/gorilla/websocket"
)

// Form is to define the response from ""
type Form struct {
	Enabled                   bool   `json:"enabled"`
	NameFieldName             string `json:"nameFieldName"`
	UnrandomizedNameFieldName string `json:"unrandomizedNameFieldName"`
	ValidFromFieldName        string `json:"validFromFieldName"`
	EncryptedValidFrom        string `json:"encryptedValidFrom"`
}

// Client is to define the over-all package of kick
type Client struct {
	Conn     *websocket.Conn
	Email    string
	Password string
	socketID string
	xsrf     string

	request tls_client.HttpClient
	form    Form
}

// CreateClient creates and adds our http_client
func CreateClient(email, password string) *Client {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(tls_client.Firefox_108),
		tls_client.WithRandomTLSExtensionOrder(),
		tls_client.WithCookieJar(tls_client.NewCookieJar(tls_client.WithLogger(tls_client.NewLogger()))),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		Email:    email,
		Password: password,
		request:  client,
	}
}
