package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func loadMigrations(dir, mode string) (migrations map[string]string, err error) {
	var files []os.FileInfo
	migrations = make(map[string]string)
	files, err = ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() == false {
			name := f.Name()
			if strings.Contains(name, "config") {
				continue
			}
			path := dir + "/" + f.Name()
			var data []byte
			if data, err = ioutil.ReadFile(path); err != nil {
				fatal(fmt.Sprintf("Error loading file: %s\n", err))
			}
			migrations[f.Name()] = string(data)
		}
	}
	return
}
