package main

import (
	"context"
	"fmt"
	"log"
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
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
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
		chromedp.Flag("disable-extensions", true),
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
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var url string
	switch fc.VehicleType {
	case reg.FHV.String():
		url = "https://portal.311.nyc.gov/article/?kanumber=KA-01244"
	default:
		url = "https://portal.311.nyc.gov/article/?kanumber=KA-01241"
	}

	// navigate
	fmt.Printf("> opening %s\n", url)
	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		return err
	}

	fmt.Printf("10s")
	time.Sleep(10 * time.Second)
	return nil
}
