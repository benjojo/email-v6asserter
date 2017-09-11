package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"regexp"
	"strings"
)

var v4r = regexp.MustCompile("(client-ip=|designates )\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}")
var v6r = regexp.MustCompile("(client-ip=|designates )(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))")

type statisticGroup struct {
	V6Count          int
	V6NotGoogleCount int
	V4Count          int
	TLSCount         int
	Total            int
}

func main() {
	inputfile := flag.String("mbox", "in.mbox", "The archive you want to read")
	flag.Parse()

	f, err := os.Open(*inputfile)

	if err != nil {
		log.Fatalf("Unable to open archive file, %s", err.Error())
	}

	input := make(chan io.Reader)
	go mboxreader(f, input)

	statisticMap := make(map[string]statisticGroup)
	for mailreader := range input {
		fmt.Print(".")
		msg, err := mail.ReadMessage(mailreader)
		if err != nil {
			log.Printf("Failed to read message (1) %s", err.Error())
			continue
		}

		OldEmailHostingTest := msg.Header.Get("X-Spam-Report")
		if strings.Contains(OldEmailHostingTest, "killersservers.co.uk") {
			// my email ( ben@benjojo.co.uk ) is a little complex, since it
			// was not always google apps. It was just a POP3 importer from my
			// shared hosting. So I'm going to filter those out.
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

			if !strings.Contains(InboundHeader, ".google.com designates") {
				group.V6NotGoogleCount++
			}

		} else if v4r.MatchString(InboundHeader) {
			group.V4Count++
		} else {
			if InboundHeader != "" {
				log.Printf("Failed to read message (3) - %s", InboundHeader)
			}
			continue
		}

		// try and find the TLS Header
		for _, v := range msg.Header["Received"] {
			if strings.Contains(v, "mx.google.com with ESMTPS") {
				group.TLSCount++
			}
		}
		group.Total++
		statisticMap[key] = group
	}

	fmt.Printf("Date,v6,v6notgoogle,v4,tls,Total\n")

	for k, v := range statisticMap {
		fmt.Printf("%s,%d,%d,%d,%d,%d\n", k, v.V6Count, v.V6NotGoogleCount, v.V4Count, v.TLSCount, v.Total)
	}
}

func mboxreader(r io.Reader, out chan io.Reader) {
	bio := bufio.NewReader(r)

	mail := ""
	bytes := 0
	toobig := false
	for {
		ln, _, err := bio.ReadLine()
		if err != nil {
			close(out)
			return
		}

		if strings.HasPrefix(string(ln), "From ") {
			// reset and send the reader down
			if toobig {
				log.Printf("Jumbo email! Was %d bytes / %d MB long", bytes, bytes/1024/1024)
			}
			nr := strings.NewReader(mail)
			out <- nr
			mail = ""
			bytes = 0
			toobig = false
			continue
		} else {
			if !toobig {
				mail += string(ln) + "\n"
			}
			bytes += len(ln)
		}

		if bytes > 1.5*1024*1024 {
			toobig = true
		}
	}
}
