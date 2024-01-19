package main 
import (
	"encoding/json"
	"fmt"
	"net/http"
)

// type RequestPayload struct {
// 	ID 		string `json:"id"`
// 	PubKey string `json:"pubkey"`
// 	CreatedAt int64          `json:"created_at"`
// 	Kind      int            `json:"kind"`
// 	Tags      [][]interface{} `json:"tags"`
// 	Content   string         `json:"content"`
// 	Sig       string         `json:"sig"`
// }

type RequestPayload struct {
	ID        string         `json:"id"`
	PubKey    string         `json:"pubkey"`
	CreatedAt int64          `json:"created_at"`
	Kind      int            `json:"kind"`
	Ta        string         `json:"ta"`
	Amt       float64         `json:"amt"`
	Content   string         `json:"content"`
	Sig       string         `json:"sig"`
}


func validateNoStrEventMiddleware (next func(http.ResponseWriter, *http.Request, RequestPayload) RequestPayload) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload RequestPayload
	
		decoder := json.NewDecoder(r.Body)


		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		// create TAHUB_RECEIVE_ADDRESS_FOR_ASSET event 
		if payload.Kind == 1 && payload.Content == "TAHUB_RECEIVE_ADDRESS_FOR_ASSET" {
			
	
				if len(payload.Ta) == 0 {
					http.Error(w, "Field 'ta' must exist and not be empty", http.StatusBadRequest)
					return
				}

				if payload.Ta == "" {
					http.Error(w, "Field 'ta' must be a non-empty string", http.StatusBadRequest)
					return
				}


				if payload.Amt < 0 || payload.Amt != float64(int64(payload.Amt)) {
					http.Error(w, "Field 'amt' must be a positive integer (u64)", http.StatusBadRequest)
					return
				}
	
	
				// If conditions are met, proceed to the next handler
				next(w, r, payload)
				return
		}

		// check for TAHUB_CREATE_USER
		if payload.Kind == 1 && payload.Content == "TAHUB_CREATE_USER" {
			next(w,r, payload)
			return 
		}


		// check for TAHUB_GET_BALANCES
		if payload.Kind == 1 && payload.Content == "TAHUB_GET_BALANCES" {
			next(w,r, payload)
			return 
		}

			// Otherwise, return an error response
			http.Error(w, "Invalid event content", http.StatusBadRequest)

	}
}

func processNoStrEvent (w http.ResponseWriter, r *http.Request, payload RequestPayload) RequestPayload {
	fmt.Fprint(w, "Valid event content")
	fmt.Println(payload)
	return payload
}

func main() {
	fmt.Println("Validate Nostr Event middleware")
	http.HandleFunc("/dispatchevent", func(w http.ResponseWriter, r *http.Request) {
		validateNoStrEventMiddleware(processNoStrEvent)(w, r)
	})
	// Start the server
	http.ListenAndServe(":8180", nil)

}