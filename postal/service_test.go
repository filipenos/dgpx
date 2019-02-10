package postal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	AuthToken = "ABC123"
)

func TestConsultAddress(t *testing.T) {

	expectAddress := Address{
		Street:       "My Street",
		Neighborhood: "My Neighborhood",
		City:         "My City",
		StateShort:   "SP",
	}

	ts := newTestServer(t, expectAddress)
	defer ts.Close()

	s := New(&Options{URLBase: ts.URL, Token: AuthToken})

	address, err := s.Consult("12345678")
	if err != nil {
		t.Fatalf("Unexpected error on consult: %v", err)
	}
	if address == nil {
		t.Fatalf("Unexpected <nil> address")
	}
	if expectAddress != *address {
		t.Fatalf("Unexpected address %#v, expect %#v", address, expectAddress)
	}
}

func TestAddressNotFound(t *testing.T) {
	ts := newTestServer(t, nil)
	defer ts.Close()

	s := New(&Options{URLBase: ts.URL, Token: AuthToken})

	address, err := s.Consult("00000000")
	if err != nil {
		t.Fatalf("Unexpected error on consult address: %v", err)
	}
	if address != nil {
		t.Fatalf("Unexpected address %#v, expect <nil>", address)
	}
}

func TestAuthenticationFail(t *testing.T) {
	ts := newTestServer(t, nil)
	defer ts.Close()

	s := New(&Options{URLBase: ts.URL, Token: "Token invalido"})

	address, err := s.Consult("12345678")
	if err == nil {
		t.Fatalf("Unexpected <nil> error, expect unauthorized")
	}
	if address != nil {
		t.Fatalf("Unexpected address %#v, expect <nil>", address)
	}
	if !strings.Contains(err.Error(), "Acesso não autorizado.") {
		t.Fatalf("Unexpected error: '%v', expect 'Acesso não autorizado.'", err)
	}
}

func newTestServer(t *testing.T, e interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", AuthToken) {
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(ServerError{
				Kind:    "ACCESS_UNAUTHORIZED",
				Message: "Acesso não autorizado.",
			}); err != nil {
				t.Fatalf("Unexpected error on create error message: %v", err)
			}
			return
		}
		if e == nil {
			http.Error(w, "Zipcode not found", http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(e); err != nil {
			t.Fatalf("Unexpected error on create test server: %v", err)
		}
	}))
}
