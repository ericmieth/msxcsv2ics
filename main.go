package main

import (
	"encoding/csv"
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	timezone := "Europe/Berlin"
	now := time.Now()

	// input file
	csvfile, err := os.Open("data/input.csv")
	if err != nil {
		log.Println(err)
	}
	defer csvfile.Close()
	r := csv.NewReader(csvfile)

	// output file
	icsfile, err := os.Create("data/output.ics")
	if err != nil {
		log.Println(err)
	}
	defer icsfile.Close()

	// ics header

	icsfile.WriteString("BEGIN:VCALENDAR\n")
	icsfile.WriteString("VERSION:2.0\n")
	icsfile.WriteString("PRODID:msxcsv2ics\n")

	// parse each record

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		/*
			record[0] 	"Betreff"
			record[1] 	"Beginnt am"
			record[2] 	"Beginnt um"
			record[3] 	"Endet am"
			record[4] 	"Endet um"
			record[5] 	"Ganztägiges Ereignis"
			record[6] 	"Erinnerung Ein/Aus"
			record[7] 	"Erinnerung am"
			record[8] 	"Erinnerung um"
			record[9] 	"Besprechungsplanung"
			record[10] 	"Erforderliche Teilnehmer"
			record[11] 	"Optionale Teilnehmer"
			record[12] 	"Besprechungsressourcen"
			record[13] 	"Abrechnungsinformationen"
			record[14] 	"Beschreibung"
			record[15] 	"Kategorien"
			record[16] 	"Ort"
			record[17] 	"Priorität"
			record[18] 	"Privat"
			record[19] 	"Reisekilometer"
			record[20] 	"Vertraulichkeit"
			record[21] 	"Zeitspanne zeigen als"
		*/

		icsfile.WriteString("BEGIN:VEVENT\n")

		// DTSTART

		DTSTART, err := time.Parse("2.1.2006-15:04:05", record[1]+"-"+record[2])
		if err != nil {
			log.Println(err)
		}
		icsDTSTART := DTSTART.Format("20060102T150405")
		icsfile.WriteString("DTSTART;TZID=" + timezone + ":" + icsDTSTART + "\n")

		DTEND, err := time.Parse("2.1.2006-15:04:05", record[3]+"-"+record[4])
		if err != nil {
			log.Println(err)
		}
		icsDTEND := DTEND.Format("20060102T150405")

		icsfile.WriteString("DTDTEND;TZID=" + timezone + ":" + icsDTEND + "\n")
		icsfile.WriteString("DTSTAMP:" + now.Format("20060102T150405") + "\n")
		icsfile.WriteString("UID:" + strings.ToUpper(uuid.NewV4().String()) + "\n")
		icsfile.WriteString("CREATED:" + now.Format("20060102T150405") + "\n")
		icsfile.WriteString("LAST-MODIFIED:" + now.Format("20060102T150405") + "\n")
		icsfile.WriteString("SUMMARY:" + record[0] + ", " + record[16] + "\n")
		if len(record[14]) != 0 {
			icsfile.WriteString("DESCRIPTION:" + record[14] + "\n")
		}
		icsfile.WriteString("CLASS:PUBLIC\n")
		icsfile.WriteString("STATUS:CONFIRMED\n")
		icsfile.WriteString("TRANSP:OPAQUE\n")
		icsfile.WriteString("END:VEVENT\n")
		fmt.Println(record[0])

		icsfile.Sync()
	}
	icsfile.WriteString("END:VCALENDAR\n")
	icsfile.Sync()

}
