package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	return Submit(fc)
}

func Submit(fc *complaint.FullComplaint) error {
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
		// load-extension=/path/to/extension from https://stackoverflow.com/questions/66970005/c-sharp-selenium-enable-default-extensions
		chromedp.Flag("load-extensions", filepath.Join(hd, "/Library/Application Support/Google/Chrome/Default/Extensions")),
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

	var url, title string
	switch fc.VehicleType {
	case reg.FHV.String():
		url = "https://portal.311.nyc.gov/article/?kanumber=KA-01244"
		title = "Car Service Complaint"
	default:
		url = "https://portal.311.nyc.gov/article/?kanumber=KA-01241"
		title = "Taxi Complaint"
	}

	// navigate
	fmt.Printf("> opening %s\n", url)
	var loginLink string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
	); err != nil {
		return err
	}

	log.Printf("waiting for page %q", title)
	if err := chromedp.Run(ctx,
		// chromedp.WaitVisible(fmt.Sprintf(`//h1[text()=%q]]`, title)),
		chromedp.TextContent(`.login-area > a`, &loginLink, chromedp.ByQuery),
	); err != nil {
		return err
	}
	log.Printf("login link %#v", loginLink)

	if loginLink == "Sign In" {
		fmt.Printf("> Sign In\n")
		if err := chromedp.Run(ctx,
			chromedp.Click(`.login-area > a`, chromedp.NodeVisible, chromedp.ByQuery),
		); err != nil {
			return err
		}

		var currentAddress string
		for currentAddress != url {
			time.Sleep(500 * time.Millisecond)
			if err := chromedp.Run(ctx,
				chromedp.Location(&currentAddress),
			); err != nil {
				return err
			}
		}
		log.Printf("checking title")
		if err := chromedp.Run(ctx); // chromedp.WaitVisible(fmt.Sprintf(`//h1[text()=%q]]`, title)),
		err != nil {
			return err
		}

		fmt.Printf("> authenticated")
	}
	switch fc.VehicleType {
	case reg.FHV.String():
		err = SubmitFHV(ctx, fc)
	default:
		err = SubmitTaxi(ctx, fc)
	}
	if err != nil {
		return err
	}
	input.Ask("Done?", "")
	return chromedp.Cancel(ctx)
}

func SubmitTaxi(ctx context.Context, fc *complaint.FullComplaint) error {
	// click = "You were not a passenger"
	return fmt.Errorf("not implemented")
}

func SubmitFHV(ctx context.Context, fc *complaint.FullComplaint) error {
	log.Printf("SubmitFHV")
	start := `//h5/a[text()[contains(.,"you were NOT a passenger")]]`
	if err := chromedp.Run(ctx,
		chromedp.Click(`//button/div/div[text()='Driver Complaint']`),
	); err != nil {
		return err
	}
	log.Printf("click start %s", start)
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(start),
		chromedp.Click(start),
	); err != nil {
		return err
	}

	// chromedp.Nodes(`document`, &nodes, chromedp.ByJSPath),
	// chromedp.Click(`#example-After`, chromedp.NodeVisible),
	// chromedp.Value(`#example-After textarea`, &example),
	// chromedp.TextContent
	return nil
}
