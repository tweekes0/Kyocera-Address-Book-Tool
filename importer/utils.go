package importer

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func SetupCSV(t *testing.T, data [][]string) (io.ReadWriter, func()) {
	
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("cannot create test file: %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()	
	
	w.WriteAll(data)
	

	r, err := os.Open(f.Name())
	if err != nil {
		log.Fatalf("cannot open csv: %v", err)
	}
	
	teardown := func() {
		defer r.Close()
	}


	return r, teardown
}