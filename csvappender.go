package helpers

import (
	"github.com/bocajim/helpers/log"
	"encoding/csv"
	"os"
)

func AppendToCsv(fileName string, fields []string) {
	fh,err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0660)
	if err!=nil {
		log.Printf(log.Warn,"Could not open file: "+err.Error())
		return
	}
	defer func() {
		fh.Sync()
		fh.Close()
	}()
	w:=csv.NewWriter(fh)
	err=w.Write(fields)	
	if err!=nil {
		log.Printf(log.Warn,"Could not write file: "+err.Error())
		return
	}
	w.Flush()
}