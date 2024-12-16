package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dapplux/email-validation/config"
	"github.com/dapplux/email-validation/infrastructure/emailvalidators"
)

func main() {

	cfg, aErr := config.AutoLoad()
	if aErr != nil {
		log.Fatalln(aErr)

		return
	}
	apiKey := cfg.ZeroBounceApiKey
	zp := emailvalidators.NewZerobounceProvider(apiKey)

	http.HandleFunc("/validate-email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		response, err := zp.ValidateEmail(r.Context(), req.Email, "")
		if err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	fmt.Printf("Server is running on port %d...\n", cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
