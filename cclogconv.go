package main

import (
    "fmt"
    "flag"
    "github.com/oschwald/geoip2-golang"
    "log"
    "net"
    "os"
    "bufio"
    "strings"
    "regexp"
)

func main() {

    var (
        optMmdbFilePath = flag.String("data", "/usr/share/GeoIP/GeoLite2-Country.mmdb", "GeoIP2 Database Filename")
        selectCc        = flag.String("cc"  , "", "Only displays line including this country's ip")
    )
    flag.Parse()
    var mmdbFilePath = fmt.Sprintf("%s", *optMmdbFilePath)
    var lineBuf = ""

    db, err := geoip2.Open( mmdbFilePath )
    if err != nil {
            log.Fatal(err)
    }
    defer db.Close()
    // If you are using strings that may be invalid, check that ip is not nil

    re, _ := regexp.Compile("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$")

    var sc = bufio.NewScanner(os.Stdin)
    for sc.Scan() {
        var ccMatchFlag = false
        lineBuf = ""
        words := strings.Fields( sc.Text() )
        for _, word := range words {
            if re.MatchString(word) {
                ip := net.ParseIP( word )
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
                lineBuf = fmt.Sprintf("%s %s ", cc , word)
            }else{
                lineBuf = fmt.Sprintf("%s ",word)
            }
        }

        if *selectCc == "" || ccMatchFlag {
            fmt.Printf(lineBuf + "\n")
        }
        lineBuf = ""
    }
    if err := sc.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }
}
