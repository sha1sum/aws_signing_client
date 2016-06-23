package aws_signing_client

import (
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"strings"
)

type testRoundTripper struct {
	err error
}

var (
	creds   = credentials.NewStaticCredentials("ID", "SECRET", "TOKEN")
	signer  *v4.Signer
	client  *http.Client
	service string
	region  string

	rt        *testRoundTripper
	newClient *http.Client
	err       error
)

func init() {
	Init()
}

func Init() {
	rt = &testRoundTripper{}
	signer = v4.NewSigner(creds)
	client = http.DefaultClient
	client.Transport = rt
	service = "es"
	region = "us-east-1"
	newClient, _ = nc()
	err = nil
}

func nc() (*http.Client, error) {
	return New(signer, client, service, region)
}

//   _   _                ____ _ _            _
//  | \ | | _____      __/ ___| (_) ___ _ __ | |_
//  |  \| |/ _ \ \ /\ / / |   | | |/ _ \ '_ \| __|
//  | |\  |  __/\ V  V /| |___| | |  __/ | | | |_
//  |_| \_|\___| \_/\_/  \____|_|_|\___|_| |_|\__|
//

// TestNewClientWithoutSigner tests the NewClient() function when it is not passed a *v4.Signer.
func TestNewClientWithoutSigner(t *testing.T) {
	Init()
	signer = nil
	_, err = nc()
	if err != (MissingSignerError{}) {
		t.Error("Error was not of type MissingSignerError")
	}
}

// TestNewClientWithoutService tests the NewClient() function when it is not passed a service string.
func TestNewClientWithoutService(t *testing.T) {
	Init()
	service = ""
	_, err = nc()
	if err != (MissingServiceError{}) {
		t.Error("Error was not of type MissingServiceError")
	}
}

// TestNewClientWithoutRegion tests the NewClient() function when it is not passed a region string.
func TestNewClientWithoutRegion(t *testing.T) {
	Init()
	region = ""
	_, err = nc()
	if err != (MissingRegionError{}) {
		t.Error("Error was not of type MissingRegionError")
	}
}

// TestNewClient tests the NewClient() function when all is right in the World.
func TestNewClient(t *testing.T) {
	Init()
	newClient, err = nc()
	switch {
	case err != nil:
		t.Errorf("An unexpected error occurred while creating a new client with valid parameters: %s", err)
	case newClient == nil:
		t.Error("A nil *http.Client was returned while creating a new client with valid parameters")
	}
}

//   ____                       _ _____     _
//  |  _ \ ___  _   _ _ __   __| |_   _| __(_)_ __
//  | |_) / _ \| | | | '_ \ / _` | | || '__| | '_ \
//  |  _ < (_) | |_| | | | | (_| | | || |  | | |_) |
//  |_| \_\___/ \__,_|_| |_|\__,_| |_||_|  |_| .__/
//                                           |_|

var passedReq *http.Request

// TestRoundTripSignsGetRequest ensures that a GET request is signed before sending.
func TestRoundTripSignsGetRequest(t *testing.T) {
	Init()
	_, err = newClient.Get("https://google.com")
	checkSignatures(t)
}

// TestRoundTripSignsPostRequest ensures that a GET request is signed before sending.
func TestRoundTripSignsPostRequest(t *testing.T) {
	Init()
	_, err = newClient.Post("https://google.com", "application/json", strings.NewReader("{}"))
	checkSignatures(t)
}

func checkSignatures(t *testing.T) {
	switch {
	case err != nil:
		t.Errorf("An unexpected error occurred while making a request: %s", err)
	case passedReq.Header == nil:
		t.Error("nil headers were returned from the signing request")
	case len(passedReq.Header["x-amz-date"]) == 0:
		t.Error("No 'x-amz-date' header was returned from the signing request")
	case len(passedReq.Header["x-amz-security-token"]) == 0:
		t.Error("No 'x-amz-security-token' header was returned from the signing request")
	}
}

func (rt *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	passedReq = req
	return &http.Response{}, rt.err
}
