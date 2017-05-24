package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {

	var (
		optMmdbFilePath = flag.String("data", "/usr/share/GeoIP2/GeoLite2-Country.mmdb", "GeoIP2 Database Filename")
		selectCc        = flag.String("cc", "", "Only displays line including this country's ip")
		nFlag           = flag.Bool("n", false, "Not adding country code")
		vFlag           = flag.Bool("v", false, "Reverse condition for cc option")
	)
	flag.Parse()
	var mmdbFilePath = fmt.Sprintf("%s", *optMmdbFilePath)
	var lineBuf = ""

	if *selectCc == "" {
		if *nFlag {
			fmt.Fprintln(os.Stderr, "n option must be used with cc option.")
		}
		if *vFlag {
			fmt.Fprintln(os.Stderr, "v option must be used with cc option.")
		}
		if *nFlag || *vFlag {
			os.Exit(1)
		}
	}

	db, err := geoip2.Open(mmdbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil

	re, _ := regexp.Compile("^(([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")

	var sc = bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		var ccMatchFlag = false
		lineBuf = ""
		words := strings.Fields(sc.Text())
		for _, word := range words {
			if re.MatchString(word) {
				ip := net.ParseIP(word)
				record, err := db.City(ip)
				if err != nil {
					log.Fatal(err)
				}
				cc := record.Country.IsoCode
				if ccMatchFlag || cc == *selectCc {
					ccMatchFlag = true
				}
				if cc == "" {
					cc = "-"
				}
				if *nFlag == false {
					lineBuf += fmt.Sprintf("%s ", cc)
				}
				lineBuf = lineBuf + fmt.Sprintf("%s ", word)
			} else {
				lineBuf += fmt.Sprintf("%s ", word)
			}
		}

		if ((*selectCc == "" || ccMatchFlag) && !*vFlag) || (!(*selectCc == "" || ccMatchFlag) && *vFlag) {
			fmt.Println(lineBuf)
		}
		lineBuf = ""
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
