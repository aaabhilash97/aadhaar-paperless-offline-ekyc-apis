package main

import (
	"fmt"
	"os"

	"github.com/aaabhilash97/aadhaar_scrapper_apis/internal/appconfig"
	cmd "github.com/aaabhilash97/aadhaar_scrapper_apis/internal/cmd/server"
)

func init() {

}

func main() {

	conf := appconfig.Init()
	if err := cmd.RunServer(conf); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	fmt.Println("Closing connections")
	conf.Close()
}
