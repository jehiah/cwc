package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func main() {
	listReg := flag.Bool("list-regulations", false, "list all regulations")
	flag.Parse()

	if *listReg {
		for _, r := range allReg {
			fmt.Printf("%s,%s\n", r.Code, r.Description)
		}
		os.Exit(1)
	}

	var yyyymmdd, hhmm, license string
	fmt.Printf("Date (YYYYMMDD): ")
	fmt.Scanln(&yyyymmdd)

	if yyyymmdd == "" {
		yyyymmdd = time.Now().Format("20060102")
		fmt.Printf(" > using %s\n", yyyymmdd)
	}

	fmt.Printf("Time (HHMM): ")
	fmt.Scanln(&hhmm)

	dt, err := time.Parse("20060102 1504", fmt.Sprintf("%s %s", yyyymmdd, hhmm))
	if err != nil {
		log.Fatalf("err %s", err)
	}

	fmt.Printf("License Plate: ")
	fmt.Scanln(&license)

	baseDir := fmt.Sprintf("/Users/jehiah/Documents/cyclists_with_cameras/%s_%s_%s", yyyymmdd, hhmm, license)
	fmt.Printf("\tcreating ~/Documents/cyclists_with_cameras/%s_%s_%s\n", yyyymmdd, hhmm, license)

	err = os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		log.Fatalf("err %s", err)
	}

	vehicle := detectLicenseType(license)

	fmt.Printf("Where? ")
	reader := bufio.NewReader(os.Stdin)
	where, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("err %s", err)
	}
	where = strings.TrimSpace(where)

	f, err := os.Create(path.Join(baseDir, "notes.txt"))
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
	exec.Command("/Users/jehiah/bin/mate", baseDir).Run()
	exec.Command("/usr/bin/open", baseDir).Run()
}

func confirm() bool {
	var s string
	fmt.Scanf("%1s\n", &s)
	return s == "y" || s == "Y"
}
