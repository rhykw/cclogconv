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
    )
    flag.Parse()
    var mmdbFilePath = fmt.Sprintf("%s", *optMmdbFilePath)

    db, err := geoip2.Open( mmdbFilePath )
    if err != nil {
            log.Fatal(err)
    }
    defer db.Close()
    // If you are using strings that may be invalid, check that ip is not nil

    re, _ := regexp.Compile("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$")

    var sc = bufio.NewScanner(os.Stdin)
    for sc.Scan() {
        words := strings.Fields( sc.Text() )
        for _, word := range words {
            if re.MatchString(word) {
                ip := net.ParseIP( word )
                record, err := db.City(ip)
                if err != nil {
                        log.Fatal(err)
                }
                cc := record.Country.IsoCode
                if cc == "" {
                    cc = "-"
                }
                fmt.Printf("%s %s ", cc , word)
            }else{
                fmt.Printf("%s ",word)
            }
        }
        fmt.Printf("\n")
    }
    if err := sc.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }
}
