package main

import (
	"fmt"
	"flag"
	"github.com/sourcegraph/go-ses"
	"io/ioutil"
	"os"
	"path/filepath"
)

func isFile(file string) bool {
	f, _ := os.Open(file)
	defer f.Close()
	if fi, err := f.Stat(); err == nil {
		if !fi.IsDir() && fi.Size() > 0 {
			return true
		}
	}
	return false
}

func main() {
	var (
		file string
		del bool
	)
	flag.StringVar(&file, "file", "", "Watcing file path")
	flag.BoolVar(&del, "del", false, "Deleting file flag")
	flag.Parse()

	if isFile(file) {
		fmt.Println("Exist:", file)
		filename := filepath.Base(file)
		content, _ := ioutil.ReadFile(file)

		res, err := ses.EnvConfig.SendEmail(
			"akiraak@gmail.com",
			"akiraak@gmail.com",
			fmt.Sprintf("Exist file: %s", filename),
			string(content))
		if err == nil {
			fmt.Printf("Sent email: %s...\n", res[:32])
		} else {
			fmt.Printf("Error sending email: %s\n", err)
		}

		if del {
			os.Remove(file)
		}
	} else {
		fmt.Println("Not exist:", file)
	}
}
