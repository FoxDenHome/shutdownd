package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func exePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}

	var fi fs.FileInfo
	fi, err = os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err = os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}

func getConfigDir() (string, error) {
	exepath, err := exePath()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exepath), nil
}
