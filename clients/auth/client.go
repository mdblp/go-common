package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/mdblp/shoreline/token"
)

// ClientInterface interface that we will implement and mock
type ClientInterface interface {
	Authenticate(req *http.Request) *token.TokenData
}

// Client holds the state of the Auth Client
type Client struct {
	authSecret     string
	tokenValidator *validator.Validator
}

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope    string   `json:"scope"`
	Roles    []string `json:"http://your-loops.com/roles"`
	IsServer bool     `json:"isServer"`
}

// Nothing to validate for tidewhisperer, roles do not matter
// TODO: should we check that the email is verified? info should be in the token
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// WithCustomCa is a Provider Option for our jwks CachingProvider
// It is used to specify a local CA cert, usefull when using a local OAuth server which use a self-signed cert
func WithCustomCA(pem string) jwks.ProviderOption {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM([]byte(pem))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
	}

	return func(p *jwks.Provider) {
		p.Client.Transport = tr
	}
}

func setupAuth0() *validator.Validator {
	//target audience is used to verify the token was issued for a specific domain or url.
	//by default it will be empty but we would (in the future) use this to authorize or deny access to some urls
	targetAudience := []string{}
	if value, present := os.LookupEnv("AUTH0_AUDIENCE"); present {
		targetAudience = []string{value}
	}
	issuerURL, err := url.Parse(os.Getenv("AUTH0_URL") + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}
	keyProvider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	// Use a custom CA cert if it's provided
	if os.Getenv("SSL_CUSTOM_CA_KEY") != "" {
		keyProvider = jwks.NewCachingProvider(issuerURL, 5*time.Minute, WithCustomCA(os.Getenv("SSL_CUSTOM_CA_KEY")))
	}
	jwtValidator, err := validator.New(
		keyProvider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		targetAudience,
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	return jwtValidator
}

// NewClient creates a new Auth Client
func NewClient(authSecret string) (*Client, error) {

	validator := setupAuth0()

	return &Client{
		authSecret:     authSecret,
		tokenValidator: validator,
	}, nil
}

// Authenticate the incomming request using either the x-tidepool-session token or the authorization Bearer token provided by OAuth
func (client *Client) Authenticate(req *http.Request) *token.TokenData {
	if sessionToken := req.Header.Get("x-tidepool-session-token"); sessionToken != "" {
		tokenData, err := token.UnpackSessionTokenAndVerify(sessionToken, client.authSecret)
		//More validations?
		if err != nil {
			log.Print("Error decoding tidepool session token")
		} else {
			return tokenData
		}

	}
	// Defaults to the auth bearer token when the tidepool token is not provided or invalid
	var parsedToken *validator.ValidatedClaims
	if rawToken, err := jwtmiddleware.AuthHeaderTokenExtractor(req); err != nil {
		log.Print("Error decoding bearer token")
		return nil
	} else if t, err := client.tokenValidator.ValidateToken(req.Context(), rawToken); err != nil {
		log.Print("Error decoding bearer token")
		return nil
	} else {
		parsedToken = t.(*validator.ValidatedClaims)
	}
	uid := strings.Split(parsedToken.RegisteredClaims.Subject, "|")[1]
	customClaims := parsedToken.CustomClaims.(*CustomClaims)
	return &token.TokenData{UserId: uid, IsServer: customClaims.IsServer, Role: customClaims.Roles[0]}
}
