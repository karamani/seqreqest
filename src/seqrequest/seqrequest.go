package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	debugMode    bool
	urlArg       string
	dataArg      string
	methodArg    string
	separatorArg string
)

func main() {
	app := cli.NewApp()
	app.Name = "fieldextract"
	app.Usage = "Retrieves the fields of data structures & prints them to stdout"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "debug mode",
			Destination: &debugMode,
		},
		cli.StringFlag{
			Name:        "url",
			Usage:       "request url",
			Destination: &urlArg,
		},
		cli.StringFlag{
			Name:        "data",
			Usage:       "parameters",
			Destination: &dataArg,
		},
		cli.StringFlag{
			Name:        "method",
			Usage:       "method",
			Value:       "GET",
			Destination: &methodArg,
		},
		cli.StringFlag{
			Name:        "separator",
			Usage:       "Output separator",
			Value:       "\t",
			Destination: &separatorArg,
		},
	}
	app.Action = func(c *cli.Context) {

		reader := bufio.NewReader(os.Stdin)
		inputString := ""
		for {
			bytes, hasMoreInLine, err := reader.ReadLine()
			if err != nil {
				if err != io.EOF {
					log.Fatalf("ERROR: %s\n", err.Error())
				}
				break
			}
			inputString += string(bytes)
			if !hasMoreInLine {

				debug(inputString)
				reqParams := strings.Split(inputString, separatorArg)

				urlString := urlArg
				dataString := dataArg
				for i, param := range reqParams {
					paramTpl := fmt.Sprintf("{%d}", i+1)
					urlString = strings.Replace(urlString, paramTpl, param, -1)
					dataString = strings.Replace(dataString, paramTpl, param, -1)
				}

				err := sendRequest(urlString, dataString)
				if err != nil {
					log.Println(err.Error())
				}

				inputString = ""
			}
		}
	}

	app.Run(os.Args)
}

func debug(msg string) {
	if debugMode {
		fmt.Println("DEBUG: " + msg)
	}
}

func sendRequest(urlString, dataString string) error {

	u, _ := url.ParseRequestURI(urlString)
	u.RawQuery = dataString

	req := &http.Request{
		Method: methodArg,
		URL:    u,
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}
