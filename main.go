package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	exitOK int = iota
	exitError

	ENDPOINT = ""
)

func main() {
	os.Exit(realMain(os.Args))
}

func realMain(args []string) int {
	// os.Args[1] is the main.go file path. hence provided bearer token is treated as the second argument.
	if len(os.Args) < 2 {
		fmt.Printf("Usage: bearer token is needed\n")
		return exitError
	}

	rate := vegeta.Rate{Freq: 1, Per: time.Second} // This is for setting the rps.
	duration := time.Second * 1                    // This is the time for conitnueing the test with the rate.

	// This is the body of the request as JSON.
	bodyJson := ``

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		Header: http.Header{
			"Authorization": {
				fmt.Sprintf("Bearer %s", os.Args[1]),
			},
			"Content-Type": {"application/json"},
		},
		URL:  ENDPOINT,
		Body: []byte(bodyJson),
	})

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, fmt.Sprintf("Attack %s", ENDPOINT)) {
		metrics.Add(res)
	}
	metrics.Close()
	log.Printf("%#v\n\n", metrics.Errors)
	vegeta.NewTextReporter(&metrics).Report(os.Stdout)

	return exitOK
}
