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
	"sync"
)

var (
	serieNumber, lottoNumber int
	fileSource               string
	wg                       sync.WaitGroup
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
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			result := strings.Split(text, " ")

			serieNumber, _ = strconv.Atoi(result[0])
			lottoNumber, _ = strconv.Atoi(result[1])
			wg.Add(1)
			go checkLot(serieNumber, lottoNumber)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	} else {
		wg.Add(1)
		go checkLot(serieNumber, lottoNumber)
	}

	wg.Wait()
}

func checkLot(serieNumber, lottoNumber int) {
	defer wg.Done()
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

	var winString string
	re_success := regexp.MustCompile(`<div class="alert alert-success">(.|\n)*?<\/div>`)

	// Check if the returned value contains "tyvärr", in this case the output will be in red, otherwise green if you won :)
	var terminalColor string
	if strings.Contains(submatchall[1][1], "grattis") {
		terminalColor = "\033[32m"
		submatchall1 := re_success.FindAllStringSubmatch(string(respondBody), -1)
		winString = submatchall1[0][0]
	} else if strings.Contains(submatchall[1][1], "tyvärr") {
		terminalColor = "\033[31m"
	}

	fmt.Printf("%s%s\n%s\n%s\n\n", terminalColor, submatchall[0][1], submatchall[1][1], winString)
}
