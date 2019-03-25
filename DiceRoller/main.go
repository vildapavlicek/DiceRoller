package main

import (
	"log"
	"os"
	"time"

	httpclient "github.com/vildapavlicek/GoLang/DiceGetter/httpClient"
	"github.com/vildapavlicek/GoLang/DiceGetter/models"
)

func main() {
	timeout := 5 * time.Second
	client := httpclient.New(timeout, 10, "https://enep8mt0cn7tq.x.pipedream.net")

	diceRolls := models.New(client)

	response, err := diceRolls.Client.GetResponse("GET", nil)
	if err != nil {
		log.Printf("failed to get correct response, error: %s ", err)
		os.Exit(1)
	}

	err = diceRolls.ParseHTML(response)
	if err != nil {
		log.Printf("Failed to parse HTML, reason: %s", err)
	}
	diceRolls.OrderResults(os.Stderr)
	diceRolls.BucketResults(os.Stdout)
	err = diceRolls.DoPost()
	if err != nil {
		log.Printf("Failed to POST data, reason: %s", err)
	}

}
