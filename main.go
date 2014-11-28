package main

import (
	"fmt"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/currency"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	AmountToCharge uint64 = 10000
)

type Home struct {
	PublishableKey string
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("home.html")
	home := Home{
		PublishableKey: os.Getenv("STRIPE_PUBLISHABLE_KEY"),
	}
	t.Execute(w, home)
}

func createDebit(token string, amount uint64, description string) *stripe.Charge {
	stripe.Key = os.Getenv("STRIPE_KEY")

	params := &stripe.ChargeParams{
		Amount:   amount,
		Currency: currency.USD,
		Card: &stripe.CardParams{
			Token: token,
		},
		Desc: description,
	}

	ch, err := charge.New(params)

	if err != nil {
		log.Fatalf("error while trying to charge a cc", err)
	}

	log.Printf("debit created successfully %v\n", ch.ID)

	return ch
}

func debitsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		createDebit(r.FormValue("stripeToken"), AmountToCharge, "testing charge description!")
		fmt.Fprint(w, "successful payment.")
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/debits", debitsHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8081", nil)
}
