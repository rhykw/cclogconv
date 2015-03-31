package main

import (
    "fmt"
    "github.com/oschwald/geoip2-golang"
    "log"
    "net"
    "os"
    "bufio"
    "strings"
    "regexp"
)

func main() {

    db, err := geoip2.Open("./GeoLite2-Country.mmdb")
    if err != nil {
            log.Fatal(err)
    }
    defer db.Close()
    // If you are using strings that may be invalid, check that ip is not nil

    var sc = bufio.NewScanner(os.Stdin)
    c := 1
    for sc.Scan() {
        words := strings.Fields( sc.Text() )
        c += 1
        for _, v := range words {
            isIp , _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$",v)
            if isIp {
                ip := net.ParseIP( v )
                record, err := db.City(ip)
                if err != nil {
                        log.Fatal(err)
                }
                fmt.Printf("%s %v ", record.Country.IsoCode , v)
            }else{
                fmt.Printf("%s ",v)
            }
        }
        fmt.Printf("\n")
    }
    if err := sc.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }
}
