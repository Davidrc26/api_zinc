package main

import (
	"encoding/json"
	"fmt"
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
		r.Post("/index", indexHandler)
		r.Get("/job/{jobID}", jobStatusHandler)
	})

	return mux
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := parseBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := models.Response{
			Status:  400,
			Message: "Error al parsear el body de la petición: " + err.Error(),
			Result:  nil,
		}
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	res := services.Search(body)
	w.WriteHeader(res.Status)
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Crear un nuevo job
	jobID := services.CreateJob()

	// Iniciar la indexación de forma asíncrona
	services.StartIndexingAsync(jobID)

	// Retornar inmediatamente con el job ID
	w.WriteHeader(http.StatusAccepted)
	response := models.Response{
		Status:  202,
		Message: "Indexación iniciada. Use el job_id para consultar el estado.",
		Result: map[string]interface{}{
			"job_id":     jobID,
			"status_url": fmt.Sprintf("/api/job/%s", jobID),
		},
	}
	_ = json.NewEncoder(w).Encode(response)
}

func jobStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jobID := chi.URLParam(r, "jobID")
	if jobID == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := models.Response{
			Status:  400,
			Message: "Job ID es requerido",
			Result:  nil,
		}
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	response := services.GetJobResponse(jobID)
	w.WriteHeader(response.Status)
	_ = json.NewEncoder(w).Encode(response)
}
