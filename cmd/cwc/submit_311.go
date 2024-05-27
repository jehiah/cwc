package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/jehiah/cwc/db"
	"github.com/jehiah/cwc/input"
	"github.com/jehiah/cwc/internal/complaint"
	"github.com/jehiah/cwc/internal/reg"
	"github.com/spf13/cobra"
	"github.com/ysmood/gson"
)

func submitComplaint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "submit Complaint",
		Run: func(cmd *cobra.Command, args []string) {
			err := runSubmitComplaint(loadDB(cmd.Flags().GetString("db")), args...)
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}
	cmd.Flags().String("db", string(db.Default), "DB path")
	return cmd
}

func runSubmitComplaint(d db.ReadWrite, args ...string) error {
	var err error
	var query string
	if len(args) > 0 {
		query = args[0]
	}
	if query == "" {
		query, err = input.Ask("Search for?", "")
		if err != nil {
			log.Fatalf("%s", err)
		}
		if query == "" {
			log.Fatalf("missing query")
		}
	} else {
		fmt.Printf("Searching for: %q\n", query)
	}

	// files, err := d.Find(query)
	// if err != nil {
	// 	return err
	// }
	// if len(files) == 0 {
	// 	fmt.Printf("no files found\n")
	// 	return nil
	// }
	// if len(files) > 1 {
	// 	return fmt.Errorf("too many matches %d", len(files))
	// }

	fc, err := d.FullComplaint(complaint.Complaint(query))
	if err != nil {
		return err
	}
	return Submit(d, fc)
}

func Submit(d db.ReadOnly, fc *complaint.FullComplaint) error {
	hd, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	var u string
	if path, exists := launcher.LookPath(); exists {
		u = launcher.New().Headless(false).Bin(
			path).UserDataDir(
			filepath.Join(hd, ".cache/cwc_chrome_profile")).MustLaunch()
	}
	browser := rod.New().ControlURL(u).MustConnect()

	var url, title string
	switch fc.VehicleType {
	case reg.FHV.String():
		url = "https://portal.311.nyc.gov/article/?kanumber=KA-01244"
		title = "Car Service Complaint"
	default:
		url = "https://portal.311.nyc.gov/article/?kanumber=KA-01241"
		title = "Taxi Complaint"
	}

	fmt.Printf("> opening %s\n", url)
	page := browser.MustPage(url)

	log.Printf("waiting for page %q", title)
	page.MustWaitLoad()
	log.Printf("on %q", page.MustInfo().Title)
	// page.MustWait(fmt.Sprintf(`() => document.title === %q`, title))

	loginLink := page.MustElement(".login-area > a").MustText()
	log.Printf("login link %#v", loginLink)

	// if unauthenticated, go to the sign in page and wait for the browser to return to the desired URL
	if loginLink == "Sign In" {
		fmt.Printf("> Sign In\n")
		page.MustElement(".login-area > a").MustClick()

		fmt.Printf("> waiting for sign in and redirect to %s", url)
		page.MustWaitNavigation()
		// WaitRequestIdle  ?

		input.Ask("continue?", "")
		// var currentAddress string
		// for currentAddress != url {
		// 	time.Sleep(500 * time.Millisecond)
		// 	if err := chromedp.Run(ctx,
		// 		chromedp.Location(&currentAddress),
		// 	); err != nil {
		// 		return err
		// 	}
		// 	// fmt.Printf("current url is %q", currentAddress)
		// }

		// TODO: re-confirm title
		// log.Printf("checking for title %q", title)
		// if err := chromedp.Run(ctx); // chromedp.WaitVisible(fmt.Sprintf(`//h1[text()=%q]]`, title)),
		// err != nil {
		// 	return err
		// }

		fmt.Printf("> authenticated")
	}
	time.Sleep(time.Millisecond * 150)
	switch fc.VehicleType {
	case reg.FHV.String():
		err = SubmitFHV(page, d, fc)
	default:
		err = SubmitTaxi(page, d, fc)
	}
	if err != nil {
		log.Printf("err %s", err)

		img, _ := page.Screenshot(true, &proto.PageCaptureScreenshot{
			Format:  proto.PageCaptureScreenshotFormatPng,
			Quality: gson.Int(90),
			Clip: &proto.PageViewport{
				X:      0,
				Y:      0,
				Width:  300,
				Height: 200,
				Scale:  1,
			},
			FromSurface: true,
		})

		if len(img) > 0 {
			log.Printf("> screnshot saved as fullScreenshot.png")
			os.WriteFile("fullScreenshot.png", img, 0o644)
		}
		input.Ask("Close?", "")

		return err
	}
	input.Ask("Done?", "")
	return nil
}

func SubmitTaxi(page *rod.Page, d db.ReadOnly, fc *complaint.FullComplaint) error {
	var err error
	log.Printf("SubmitTaxi")
	// click = "You were not a passenger"
	sel := `//h5/a[text()[contains(.,"NOT a passenger")]]`
	err = page.MustElementX(sel).ScrollIntoView()
	if err != nil {
		return err
	}
	page.MustElementX(sel).MustClick()
	time.Sleep(time.Second * 2)

	if err = checkStep(page, "What"); err != nil {
		return err
	}

	log.Printf("> filling fields")
	page.MustElement(`#n311_attendhearing_1`).MustClick()                                          // yes
	page.MustElement(`#n311_coloroftaxi_1`).MustClick()                                            // yellow
	page.MustElement(`#n311_problemdetailid_select`).Input("87779989-ee94-e811-a961-000d3a1993e0") // Driver complaint - non passenger
	time.Sleep(time.Second)
	page.MustElement(`#n311_additionaldetailsid_select`).Timeout(time.Second).Input("eb4e791a-374e-e811-a94d-000d3a360e00") // Unsafe Driving - Non-Passenger
	page.MustElement(`#n311_taximedallionnumber_name`).Timeout(time.Second).Input(fc.License)
	// page.MustElement(`#n311_datetimeobserved`).Timeout(time.Second).Input(fc.Time.Format("2006-01-02T15:04:05.000Z"))
	page.MustElement(`#n311_datetimeobserved_datepicker_description`).Timeout(time.Second).Input(fc.Time.Format("1/2/2006 3:04 PM"))

	// if err := chromedp.Run(ctx,
	// 	chromedp.Click(`#n311_attendhearing_1`, chromedp.ByID),                                                       // yes
	// 	chromedp.Click(`#n311_coloroftaxi_1`, chromedp.ByID, chromedp.NodeVisible),                                   // yellow
	// 	chromedp.SetValue(`#n311_problemdetailid_select`, "87779989-ee94-e811-a961-000d3a1993e0", chromedp.ByID),     // Driver complaint - non passenger
	// 	chromedp.SetValue(`#n311_additionaldetailsid_select`, "eb4e791a-374e-e811-a94d-000d3a360e00", chromedp.ByID), // Unsafe Driving - Non-Passenger
	// 	chromedp.SetValue(`#n311_taximedallionnumber_name`, fc.License, chromedp.ByID),
	// 	chromedp.SetValue(`#n311_datetimeobserved`, fc.Time.Format("2006-01-02T15:04:05.000Z"), chromedp.ByID),
	// 	chromedp.SetValue(`#n311_datetimeobserved_datepicker_description`, fc.Time.Format("1/2/2006 3:04 PM"), chromedp.ByID),
	// ); err != nil {
	// 	return err
	// }

	// upload attachments; video first
	for _, file := range uploads(d.FullPath(fc.Complaint), fc) {
		f := filepath.Join(d.FullPath(fc.Complaint), file)
		err := uploadFile(page, f)
		if err != nil {
			return err
		}
	}

	page.MustElement(`#n311_description`).MustInput(fc.Description)
	fmt.Println("press Next when ready")

	err = upload311Location(page, fc)
	if err != nil {
		return err
	}

	return nil
}

func SubmitFHV(page *rod.Page, d db.ReadOnly, fc *complaint.FullComplaint) error {
	log.Printf("SubmitFHV")
	sel := `//button/div/div[text()='Driver Complaint']`
	err := page.MustElementX(sel).ScrollIntoView()
	if err != nil {
		return err
	}
	page.MustElementX(sel).MustClick()
	// if err := chromedp.Run(ctx,
	// 	chromedp.ScrollIntoView(sel),
	// 	chromedp.Sleep(time.Millisecond*100),
	// 	chromedp.Click(sel),
	// 	chromedp.Sleep(time.Millisecond*100),
	// ); err != nil {
	// 	return err
	// }
	time.Sleep(time.Millisecond * 100)
	start := `//h5/a[text()[contains(.,"not a passenger")]]`
	log.Printf("click start %s", start)
	err = page.MustElementX(start).ScrollIntoView()
	if err != nil {
		return err
	}
	page.MustElementX(start).MustClick()
	// if resp, err := chromedp.RunResponse(ctx,
	// 	chromedp.ScrollIntoView(start),
	// 	chromedp.Sleep(time.Millisecond*10),
	// 	chromedp.Click(start),
	// ); err != nil {
	// 	return err
	// } else {
	// 	fmt.Println("RunResponse status code:", resp.Status)
	// }
	time.Sleep(time.Second)
	waitStep(page, "What")

	log.Printf("> filling fields")
	page.MustElement(`#n311_attendhearing_1`).MustClick() // yes
	page.MustElement(`#n311_coloroftaxi_2`).MustClick()   // yellow
	page.MustElement(`#n311_licensenumber`).ScrollIntoView()
	page.MustElement(`#n311_licensenumber`).MustInput(fc.License)
	log.Printf("n311_datetimeobserved")
	// n311_datetimeobserved is not visible
	// el := page.MustElement(`#n311_datetimeobserved`)
	// page.Context(el.GetContext()).InsertText(fc.Time.Format("2006-01-02T15:04:05.000Z"))
	// page.MustElement(`#n311_datetimeobserved`).MustInput(fc.Time.Format("2006-01-02T15:04:05.000Z"))
	page.MustElement(`#n311_datetimeobserved_datepicker_description`).MustInput(fc.Time.Format("1/2/2006 3:04 PM"))

	// if err := chromedp.Run(ctx,
	// 	chromedp.WaitVisible(`#n311_attendhearing_1`, chromedp.ByID),
	// 	chromedp.Click(`#n311_attendhearing_1`, chromedp.ByID),
	// 	chromedp.Click(`#n311_coloroftaxi_2`, chromedp.ByID, chromedp.NodeVisible),
	// 	chromedp.ScrollIntoView(`n311_licensenumber`, chromedp.ByID),
	// 	chromedp.SetValue(`#n311_licensenumber`, fc.License, chromedp.ByID),
	// 	chromedp.SetValue(`#n311_datetimeobserved`, fc.Time.Format("2006-01-02T15:04:05.000Z"), chromedp.ByID),
	// 	chromedp.SetValue(`#n311_datetimeobserved_datepicker_description`, fc.Time.Format("1/2/2006 3:04 PM"), chromedp.ByID),
	// ); err != nil {
	// 	return err
	// }

	// upload attachments; video first
	for _, file := range uploads(d.FullPath(fc.Complaint), fc) {
		f := filepath.Join(d.FullPath(fc.Complaint), file)
		err := uploadFile(page, f)
		if err != nil {
			return err
		}
	}

	// do this last so we can wait on a 'RunResponse'
	page.MustElement(`#n311_description`).MustInput(fc.Description)
	fmt.Println("press Next when ready")

	err = upload311Location(page, fc)
	if err != nil {
		return err
	}

	return nil
}

func uploads(path string, fc *complaint.FullComplaint) []string {
	var out []string
	for _, f := range fc.Videos {
		// check for file size
		fileName := filepath.Join(path, f)
		fi, err := os.Stat(fileName)
		if err != nil {
			log.Printf("error stat %q %s", f, err)
		} else {
			if fi.Size() > 74000000 {
				log.Printf("skipping %q, too big %d", f, fi.Size())
				continue
			}
		}
		log.Printf("%s size %d", filepath.Base(f), fi.Size())
		out = append(out, f)
	}
	for _, p := range fc.Photos {
		switch filepath.Ext(strings.ToLower(p)) {
		case ".heic":
		default:
			out = append(out, p)
		}
	}
	if len(out) > 3 {
		out = out[:3]
	}
	return out
}

func uploadFile(page *rod.Page, f string) error {
	log.Printf("uploading %q", f)
	page.MustElement(`#attachments-addbutton`).MustScrollIntoView()
	page.MustElement(`#attachments-addbutton`).MustClick()
	page.MustElement(`input[type="file"]`).MustWaitVisible()
	page.MustElement(`input[type="file"]`).MustSetFiles(f)
	page.MustElementX(`//button[text() = 'Add Attachment']`).MustClick()
	sel := fmt.Sprintf(`//p[contains(@class, 'attachmentTitle') and text()='%s']`, filepath.Base(f))
	page.MustElementX(sel).MustWaitVisible()

	// ctx, cancel := context.WithTimeout(ctx, timeout)
	// defer cancel()
	// if err := chromedp.Run(ctx,
	// 	chromedp.ScrollIntoView(`#attachments-addbutton`, chromedp.ByID),
	// 	chromedp.Click(`#attachments-addbutton`, chromedp.ByID),      // add attachment
	// 	chromedp.WaitVisible(`input[type="file"]`, chromedp.ByQuery), // wait for modal
	// 	chromedp.SetUploadFiles(`input[type="file"]`, []string{f}, chromedp.ByQuery),
	// 	chromedp.Click(`//button[text() = 'Add Attachment']`), // upload
	// 	// TODO: check for error uploading?
	// 	chromedp.WaitVisible(fmt.Sprintf(`//p[contains(@class, 'attachmentTitle') and text()='%s']`, filepath.Base(f))), // the final results table
	// ); err != nil {
	// 	return err
	// }
	return nil
}

func waitStep(page *rod.Page, expectedStep string) {
	fmt.Printf("waiting for step %q\n", expectedStep)
	sel := fmt.Sprintf(`() => document.querySelector(".step.active")?.innerText === %q`, expectedStep)
	page.Timeout(time.Minute).MustWait(sel)
}

func checkStep(page *rod.Page, expected string) error {
	step := strings.TrimSpace(page.MustElement(".step.active").MustText())
	fmt.Println("on step", step)
	if step != expected {
		return fmt.Errorf("on step %q expected %q", step, expected)
	}
	return nil
}

func upload311Location(page *rod.Page, fc *complaint.FullComplaint) error {
	// https://portal.311.nyc.gov/sr-step/?id=4a3484e4-cd1e-ee11-a81c-6045bdb05de8&stepid=9241458f-fb0d-e811-8127-1458d04d2538
	var err error

	time.Sleep(time.Second * 2) // TODO: wait for spinner to go away
	waitStep(page, "Where")

	time.Sleep(time.Second)
	page.MustElement(`#n311_locationtypeid_select`).Input("a7c99a56-e64e-e811-a951-000d3a3606de") // street
	time.Sleep(time.Second)
	fmt.Println("> filling location", fc.Location)
	page.MustElement(`#n311_additionallocationdetails`).MustInput(fc.Location)

	// fmt.Println("> filling location ", fc.Location)
	// if resp, err := chromedp.RunResponse(ctx,
	// 	chromedp.SetValue(`#n311_locationtypeid_select`, "a7c99a56-e64e-e811-a951-000d3a3606de", chromedp.ByID), // street
	// 	chromedp.Sleep(time.Second),
	// 	// Bug - the RunResponse is stalling here
	// 	chromedp.SetValue(`#n311_additionallocationdetails`, fc.Location, chromedp.ByID),
	// 	chromedp.Sleep(time.Millisecond*100),
	// ); err != nil {
	// 	return err
	// } else {
	// 	fmt.Println("location RunResponse status code:", resp.Status)
	// }

	if fc.Address != "" {
		fmt.Println("> filling address", fc.Address)
		err = page.MustElement(`#SelectAddressWhere`).Timeout(time.Second * 5).WaitVisible()
		if err != nil {
			return err
		}
		err = page.MustElement(`#SelectAddressWhere`).Timeout(time.Second*5).Click(proto.InputMouseButtonLeft, 1)
		if err != nil {
			return err
		}
		// page.MustElement(`#SelectAddressWhere`).Timeout(time.Second * 5).MustClick()
		time.Sleep(time.Second)
		page.MustElement(`#address-search-box-input`).Timeout(time.Second * 5).MustWaitVisible()
		page.MustElement(`#address-search-box-input`).MustInput(fc.Address)
		page.MustElement(`#ui-id-2`).Timeout(time.Second * 2).MustWaitVisible()
		page.MustElement(`#ui-id-2`).MustClick()
		time.Sleep(time.Millisecond * 100)
		page.MustElement(`#SelectAddressMap`).Timeout(time.Millisecond * 100).MustScrollIntoView()
		page.MustElement(`#SelectAddressMap`).Timeout(time.Millisecond * 100).MustClick()
	}

	// if fc.Address != "" {
	// 	fmt.Println("> filling address ", fc.Address)
	// 	if resp, err := chromedp.RunResponse(ctx,
	// 		chromedp.Sleep(time.Second),
	// 		chromedp.Click(`#SelectAddressWhere`, chromedp.ByID),
	// 		chromedp.Sleep(time.Second),
	// 		chromedp.WaitVisible(`#address-search-box-input`, chromedp.ByID),
	// 		chromedp.SetValue(`#address-search-box-input`, fc.Address, chromedp.ByID), // this is search pre-population
	// 		chromedp.Sleep(time.Second),
	// 		// wait for #ui-id-1 > li > div
	// 		chromedp.WaitVisible(`#ui-id-2`, chromedp.ByID),
	// 		chromedp.Click(`#ui-id-2`, chromedp.ByID, chromedp.NodeVisible),
	// 		chromedp.Sleep(time.Millisecond*100),
	// 		chromedp.ScrollIntoView(`#SelectAddressMap`, chromedp.ByID),
	// 		chromedp.Click(`#SelectAddressMap`, chromedp.ByID, chromedp.NodeVisible),
	// 	); err != nil {
	// 		return err
	// 	} else {
	// 		fmt.Println("address RunResponse status code:", resp.Status)
	// 	}

	// 	// ui-id-1 > li > div
	// 	// document.getElementById('SelectAddressMap').click()

	// }
	// .address-picker-input ?
	return nil
}

func save311(page *rod.Page) {
	// #n311_name .value

}
