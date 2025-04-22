package main

import (
	"flag"
	"log"
	"os"

	"github.com/FoxDenHome/shutdownd/util"
)

var hostnameFlag = flag.String("hostname", "", "Hostname to generate certificate for (default: current hostname)")
var certFileFlag = flag.String("file", "", "Path to the certificate file (required)")

func main() {
	log.Printf("ShutdownD certificate generator version %s", util.Commit())

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
