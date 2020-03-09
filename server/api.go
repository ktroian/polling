package main

import (
	"log"
	"fmt"
	"context"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/gorilla/mux"
)

type Handler struct {
	updates chan bool
	companies *AllCompanies
}

func (h *Handler) createCompany(w http.ResponseWriter, r *http.Request) {
	var company Company

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Not enough data")
		log.Println(err)
	}
	
	json.Unmarshal(reqBody, &company)
	h.companies.add(company)
	h.updates <- true
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(company)
}


func (h *Handler) getAllCompanies(w http.ResponseWriter, r *http.Request) {
	companies := h.companies.get()
	json.NewEncoder(w).Encode(companies)
}


func (h *Handler) deleteCompany(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.companies.delete(id)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Fatal(err)
	} else {
		h.updates <- true
		fmt.Fprintf(w, "Success!")
	}
}

func serveRoutes(ctx context.Context, updates chan bool, companies *AllCompanies, port int) {
	go func() {
		// routing goroutine
		h := Handler{ updates, companies }
		router := mux.NewRouter().StrictSlash(true)
		defer log.Println("Server Stopped")

		router.HandleFunc("/company", h.createCompany).Methods("POST")
		router.HandleFunc("/companies", h.getAllCompanies).Methods("GET")
		router.HandleFunc("/companies/{id}", h.deleteCompany).Methods("DELETE")
		log.Println("Starting server on the port", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
	}()
}
