package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func main() {
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

	fhv := isFHV(license)

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

	fmt.Fprintf(f, "%s %s %s %s\n", dt.Format("2006/01/02 3:04pm"), fhvStr(fhv), license, where)

	reg := getReg(fhv)

	fmt.Fprintf(f, "\n%s\n", SelectSample(reg, where))
	f.Close()

	fmt.Printf("done\n")

	var url string
	if fhv {
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

func fhvStr(fhv bool) string {
	if fhv {
		return "FHV"
	}
	return "Taxi"
}

func isFHV(license string) bool {
	if len(license) > 4 {
		return true
	}
	fmt.Printf("Taxi? y/n: ")
	taxi := confirm()
	return !taxi
}

func confirm() bool {
	var s string
	fmt.Scanf("%1s\n", &s)
	return s == "y" || s == "Y"
}
