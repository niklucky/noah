package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
			path := dir + "/" + f.Name()
			var data []byte
			fmt.Println("Loading file: ", path)
			if data, err = ioutil.ReadFile(path); err != nil {
				fmt.Println("Error loading file: ", err)
			}
			migrations[f.Name()] = string(data)
		}
	}
	return
}
