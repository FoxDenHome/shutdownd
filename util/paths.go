package util

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var configDir = ""

func ExePath() (string, error) {
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

func GetConfigDir(logger Logger) (string, error) {
	if configDir != "" {
		return configDir, nil
	}

	exePath, err := ExePath()
	if err != nil {
		return "", err
	}
	configDir = filepath.Dir(exePath)
	logger.Info(1, fmt.Sprintf("Config dir %s from EXE path %s", configDir, exePath))
	return configDir, err
}
