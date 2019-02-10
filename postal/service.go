package postal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	base           = "https://service-homolog.digipix.com.br/v0b"
	endpoint       = "/shipments/zipcode/%s"
	authentication = "/auth"
)

//service internal use
type service struct {
	base   string
	token  string
	client *http.Client
}

//Options used to pass options to service
type Options struct {
	URLBase string
	Token   string
}

//Address represents the address of postal code
type Address struct {
	Name         string `json:"name"`
	Zipcode      string `json:"zipcode"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	StateShort   string `json:"state_short"`
	City         string `json:"city"`
}

//ServerError represents the error returned by server
type ServerError struct {
	Error   string `json:"error"`
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

//Consult get information about postal
//if it returns null, it means that the address was not found
func (s service) Consult(postalCode string) (*Address, error) {
	log.Printf("Consulting postal code: %s", postalCode)

	validPostalCode, isValid := ValidatePostalCode(postalCode)
	if !isValid {
		return nil, fmt.Errorf("Invalid postal code: '%s'", validPostalCode)
	}

	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf(s.base+endpoint, validPostalCode), nil)
	if err != nil {
		return nil, fmt.Errorf("error on create request: %v", err)
	}
	if s.token != "" {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	}

	resp, err := s.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error on request: %v", err)
	}
	defer resp.Body.Close()
	log.Printf("Server return status: %d", resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on read body: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		address := &Address{}
		if err := json.Unmarshal(b, &address); err != nil {
			return nil, fmt.Errorf("error on decode json: %v", err)
		}
		log.Println("Success")
		return address, nil

	case http.StatusUnauthorized:
		serverError := &ServerError{}
		if err := json.Unmarshal(b, &serverError); err != nil {
			return nil, fmt.Errorf("error on decode json error: %v", err)
		}
		return nil, fmt.Errorf("Server return unauthorized: %s", serverError.Message)

	case http.StatusNotFound:
		return nil, nil

	default:
		return nil, fmt.Errorf("Server return: %s", string(b))
	}

}

//New create new service to use
func New(opts *Options) (s service) {
	s.client = http.DefaultClient
	s.base = base
	if opts != nil {
		if opts.URLBase != "" {
			log.Printf("Changing default url base to: %s", opts.URLBase)
			s.base = opts.URLBase
		}
		if opts.Token != "" {
			log.Println("Setting token authentication")
			s.token = opts.Token
		}
	}
	return
}
