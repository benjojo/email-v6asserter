package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"net/mail"
	"os"
	"regexp"
	"strings"
)

var v4r = regexp.MustCompile("client-ip=\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}")
var v6r = regexp.MustCompile("client-ip=(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))")

type statisticGroup struct {
	V6Count  int
	V4Count  int
	TLSCount int
	Total    int
}

func main() {
	inputfile := flag.String("archive", "example-UNSAFE.tar.gz", "The archive you want to read")
	flag.Parse()

	f, err := os.Open(*inputfile)

	if err != nil {
		log.Fatalf("Unable to open archive file, %s", err.Error())
	}

	gzr, err := gzip.NewReader(f)

	if err != nil {
		log.Fatalf("Unable to gzip decompress archive, %s", err.Error())
	}

	tarreader := tar.NewReader(gzr)
	statisticMap := make(map[string]statisticGroup)
	for {
		fileinfo, err := tarreader.Next()
		if err != nil {
			break
		}

		if !strings.Contains(fileinfo.Name, ".mbox") && !strings.Contains(fileinfo.Name, ".txt") {
			continue
		}

		msg, err := mail.ReadMessage(tarreader)
		if err != nil {
			log.Printf("Failed to read message (1) %s", err.Error())
			continue
		}

		InboundHeader := msg.Header.Get("Received-SPF")
		time, err := msg.Header.Date()
		if err != nil {
			log.Printf("Failed to read message (2) %s", err.Error())

			continue
		}

		key := time.Format("Jan 2006")
		group := statisticMap[key]
		if v6r.MatchString(InboundHeader) {
			group.V6Count++
		} else if v4r.MatchString(InboundHeader) {
			group.V4Count++
		} else {
			log.Printf("Failed to read message (3) - %s", InboundHeader)
			continue
		}
		group.Total++
		statisticMap[key] = group
	}

	fmt.Printf("Date,v6,v4,Total\n")

	for k, v := range statisticMap {
		fmt.Printf("%s,%d,%d,%d\n", k, v.V6Count, v.V4Count, v.Total)
	}
}
