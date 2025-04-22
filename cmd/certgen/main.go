package main

import (
	"flag"
	"os"

	"github.com/FoxDenHome/shutdownd/util"
)

var hostnameFlag = flag.String("hostname", "", "Hostname to generate certificate for (default: current hostname)")
var certFileFlag = flag.String("file", "", "Path to the certificate file (required)")

func main() {
	flag.Parse()

	if *certFileFlag == "" {
		flag.Usage()
		os.Exit(1)
		return
	}

	err := util.GenerateSelfSignedCert(*certFileFlag, *hostnameFlag)
	if err != nil {
		panic(err)
	}
}
