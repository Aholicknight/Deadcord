package core

import (
	"io"
	"log"
	"os"
)

func InitLogger() {
	log_file, err := os.OpenFile("./deadcord.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	out_file := io.MultiWriter(os.Stdout, log_file)
	log.SetOutput(out_file)
}
