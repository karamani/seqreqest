package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/karamani/iostreams"
)

var (
	debugMode    bool
	urlArg       string
	dataArg      string
	methodArg    string
	separatorArg string
	throughMode  bool
	fakeMode     bool
)

func main() {
	app := cli.NewApp()
	app.Name = "fieldextract"
	app.Usage = "Retrieves the fields of data structures & prints them to stdout"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "debug mode",
			EnvVar:      "SEQREQUEST_DEBUG",
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
		cli.BoolFlag{
			Name:        "through",
			Usage:       "translate stdin to stdout",
			Destination: &throughMode,
		},
		cli.BoolFlag{
			Name:        "fake",
			Usage:       "fake mode (without http requests)",
			Destination: &fakeMode,
		},
	}

	app.Action = func(c *cli.Context) {

		// this func's called for each stdin's row
		process := func(row []byte) error {

			debug(string(row))

			reqParams := strings.Split(string(row), separatorArg)

			urlString := urlArg
			dataString := dataArg
			for i, param := range reqParams {
				paramTpl := fmt.Sprintf("{%d}", i+1)
				urlString = strings.Replace(urlString, paramTpl, param, -1)
				dataString = strings.Replace(dataString, paramTpl, param, -1)
			}

			debug("Url: %s", urlString)
			debug("Data: %s", dataString)

			if fakeMode {
				return nil
			}

			err := sendRequest(urlString, dataString)
			if err != nil {
				log.Println(err.Error())
			}

			return nil
		}

		iostreams.ThroughMode = throughMode

		err := iostreams.ProcessStdin(process)
		if err != nil {
			log.Panicln(err.Error())
		}
	}

	app.Run(os.Args)
}

func sendRequest(urlString, dataString string) error {

	u, _ := url.ParseRequestURI(urlString)
	u.RawQuery = dataString

	debug(fmt.Sprintf("%#v", u))

	req := &http.Request{
		Method: methodArg,
		URL:    u,
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	debug("Response: %s\n", string(body))

	return nil
}

func debug(format string, args ...interface{}) {
	if debugMode {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}
