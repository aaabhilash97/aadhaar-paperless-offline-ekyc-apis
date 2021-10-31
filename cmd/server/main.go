package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aaabhilash97/aadhaar_scrapper_apis/internal/appconfig"
	cmd "github.com/aaabhilash97/aadhaar_scrapper_apis/internal/cmd/server"
)

var gitCommit, gitTag string

func init() {
	version := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *version {
		msg := fmt.Sprintf(`Git Tag:      %s
Git commit:   %s`, gitTag, gitCommit)

		fmt.Println(msg)
		os.Exit(0)
	}
}

func main() {

	conf := appconfig.Init()
	if err := cmd.RunServer(conf); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	fmt.Println("Closing connections")
	conf.Close()
}
