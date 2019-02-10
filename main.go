package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/filipenos/dgpx/postal"
)

var (
	port  = 8080
	token = ""
)

func init() {
	flag.IntVar(&port, "port", port, "port to listen")
	flag.StringVar(&token, "token", "", "token used in authorization")
	flag.Parse()

	if v := os.Getenv("TOKEN"); v != "" {
		log.Println("Using token from env var")
		token = v
	}
}

func main() {
	http.HandleFunc("/", ZipCodeHandle)

	log.Printf("Running on port :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

//ZipCodeHandle handler used to show postal code
func ZipCodeHandle(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = make(map[string]interface{}, 0)
		temp = "templates/consult.html"
		cep  = strings.TrimSpace(r.FormValue("cep"))
	)

	switch r.Method {
	case http.MethodGet:
		log.Println("Consult postal code")

	case http.MethodPost:
		log.Println("Show results")

		if cep == "" {
			ctx["Msg"] = "CEP não pode estar vazio"
		} else {
			s := postal.New(&postal.Options{Token: token})

			if address, err := s.Consult(cep); err != nil {
				ctx["Msg"] = "Ocorreu um erro ao consultar o CEP: " + err.Error()
			} else {
				ctx["Address"] = address
				temp = "templates/result.html"
			}
		}

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	t := template.Must(template.ParseFiles("templates/base.html", temp))
	if err := t.Execute(w, ctx); err != nil {
		log.Fatalf("Unexpected error on execute template: %v", err)
	}
}
