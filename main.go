package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var (
	serieNumber, lottoNumber int
)

func main() {

	// Parse number from cmd arguments
	flag.IntVar(&serieNumber, "s", serieNumber, "Serie number")
	flag.IntVar(&lottoNumber, "l", lottoNumber, "Lotto number")
	flag.Parse()

	// Crate request
	request := fmt.Sprintf("https://www.bingolotto.se/ratta-lotten/?S=%d&L=%d", serieNumber, lottoNumber)

	// Send a get request
	respond, _ := http.Get(request)

	// Read the response body
	respondBody, _ := io.ReadAll(respond.Body)

	// Create the regex
	re := regexp.MustCompile(`<p.*class=\"lead\">(.*)</p>`)

	// Printout the output when match the regex from the response
	submatchall := re.FindAllStringSubmatch(string(respondBody), -1)
	for _, element := range submatchall {
		fmt.Println(element[1])
	}
}
