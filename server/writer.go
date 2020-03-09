package main

import (
	"os"
	"log"
	"context"
	"encoding/csv"
)

func initWriter(ctx context.Context, update chan bool, companies *AllCompanies, filename string) {
	filedata := readFile(filename)

	for _, i := range filedata {
		companies.add(load(i))
	}

	go func() {
		// goroutine 'owns' file writing
		file, err := os.OpenFile(filename, os.O_WRONLY, 0755)

		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		for {
			select {
			case <- ctx.Done():
				return

			case <- update:
				data := [][]string{}
				for _, company := range companies.get() {
					data = append(data, company.dump())
				}
				write(file, data)
			}
		}
	}()
}

func write(file *os.File, data [][]string) {
	file.Truncate(0)
	file.Seek(0, 0)
	
	w := csv.NewWriter(file)
	w.WriteAll(data)

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func readFile(filename string) [][]string {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(file)

	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return data
}