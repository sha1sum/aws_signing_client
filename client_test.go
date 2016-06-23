package aws_signing_client

import (
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

var (
	creds   = credentials.NewStaticCredentials("ID", "SECRET", "TOKEN")
	signer  *v4.Signer
	client  *http.Client
	service string
	region  string

	newClient *http.Client
)

func init() {
	Init()
}

func Init() {
	signer = v4.NewSigner(creds)
	client = http.DefaultClient
	service = "es"
	region = "us-east-1"
	newClient, _ = nc()
}

func nc() (*http.Client, error) {
	return NewClient(signer, client, service, region)
}

// TestNewClientWithoutSigner tests the NewClient() function when it is not passed a *v4.Signer.
func TestNewClientWithoutSigner(t *testing.T) {
	Init()
	signer = nil
	_, err := nc()
	if err != (MissingSignerError{}) {
		t.Error("Error was not of type MissingSignerError")
	}
}

// TestNewClientWithoutService tests the NewClient() function when it is not passed a service string.
func TestNewClientWithoutService(t *testing.T) {
	Init()
	service = ""
	_, err := nc()
	if err != (MissingServiceError{}) {
		t.Error("Error was not of type MissingServiceError")
	}
}

// TestNewClientWithoutRegion tests the NewClient() function when it is not passed a region string.
func TestNewClientWithoutRegion(t *testing.T) {
	Init()
	region = ""
	_, err := nc()
	if err != (MissingRegionError{}) {
		t.Error("Error was not of type MissingRegionError")
	}
}

// TestNewClient tests the NewClient() function when all is right in the World.
func TestNewClient(t *testing.T) {
	Init()
	newClient, err := nc()
	switch {
	case err != nil:
		t.Errorf("An unexpected error occurred while creating a new client with valid parameters: %s", err)
	case newClient == nil:
		t.Error("A nil *http.Client was returned while creating a new client with valid parameters")
	}
}

