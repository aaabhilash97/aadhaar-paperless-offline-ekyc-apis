package main

import (
	"fmt"
	"os"

	"github.com/aaabhilash97/aadhaar-paperless-offline-ekyc-apis/internal/appconfig"
	cmd "github.com/aaabhilash97/aadhaar-paperless-offline-ekyc-apis/internal/cmd/server"
)

var gitCommit, gitTag string

func init() {
	os.Setenv("VERSION_INFO_GIT_COMMIT", gitCommit)
	os.Setenv("VERSION_INFO_GIT_TAG", gitTag)
}

func main() {

	conf := appconfig.Init()
	if err := cmd.RunServer(conf); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	fmt.Println("Closing connections")
	conf.Close()
}
