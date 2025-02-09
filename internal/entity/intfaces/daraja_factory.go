package intfaces

import (
	"github.com/harmannkibue/golang-mpesa-sdk/pkg/daraja"
	"strings"
)

// DarajaFactory is responsible for creating Daraja instances per tenant
type DarajaFactory interface {
	GetDarajaInstance() (*daraja.DarajaService, error)
}

func NewDarajaFactory() DarajaFactory {
	return &TenantMpesaConfig{}
}

// TenantMpesaConfig holds M-PESA credentials for a tenant
type TenantMpesaConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	Shortcode      string
	PassKey        string
	Environment    string
}

// DarajaEnvironment represents the possible environments for Safaricom Daraja API.
type DarajaEnvironment int

const (
	// Sandbox Environment (Test mode) -> 1
	Sandbox DarajaEnvironment = 1

	// Production Environment (Live API mode) -> 2
	Production DarajaEnvironment = 2
)

// String returns the string representation of the DarajaEnvironment enum.
func (e DarajaEnvironment) String() string {
	switch e {
	case Sandbox:
		return "sandbox"
	case Production:
		return "production"
	default:
		return "unknown"
	}
}

// ParseDarajaEnvironment converts a string to the correct Daraja environment integer (1 or 2).
func ParseDarajaEnvironment(env string) int {
	lowerEnv := strings.ToLower(env)

	if lowerEnv == "production" {
		return int(Production) // 2
	}
	return int(Sandbox) // 1 (default)
}

// GetDarajaInstance returns the correct Daraja SDK instance for a given tenant
func (c TenantMpesaConfig) GetDarajaInstance() (*daraja.DarajaService, error) {

	// TODO: This is just dummy sandbox for testing but essentially the system should be multi tenant in the way it loads and initializes daraja instances -.
	darajaService, err := daraja.New("YCWLzwRKndyK7Il3NbyshtxJHZ9aYjTGEkjwdj0NpET46kAo", "Ar7cNir15oZdA9F5YVr8AMiF57Lb7errEmtR10OdRwrAXudTDAWDB0Aa12hGI5bR", "bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919", daraja.SANDBOX)

	if err != nil {
		return nil, err
	}

	return darajaService, nil
}
