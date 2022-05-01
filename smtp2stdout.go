package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"os"
	"path/filepath"

	"github.com/mhale/smtpd"
	"github.com/google/uuid"
)

func writeMessage(fname string, data []byte) error  {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
	outfile := filepath.Join(pwd, "queue", fname + ".txt")
	outputerr := ioutil.WriteFile(outfile, data, os.ModePerm)
	if outputerr != nil {
		fmt.Println(os.Stderr, outputerr)
		os.Exit(1)
	}
	return nil
}
func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	uuidobj, _ := uuid.NewUUID()
	log.Printf("Received mail from %s for %s with subject %s write to %s", from, to[0], subject, uuidobj.String())
	writeMessage(uuidobj.String(), data)
	return nil
}

func main() {
	logFile, _ := os.OpenFile("logfile.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	defer logFile.Close()
	log.SetOutput(logFile)

	smtpd.ListenAndServe("172.30.112.1:25", mailHandler, "MyServerApp", "")
}
