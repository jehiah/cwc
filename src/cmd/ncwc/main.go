package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"db"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	var err error
	listReg := flag.Bool("list-regulations", false, "list all regulations")
	short := flag.Bool("short", false, "short format")
	flag.Parse()

	if *listReg {
		for _, r := range allReg {
			desc := r.Description
			if *short && r.Short != "" {
				desc = r.Short
			}
			fmt.Printf("%s,%s\n", r.Code, desc)
		}
		os.Exit(1)
	}

	var yyyymmdd, hhmm, license string
	fmt.Printf("Date (YYYYMMDD) or Filename: ")
	fmt.Scanln(&yyyymmdd)

	var dt time.Time
	switch {
	case strings.HasPrefix(yyyymmdd, "/"):
		f, err := os.Open(yyyymmdd)
		if err != nil {
			log.Fatal(err)
		}
		x, err := exif.Decode(f)
		if err != nil {
			log.Fatalf("failed parsing exif from %s %s", yyyymmdd, err)
		}
		dt, err = x.DateTime()
		if err != nil {
			log.Fatalf("no EXIF date time %s", err)
		}
		fmt.Printf("> using EXIF time %s\n", dt.Format("2006/01/02 3:04pm"))
	case yyyymmdd == "":
		yyyymmdd = time.Now().Format("20060102")
		fmt.Printf(" > using %s\n", yyyymmdd)
		fallthrough
	default:
		fmt.Printf("Time (HHMM): ")
		fmt.Scanln(&hhmm)

		dt, err = time.Parse("20060102 1504", fmt.Sprintf("%s %s", yyyymmdd, hhmm))
		if err != nil {
			log.Fatalf("err %s", err)
		}
	}

	fmt.Printf("License Plate: ")
	fmt.Scanln(&license)

	vehicle := detectLicenseType(license)

	fmt.Printf("Where? ")
	reader := bufio.NewReader(os.Stdin)
	where, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("err %s", err)
	}
	where = strings.TrimSpace(where)

	complaint, err := db.Default.New(dt, license)
	if err != nil {
		log.Fatalf("err %s", err)
	}
	f, err := db.Default.Create(complaint)
	if err != nil {
		log.Fatalf("err %s", err)
	}

	fmt.Fprintf(f, "%s %s %s %s\n", dt.Format("2006/01/02 3:04pm"), vehicle, license, where)

	reg := getReg(vehicle)

	fmt.Fprintf(f, "\n%s\n", SelectSample(reg, where))
	f.Close()

	fmt.Printf("done\n")

	var url string
	if vehicle == FHV {
		url = "https://www1.nyc.gov/apps/311universalintake/form.htm?serviceName=TLC+FHV+Driver+Unsafe+Driving"
	} else {
		url = "https://www1.nyc.gov/apps/311universalintake/form.htm?serviceName=TLC+Taxi+Driver+Unsafe+Driving+Non-Passenger"
	}
	err = exec.Command("/usr/bin/open", "-a", "/Applications/Google Chrome.app/", url).Run()
	if err != nil {
		log.Printf("%s", err)
	}
	db.Default.Edit(complaint)
	db.Default.ShowInFinder(complaint)
}

func confirm() bool {
	var s string
	fmt.Scanf("%1s\n", &s)
	return s == "y" || s == "Y"
}
