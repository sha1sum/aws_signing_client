package aws_signing_client

import (
	"net/http"

	"bytes"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/private/protocol/rest"
)

type (
	// Signer implements the http.RoundTripper interface and houses an optional RoundTripper that will be called between
	// signing and response.
	Signer struct {
		transport http.RoundTripper
		v4        *v4.Signer
		service   string
		region    string
	}

	// MissingSignerError is an implementation of the error interface that indicates that no AWS v4.Signer was
	// provided in order to create a client.
	MissingSignerError struct{}

	// MissingServiceError is an implementation of the error interface that indicates that no AWS service was
	// provided in order to create a client.
	MissingServiceError struct{}

	// MissingRegionError is an implementation of the error interface that indicates that no AWS region was
	// provided in order to create a client.
	MissingRegionError struct{}
)

// NewClient obtains an HTTP client with a RoundTripper that signs AWS requests for the provided service. An
// existing client can be specified for the `client` value, or--if nil--a new HTTP client will be created.
func NewClient(v4 *v4.Signer, client *http.Client, service string, region string) (*http.Client, error) {
	c := client
	switch {
	case v4 == nil:
		return nil, MissingSignerError{}
	case service == "":
		return nil, MissingServiceError{}
	case region == "":
		return nil, MissingServiceError{}
	case c == nil:
		c = http.DefaultClient
	}
	s := &Signer{
		transport: client.Transport,
		v4:        v4,
		service:   service,
		region:    region,
	}
	if s.transport == nil {
		s.transport = http.DefaultTransport
	}
	c.Transport = s
	return c, nil
}

// RoundTrip implements the http.RoundTripper interface and is used to wrap HTTP requests in order to sign them for AWS
// API calls. The scheme for all requests will be changed to HTTPS. NOTE: This method reads the entire request body in
// order to implement the io.ReadSeeker interface and may be troublesome for very large request bodies.
func (s *Signer) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "https"
	if strings.Contains(req.URL.RawPath, "%2C") {
		req.URL.RawPath = rest.EscapePath(req.URL.RawPath, false)
	}
	t := time.Now()
	req.Header.Set("Date", t.Format(time.RFC3339))
	defer req.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	r := bytes.NewReader(buf.Bytes())
	head, err := s.v4.Sign(req, r, s.service, s.region, t)
	if err != nil {
		return nil, err
	}
	req.Body = r
	req.Header = head
	return s.transport.RoundTrip(req)
}

// Error implements the error interface.
func (err *MissingSignerError) Error() string {
	return "No signer was provided. Cannot create client. Try using the elastic_aws.NewSigner() function."
}

// Error implements the error interface.
func (err *MissingServiceError) Error() string {
	return "No AWS service abbreviation was provided. Cannot create client."
}

// Error implements the error interface.
func (err *MissingRegionError) Error() string {
	return "No AWS region was provided. Cannot create client."
}
