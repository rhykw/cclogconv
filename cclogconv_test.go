package cclogconv

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestRun_parseError(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cl := CCLogConv{Out: outStream, Err: errStream}
	args := strings.Split("cclogconv --broken", " ")

	exitCode := cl.Run(args[1:])
	if exitCode != ExitCodeError {
		t.Errorf("expected %d to eq %d", exitCode, ExitCodeError)
	}

	expected := "flag provided but not defined: -broken"
	if !strings.Contains(errStream.String(), expected) {
		t.Fatalf("expected %q to contain %q", errStream.String(), expected)
	}
}

func TestRun_nFlagError(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cl := CCLogConv{Out: outStream, Err: errStream}
	args := strings.Split("cclogconv -n", " ")

	exitCode := cl.Run(args[1:])
	if exitCode != ExitCodeError {
		t.Errorf("expected %d to eq %d", exitCode, ExitCodeError)
	}

	expected := "n option must be used with cc option."
	if !strings.Contains(errStream.String(), expected) {
		t.Fatalf("expected %q to contain %q", errStream.String(), expected)
	}
}

func TestRun_vFlagError(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cl := CCLogConv{Out: outStream, Err: errStream}
	args := strings.Split("cclogconv -v", " ")

	exitCode := cl.Run(args[1:])
	if exitCode != ExitCodeError {
		t.Errorf("expected %d to eq %d", exitCode, ExitCodeError)
	}

	expected := "v option must be used with cc option."
	if !strings.Contains(errStream.String(), expected) {
		t.Fatalf("expected %q to contain %q", errStream.String(), expected)
	}
}

func TestRun_filter_noflags(t *testing.T) {
	out := new(bytes.Buffer)
	in := bytes.NewBufferString("foo bar 182.22.59.229 aaa bbb\nboo moo 98.138.219.231 ccc def\n")
	filter := filter{
		in:  in,
		out: out,
	}

	err := filter.start("test/tmp/GeoLite2-Country_20181113/GeoLite2-Country.mmdb", "", false, false)
	if err != nil {
		t.Errorf("expected %v to eq nil", err)
	}

	expected := []byte("foo bar JP 182.22.59.229 aaa bbb\nboo moo US 98.138.219.231 ccc def\n")
	if bytes.Compare(expected, out.Bytes()) != 0 {
		t.Errorf("expected \n%#v to eq \n%#v", out.Bytes(), expected)
	}
}

func TestRun_filter_with_ccOpt(t *testing.T) {
	out := new(bytes.Buffer)
	in := bytes.NewBufferString("foo bar 182.22.59.229 aaa bbb\nboo moo 98.138.219.231 ccc def\n")

	filter := filter{
		in:  in,
		out: out,
	}

	err := filter.start("test/tmp/GeoLite2-Country_20181113/GeoLite2-Country.mmdb", "JP", false, false)
	if err != nil {
		t.Errorf("expected %v to eq nil", err)
	}

	expected := []byte("foo bar JP 182.22.59.229 aaa bbb\n")
	if bytes.Compare(expected, out.Bytes()) != 0 {
		t.Errorf("expected \n%#v to eq \n%#v", out.Bytes(), expected)
	}
}

func TestRun_filter_with_ccOpt_nFlag(t *testing.T) {
	out := new(bytes.Buffer)
	in := bytes.NewBufferString("foo bar 182.22.59.229 aaa bbb\nboo moo 98.138.219.231 ccc def\n")

	filter := filter{
		in:  in,
		out: out,
	}

	err := filter.start("test/tmp/GeoLite2-Country_20181113/GeoLite2-Country.mmdb", "JP", true, false)
	if err != nil {
		t.Errorf("expected %v to eq nil", err)
	}

	expected := []byte("foo bar 182.22.59.229 aaa bbb\n")
	if bytes.Compare(expected, out.Bytes()) != 0 {
		t.Errorf("expected \n%#v to eq \n%#v", out.Bytes(), expected)
	}
}

func TestRun_filter_with_ccOpt_vFlag(t *testing.T) {
	out := new(bytes.Buffer)
	in := bytes.NewBufferString("foo bar 182.22.59.229 aaa bbb\nboo moo 98.138.219.231 ccc def\n")

	filter := filter{
		in:  in,
		out: out,
	}

	err := filter.start("test/tmp/GeoLite2-Country_20181113/GeoLite2-Country.mmdb", "JP", false, true)
	if err != nil {
		t.Errorf("expected %v to eq nil", err)
	}

	expected := []byte("boo moo US 98.138.219.231 ccc def\n")
	if bytes.Compare(expected, out.Bytes()) != 0 {
		t.Errorf("expected \n%#v to eq \n%#v", out.Bytes(), expected)
	}
}

func TestRun_filter_with_ccOpt_nFlag_vFlag(t *testing.T) {
	out := new(bytes.Buffer)
	in := bytes.NewBufferString("foo bar 182.22.59.229 aaa bbb\nboo moo 98.138.219.231 ccc def\n")

	filter := filter{
		in:  in,
		out: out,
	}

	err := filter.start("test/tmp/GeoLite2-Country_20181113/GeoLite2-Country.mmdb", "JP", true, true)
	if err != nil {
		t.Errorf("expected %v to eq nil", err)
	}

	expected := []byte("boo moo 98.138.219.231 ccc def\n")
	if bytes.Compare(expected, out.Bytes()) != 0 {
		t.Errorf("expected \n%#v to eq \n%#v", out.Bytes(), expected)
	}
}

func TestRun_filter_manyLines(t *testing.T) {
	out := new(bytes.Buffer)
	in, _ := os.Open("test/tmp/test.log")

	filter := filter{
		in:  in,
		out: out,
	}

	err := filter.start("test/tmp/GeoLite2-Country_20181113/GeoLite2-Country.mmdb", "JP", false, false)
	if err != nil {
		t.Errorf("expected %v to eq nil", err)
	}

	expected, _ := ioutil.ReadFile("test/tmp/expected_JP.txt")
	if bytes.Compare(expected, out.Bytes()) != 0 {
		t.Errorf("expected \n%#v to eq \n%#v", out.Bytes(), expected)
	}
}
