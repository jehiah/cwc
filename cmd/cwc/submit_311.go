package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/jehiah/cwc/db"
	"github.com/jehiah/cwc/input"
	"github.com/jehiah/cwc/internal/complaint"
	"github.com/jehiah/cwc/internal/reg"
	"github.com/spf13/cobra"
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
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.CombinedOutput(os.Stderr), // stdout and stderr from the browser
		// Headless,

		// After Puppeteer's default behavior.
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		// chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees"),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("force-color-profile", "srgb"),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("safebrowsing-disable-auto-update", true),
		chromedp.Flag("enable-automation", true),
		chromedp.Flag("password-store", "basic"),
		chromedp.Flag("use-mock-keychain", true),

		chromedp.UserDataDir(filepath.Join(hd, ".cache/cwc_chrome_profile")),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	defer chromedp.Cancel(ctx)

	// set a parent timeout so we bound our total time in case something hangs
	ctx, cancel = context.WithTimeout(ctx, time.Minute*10)
	defer cancel()

	var url, title string
	switch fc.VehicleType {
	case reg.FHV.String():
		url = "https://portal.311.nyc.gov/article/?kanumber=KA-01244"
		title = "Car Service Complaint"
	default:
		url = "https://portal.311.nyc.gov/article/?kanumber=KA-01241"
		title = "Taxi Complaint"
	}

	// open start page
	fmt.Printf("> opening %s\n", url)
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
	); err != nil {
		return err
	}

	// get the contents of the login area.
	log.Printf("waiting for page %q", title)
	var loginLink string
	if err := chromedp.Run(ctx,
		// chromedp.WaitVisible(fmt.Sprintf(`//h1[text()=%q]]`, title)),
		chromedp.TextContent(`.login-area > a`, &loginLink, chromedp.ByQuery),
	); err != nil {
		return err
	}
	log.Printf("login link %#v", loginLink)

	// if unauthenticated, go to the sign in page and wait for the browser to return to the desired URL
	if loginLink == "Sign In" {
		fmt.Printf("> Sign In\n")
		if resp, err := chromedp.RunResponse(ctx,
			chromedp.Click(`.login-area > a`, chromedp.NodeVisible, chromedp.ByQuery),
		); err != nil {
			return err
		} else {
			fmt.Println("RunResponse status code:", resp.Status)
		}

		fmt.Printf("> waiting for sign in and redirect to %s", url)
		var currentAddress string
		for currentAddress != url {
			time.Sleep(500 * time.Millisecond)
			if err := chromedp.Run(ctx,
				chromedp.Location(&currentAddress),
			); err != nil {
				return err
			}
			// fmt.Printf("current url is %q", currentAddress)
		}

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
		err = SubmitFHV(ctx, d, fc)
	default:
		err = SubmitTaxi(ctx, d, fc)
	}
	if err != nil {
		log.Printf("err %s", err)
		var buf []byte
		chromedp.Run(ctx,
			chromedp.FullScreenshot(&buf, 90),
		)
		if len(buf) > 0 {
			log.Printf("> screnshot saved as fullScreenshot.png")
			os.WriteFile("fullScreenshot.png", buf, 0o644)
		}
		input.Ask("Close?", "")

		return err
	}
	input.Ask("Done?", "")
	return nil
}

func SubmitTaxi(ctx context.Context, d db.ReadOnly, fc *complaint.FullComplaint) error {
	// click = "You were not a passenger"

	if resp, err := chromedp.RunResponse(ctx,
		chromedp.Click(`//h5/a[text()[contains(.,"NOT a passenger")]]`),
	); err != nil {
		return err
	} else {
		fmt.Println("RunResponse status code:", resp.Status)
	}

	time.Sleep(time.Second)

	if err := checkStep(ctx, "What"); err != nil {
		return err
	}

	log.Printf("> filling fields")
	if err := chromedp.Run(ctx,
		chromedp.Click(`#n311_attendhearing_1`, chromedp.ByID),                                                       // yes
		chromedp.Click(`#n311_coloroftaxi_1`, chromedp.ByID, chromedp.NodeVisible),                                   // yellow
		chromedp.SetValue(`#n311_problemdetailid_select`, "87779989-ee94-e811-a961-000d3a1993e0", chromedp.ByID),     // Driver complaint - non passenger
		chromedp.SetValue(`#n311_additionaldetailsid_select`, "eb4e791a-374e-e811-a94d-000d3a360e00", chromedp.ByID), // Unsafe Driving - Non-Passenger
		chromedp.SetValue(`#n311_taximedallionnumber_name`, fc.License, chromedp.ByID),
		chromedp.SetValue(`#n311_datetimeobserved`, fc.Time.Format("2006-01-02T15:04:05.000Z"), chromedp.ByID),
		chromedp.SetValue(`#n311_datetimeobserved_datepicker_description`, fc.Time.Format("1/2/2006 3:04 PM"), chromedp.ByID),
	); err != nil {
		return err
	}

	// upload attachments; video first
	for _, file := range uploads(fc) {
		f := filepath.Join(d.FullPath(fc.Complaint), file)
		err := uploadFile(ctx, f, time.Second*90)
		if err != nil {
			return err
		}
	}

	// do this last so we can wait on a 'RunResponse'
	fmt.Println("waiting for Next")
	if resp, err := chromedp.RunResponse(ctx,
		chromedp.SetValue(`#n311_description`, fc.Description, chromedp.ByID),
	); err != nil {
		return err
	} else {
		fmt.Println("RunResponse status code:", resp.Status)
	}

	err := upload311Location(ctx, fc)
	if err != nil {
		return err
	}

	return nil
}

func SubmitFHV(ctx context.Context, d db.ReadOnly, fc *complaint.FullComplaint) error {
	log.Printf("SubmitFHV")
	sel := `//button/div/div[text()='Driver Complaint']`
	if err := chromedp.Run(ctx,
		chromedp.ScrollIntoView(sel),
		chromedp.Sleep(time.Millisecond*100),
		chromedp.Click(sel),
		chromedp.Sleep(time.Millisecond*100),
	); err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 100)
	start := `//h5/a[text()[contains(.,"not a passenger")]]`
	log.Printf("click start %s", start)
	if resp, err := chromedp.RunResponse(ctx,
		chromedp.ScrollIntoView(start),
		chromedp.Sleep(time.Millisecond*10),
		chromedp.Click(start),
	); err != nil {
		return err
	} else {
		fmt.Println("RunResponse status code:", resp.Status)
	}

	time.Sleep(time.Second * 2)
	if err := checkStep(ctx, "What"); err != nil {
		return err
	}

	log.Printf("> filling fields")
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(`#n311_attendhearing_1`, chromedp.ByID),
		chromedp.Click(`#n311_attendhearing_1`, chromedp.ByID),
		chromedp.Click(`#n311_coloroftaxi_2`, chromedp.ByID, chromedp.NodeVisible),
		chromedp.ScrollIntoView(`n311_licensenumber`, chromedp.ByID),
		chromedp.SetValue(`#n311_licensenumber`, fc.License, chromedp.ByID),
		chromedp.SetValue(`#n311_datetimeobserved`, fc.Time.Format("2006-01-02T15:04:05.000Z"), chromedp.ByID),
		chromedp.SetValue(`#n311_datetimeobserved_datepicker_description`, fc.Time.Format("1/2/2006 3:04 PM"), chromedp.ByID),
	); err != nil {
		return err
	}

	// upload attachments; video first
	for _, file := range uploads(fc) {
		f := filepath.Join(d.FullPath(fc.Complaint), file)
		err := uploadFile(ctx, f, time.Second*90)
		if err != nil {
			return err
		}
	}

	// do this last so we can wait on a 'RunResponse'
	fmt.Println("waiting for Next")
	if resp, err := chromedp.RunResponse(ctx,
		chromedp.SetValue(`#n311_description`, fc.Description, chromedp.ByID),
	); err != nil {
		return err
	} else {
		fmt.Println("RunResponse status code:", resp.Status)
	}

	err := upload311Location(ctx, fc)
	if err != nil {
		return err
	}

	return nil
}

func uploads(fc *complaint.FullComplaint) []string {
	// TODO: check for file size
	uploads := fc.Videos
	for _, p := range fc.Photos {
		switch filepath.Ext(strings.ToLower(p)) {
		case ".heic":
		default:
			uploads = append(uploads, p)
		}
	}
	if len(uploads) > 3 {
		uploads = uploads[:3]
	}
	return uploads
}

func uploadFile(ctx context.Context, f string, timeout time.Duration) error {
	log.Printf("uploading %q", f)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := chromedp.Run(ctx,
		chromedp.ScrollIntoView(`#attachments-addbutton`, chromedp.ByID),
		chromedp.Click(`#attachments-addbutton`, chromedp.ByID),      // add attachment
		chromedp.WaitVisible(`input[type="file"]`, chromedp.ByQuery), // wait for modal
		chromedp.SetUploadFiles(`input[type="file"]`, []string{f}, chromedp.ByQuery),
		chromedp.Click(`//button[text() = 'Add Attachment']`), // upload
		// TODO: check for error uploading?
		chromedp.WaitVisible(fmt.Sprintf(`//p[contains(@class, 'attachmentTitle') and text()='%s']`, filepath.Base(f))), // the final results table
	); err != nil {
		return err
	}
	return nil
}

func checkStep(ctx context.Context, expected string) error {
	var step string
	err := chromedp.Run(ctx,
		chromedp.TextContent(".step.active", &step, chromedp.ByQueryAll),
	)
	step = strings.TrimSpace(step)
	fmt.Println("on step", step)
	if err != nil {
		return err
	}
	if step != expected {
		return fmt.Errorf("on step %q expected %q", step, expected)
	}
	return nil
}

func upload311Location(ctx context.Context, fc *complaint.FullComplaint) error {
	// https://portal.311.nyc.gov/sr-step/?id=4a3484e4-cd1e-ee11-a81c-6045bdb05de8&stepid=9241458f-fb0d-e811-8127-1458d04d2538
	// wait for title?

	time.Sleep(time.Second * 2) // TODO: wait for spinner to go away

	if err := checkStep(ctx, "Where"); err != nil {
		return err
	}

	fmt.Println("> filling location ", fc.Location)
	if resp, err := chromedp.RunResponse(ctx,
		chromedp.SetValue(`#n311_locationtypeid_select`, "a7c99a56-e64e-e811-a951-000d3a3606de", chromedp.ByID), // street
		chromedp.Sleep(time.Second),
		// Bug - the RunResponse is stalling here
		chromedp.SetValue(`#n311_additionallocationdetails`, fc.Location, chromedp.ByID),
		chromedp.Sleep(time.Millisecond*100),
	); err != nil {
		return err
	} else {
		fmt.Println("location RunResponse status code:", resp.Status)
	}

	if fc.Address != "" {
		fmt.Println("> filling address ", fc.Address)
		if resp, err := chromedp.RunResponse(ctx,
			chromedp.Sleep(time.Second),
			chromedp.Click(`#SelectAddressWhere`, chromedp.ByID),
			chromedp.Sleep(time.Second),
			chromedp.WaitVisible(`#address-search-box-input`, chromedp.ByID),
			chromedp.SetValue(`#address-search-box-input`, fc.Address, chromedp.ByID), // this is search pre-population
			chromedp.Sleep(time.Second),
			// wait for #ui-id-1 > li > div
			chromedp.WaitVisible(`#ui-id-2`, chromedp.ByID),
			chromedp.Click(`#ui-id-2`, chromedp.ByID, chromedp.NodeVisible),
			chromedp.Sleep(time.Millisecond*100),
			chromedp.ScrollIntoView(`#SelectAddressMap`, chromedp.ByID),
			chromedp.Click(`#SelectAddressMap`, chromedp.ByID, chromedp.NodeVisible),
		); err != nil {
			return err
		} else {
			fmt.Println("address RunResponse status code:", resp.Status)
		}

		// ui-id-1 > li > div
		// document.getElementById('SelectAddressMap').click()

	}
	// .address-picker-input ?
	return nil
}

func save311(ctx context.Context) {
	// #n311_name .value

}
