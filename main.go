package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	serieNumber, lottoNumber int
	fileSource               string
)

func main() {
	// Parse number from cmd arguments
	flag.IntVar(&serieNumber, "s", serieNumber, "Serie number")
	flag.IntVar(&lottoNumber, "l", lottoNumber, "Lotto number")
	flag.StringVar(&fileSource, "source", fileSource, "File Source")
	flag.Parse()

	if fileSource != "" {
		file, err := os.Open(fileSource)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			result := strings.Split(text, " ")

			serieNumber, _ = strconv.Atoi(result[0])
			lottoNumber, _ = strconv.Atoi(result[1])
			checkLot(serieNumber, lottoNumber)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	} else {
		checkLot(serieNumber, lottoNumber)
	}
}

func checkLot(serieNumber, lottoNumber int) {
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
	fmt.Println()
}
