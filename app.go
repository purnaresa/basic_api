package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/rs/cors"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
)

var fruits []Fruits

type Fruits struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Price   string `json:"price"`
	Picture string `json:"picture"`
}

func init() {
	banana := Fruits{
		ID:      1,
		Name:    "Banana",
		Price:   "Rp 5.000",
		Picture: "http://saltmarshrunning.com/wp-content/uploads/2014/09/bananasf.jpg",
	}
	apple := Fruits{
		ID:      2,
		Name:    "Apple",
		Price:   "Rp 3.000",
		Picture: "http://weknowyourdreams.com/images/apple/apple-06.jpg",
	}

	fruits = append(fruits, banana)
	fruits = append(fruits, apple)
}

func main() {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH"},
	})

	mux := goji.NewMux()
	mux.Use(corsHandler.Handler)

	mux.HandleFuncC(pat.Get("/fruit"), Read)
	mux.HandleFuncC(pat.Get("/fruit/:id"), ReadDetail)
	mux.HandleFuncC(pat.Post("/fruit"), Create)
	mux.HandleFuncC(pat.Post("/fruit/:id/delete"), Delete)
	log.Printf("Run on localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func Read(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fruits)
}

func ReadDetail(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idString := pat.Param(ctx, "id")
	id, errParse := strconv.ParseInt(idString, 10, 64)
	if errParse != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, fruit := range fruits {
		if fruit.ID == id {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(fruit)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.FormValue("id")
	name := r.FormValue("name")
	price := r.FormValue("price")
	picture := r.FormValue("picture")

	idInt, errParse := strconv.ParseInt(id, 10, 64)
	if errParse != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fruit := Fruits{
		ID:      idInt,
		Name:    name,
		Price:   price,
		Picture: picture,
	}

	fruits = append(fruits, fruit)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fruit)
}

func Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idString := pat.Param(ctx, "id")
	id, errParse := strconv.ParseInt(idString, 10, 64)
	if errParse != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for index, fruit := range fruits {
		if fruit.ID == id {
			fruits = append(fruits[:index], fruits[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
