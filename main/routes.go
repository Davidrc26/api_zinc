package main

import (
	"encoding/json"
	"net/http"

	"github.com/Davidrc26/api_zinc.git/models"
	"github.com/Davidrc26/api_zinc.git/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

func Routes() *chi.Mux {
	mux := chi.NewMux()

	mux.Use(
		middleware.Logger,
		middleware.Recoverer,
	)
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	mux.Route("/api", func(r chi.Router) {
		r.Use(cors.Handler)
		r.Post("/search", searchHandler)
	})

	return mux
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := parseBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m := map[string]interface{}{"msg": "Error realizando la busqueda"}
		_ = json.NewEncoder(w).Encode(m)
	}
	res, err := services.Search(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]interface{}{"msg": "Error en el servidor. Detalles: \n" + err.Error()}
		_ = json.NewEncoder(w).Encode(m)
	}
	_ = json.NewEncoder(w).Encode(res)
}

func parseBody(r *http.Request) (*models.SearchBody, error) {
	body := r.Body
	defer body.Close()
	var searchParams models.SearchBody

	err := json.NewDecoder(body).Decode(&searchParams)
	if err != nil {
		return nil, err
	}

	return &searchParams, nil
}
