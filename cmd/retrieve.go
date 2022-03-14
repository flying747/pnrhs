package main

import (
	"log"
	"net/http"
	"time"

	"github.com/pnrsh/pnrsh/pkg/delta/pnr"
)

func RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Add("Location", "/?error=t")
		w.WriteHeader(302)
		return
	}

	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	confirmationCode := r.Form.Get("confirmation_code")

	if len(confirmationCode) != 6 || len(firstName) == 0 || len(lastName) == 0 {
		w.Header().Add("Location", "/?error=t")
		w.WriteHeader(302)
		return
	}

	retrievedPNR, err := pnr.Retrieve(firstName, lastName, confirmationCode)
	if err != nil {
		w.Header().Add("Location", "/?error=t")
		w.WriteHeader(302)

		go func() {
			time.Sleep(time.Second)
			log.Fatal("shutting down due to PNR error")
		}()

		return
	}

	t := Parse("show.html")

	t.Execute(w, struct {
		PNR              pnr.PNR
		ConfirmationCode string
	}{
		retrievedPNR,
		confirmationCode,
	})
}
