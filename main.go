package main

import (
	"net/http"
	"registro_llamadas/filtros"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()
	cosrsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"*"},
	})
	handlerCros := cosrsOptions.Handler(r)
	r.HandleFunc("/llamadas", filtros.GetLlamadas).Methods("GET")
	r.HandleFunc("/llamadas", filtros.GetLlamadasPorDisposition).Methods("POST")
	r.HandleFunc("/llamadas/fecha", filtros.GetLlamadasPorRangoFecha).Methods("POST")
	r.HandleFunc("/llamadas/src", filtros.GetLlamadasPorSrc).Methods("GET")
	http.ListenAndServe(":8080", handlerCros)
}
