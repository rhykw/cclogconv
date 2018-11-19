package cclogconv

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

const (
	name    = "cclogconv"
	version = "1.3.9"
)

// These are the exit code definitions.
const (
	ExitCodeOK = iota
	ExitCodeError
)

// CCLogConv type
type CCLogConv struct {
	Out, Err io.Writer
}

// Run CCLogConv
func (cl CCLogConv) Run(args []string) int {

	// Define option flag parse
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(cl.Err)

	var (
		mmdbFilePath string
		selectCC     string
		notAdd       bool
		reverseCond  bool
	)
	flags.BoolVar(&notAdd, "n", false, "Not adding country code")
	flags.BoolVar(&reverseCond, "v", false, "Reverse condition for cc option")
	flags.StringVar(&selectCC, "cc", "", "Only displays line including this country's ip")
	flags.StringVar(&mmdbFilePath, "data", "/usr/share/GeoIP2/GeoLite2-Country.mmdb", "GeoIP2 Database Filename")

	if err := flags.Parse(args[0:]); err != nil {
		return ExitCodeError
	}

	if selectCC == "" {
		exitCode := ExitCodeOK
		if notAdd {
			fmt.Fprintln(cl.Err, "n option must be used with cc option.")
			exitCode = ExitCodeError
		}
		if reverseCond {
			fmt.Fprintln(cl.Err, "v option must be used with cc option.")
			exitCode = ExitCodeError
		}
		if exitCode != ExitCodeOK {
			return exitCode
		}
	}

	filter := filter{
		in:  os.Stdin,
		out: cl.Out,
	}

	if err := filter.start(mmdbFilePath, selectCC, notAdd, reverseCond); err != nil {
		fmt.Fprintf(cl.Err, "%s\n", err)
	}

	return ExitCodeOK
}

type filter struct {
	in  io.Reader
	out io.Writer
}

func (f filter) start(mmdbFilePath string, selectCc string, nFlag bool, vFlag bool) error {

	var lineBuf = ""

	db, err := geoip2.Open(mmdbFilePath)
	if err != nil {
		return err
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil

	re, _ := regexp.Compile("^(([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")

	var sc = bufio.NewScanner(f.in)
	for sc.Scan() {
		var ccMatchFlag = false
		lineBuf = ""
		words := strings.Fields(sc.Text())
		for _, word := range words {
			if re.MatchString(word) {
				ip := net.ParseIP(word)
				record, err := db.City(ip)
				if err != nil {
					return err
				}
				cc := record.Country.IsoCode
				if ccMatchFlag || cc == selectCc {
					ccMatchFlag = true
				}
				if cc == "" {
					cc = "-"
				}
				if nFlag == false {
					lineBuf += fmt.Sprintf("%s ", cc)
				}
				lineBuf = lineBuf + fmt.Sprintf("%s ", word)
			} else {
				lineBuf += fmt.Sprintf("%s ", word)
			}
		}

		if ((selectCc == "" || ccMatchFlag) && !vFlag) || (!(selectCc == "" || ccMatchFlag) && vFlag) {
			fmt.Fprintln(f.out, strings.TrimRight(lineBuf, " "))
		}
		lineBuf = ""
	}

	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}
