package main

import (
	"os"
)

func ce(err error, msg ...string) bool {
	return cep(err, false, msg...)
}

func cep(err error, p bool, msg ...string) bool {
	if err != nil {
		if p {
			panic(err)
		} else {
			debug(msg...)
			debug(err.Error())
			return true
		}
	}
	return false
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
