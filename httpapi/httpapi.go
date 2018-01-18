package httpapi

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func Cmd(postgresConnString string) error {
	r := mux.NewRouter()

	db, err := sql.Open("postgres", postgresConnString)
	if err != nil {
		panic(err)
	}

	webhookEndpoint := WebHookEndpoint{db: db, router: r}

	r.HandleFunc("/", webhookEndpoint.GetList()).Methods("GET")
	r.HandleFunc("/{id}", webhookEndpoint.GetOne()).Methods("GET").Name("GET_ONE_WEBHOOK")
	r.HandleFunc("/", webhookEndpoint.Insert()).Methods("POST")

	return http.ListenAndServe(":8080", r)
}
