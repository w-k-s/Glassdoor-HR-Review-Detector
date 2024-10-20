package server

import (
	"encoding/json"
	"log"
	"net/http"

	pd "github.com/josestg/problemdetail"
)

type RESTHandler = func(r *http.Request) (body any, err error, status int)

func RESTEndpoint(api RESTHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err, status := api(r)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		if err == nil {
			if body == nil {
				return
			}
			if err := json.NewEncoder(w).Encode(body); err != nil {
				log.Printf("Failed to encode json. Error: %s", err)
				http.Error(w, "Error encoding response", http.StatusInternalServerError)
			}
			return
		}

		log.Printf("%d %s %s %s", status, r.Method, r.RequestURI, err)

		problemDetail := pd.New(
			pd.Untyped,
			pd.WithValidateLevel(pd.LStandard),
			pd.WithInstance(r.RequestURI),
			pd.WithDetail(err.Error()),
		)
		pd.WriteJSON(w, problemDetail, status)
	})
}
