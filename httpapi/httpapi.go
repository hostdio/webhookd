package httpapi

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type HeaderOverride []Header

func (h HeaderOverride) JSON() []byte {
	byt, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return byt
}

type Webhook struct {
	ID             string          `json:"id"`
	URL            string          `json:"url"`
	HeaderOverride *HeaderOverride `json:"headerOverride,omitempty"`
}

type ResponsePayload struct {
	NextPageToken string `json:"nextPageToken"`
}

type WebhooksResponsePayload struct {
	*ResponsePayload
	Webhooks []Webhook `json:"webhooks"`
}

type InsertResponsePayload struct {
	Ref string `json:"_ref"`
}

func Cmd(postgresConnString string) error {
	r := mux.NewRouter()

	db, err := sql.Open("postgres", postgresConnString)
	if err != nil {
		panic(err)
	}

	r.HandleFunc("/", GetList(db)).Methods("GET")
	// r.HandleFunc("/{id}", GetOne(db)).Methods("GET").Name("GET_ONE_WEBHOOK")
	r.HandleFunc("/", Insert(db, r)).Methods("POST")

	return http.ListenAndServe(":8080", r)
}

func GetList(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT * FROM webhooks")
		if err != nil {
			panic(err)
		}

		webhooks := []Webhook{}
		for rows.Next() {
			w := Webhook{}
			var headerOverrideByt []byte
			if scanErr := rows.Scan(&w.ID, &w.URL, &headerOverrideByt); scanErr != nil {
				panic(scanErr)
			}
			if headerOverrideByt != nil {
				if err := json.Unmarshal(headerOverrideByt, &w.HeaderOverride); err != nil {
					panic(err)
				}
			}
			webhooks = append(webhooks, w)
		}
		payload := WebhooksResponsePayload{Webhooks: webhooks}
		byt, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}
		w.Write(byt)

		w.Header().Set("content-type", "application/json")
	}
}

func Insert(db *sql.DB, router *mux.Router) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		stmt, prepErr := db.Prepare("INSERT INTO webhooks (id, url, header_override) VALUES ($1, $2, $3)")
		if prepErr != nil {
			panic(prepErr)
		}
		defer stmt.Close()
		id := uuid.NewV4().String()

		byt, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			panic(readErr)
		}
		defer r.Body.Close()

		var webhook Webhook
		if unmarsalErr := json.Unmarshal(byt, &webhook); unmarsalErr != nil {
			panic(unmarsalErr)
		}

		_, insertErr := stmt.Exec(id, webhook.URL, webhook.HeaderOverride.JSON())
		if insertErr != nil {
			panic(insertErr)
		}

		url, urlErr := router.Get("GET_ONE_WEBHOOK").URL("id", id)
		if urlErr != nil {
			panic(urlErr)
		}
		payload := InsertResponsePayload{
			Ref: url.String(),
		}

		payloadByt, marshallErr := json.Marshal(payload)
		if marshallErr != nil {
			panic(marshallErr)
		}

		w.Write(payloadByt)

	}
}

func GetOne(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}
}
