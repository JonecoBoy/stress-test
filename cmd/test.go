/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/JonecoBoy/stress-test/appContext"
	"github.com/JonecoBoy/stress-test/reporter"
	"github.com/JonecoBoy/stress-test/requester"
	"time"

	"github.com/spf13/cobra"
)

var reportFormat string
var url string
var requestNumber uint
var concurrentRequests uint
var quiet bool
var addTimeStamp bool
var outputPath string
var outputType string

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "start streess test",
	Long: `You can use this command to start the streess test.
			Flags:
			-u --url: set target URL
			-r --requests: set number of requests
			-c --concurrency: set number of concurrent requests
			-q --quiet: dont log requests statuses
			Example:
			stresser test --url http://localhost:8080/api/v1/test --requests 100 --concurrency 10
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := &appContext.Context{
			URL:           url,
			TotalRequests: int(requestNumber),
			Concurrency:   int(concurrentRequests),
			Quiet:         bool(quiet),
		}

		fmt.Println("Starting stress test...")
		fmt.Printf("URL: %s, Requests: %d, Concurrency: %d\n", ctx.URL, ctx.TotalRequests, ctx.Concurrency)

		req := requester.NewRequester(ctx)
		rep := reporter.NewReporter(ctx)

		start := time.Now()
		req.Start(ctx)
		ctx.TotalTime = time.Since(start)

		report, err := rep.PrepareReport(outputType)
		if err != nil {
			fmt.Println("Error preparing report:", err)
		}
		rep.CliReport(report)

		// Remove the extension from the filename and add .html

		if outputPath == "" {
			outputPath = "report-"
		}
		if addTimeStamp {
			outputPath += time.Now().Format("2006-01-02-15-04-05")
		}

		err = rep.GenerateHTMLReport(outputPath)
		if err != nil {
			fmt.Println("Error generating HTML report:", err)
		}

		err = rep.LogToFile(report, outputPath, outputType)
		if err != nil {
			fmt.Println("Error writing report to file:", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	//testCmd.AddCommand(reporterCmd)

	testCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "set target URL (Required)")
	testCmd.PersistentFlags().UintVarP(&requestNumber, "requests", "r", 100, "set number of total requests")
	testCmd.PersistentFlags().UintVarP(&concurrentRequests, "concurrency", "c", 10, "set number of concurrent requests")
	testCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", "", "set output file for the report")
	testCmd.PersistentFlags().BoolVarP(&addTimeStamp, "timestamp", "t", true, "add timestamp to the output file name")
	testCmd.PersistentFlags().StringVarP(&outputType, "encode", "e", "", "set output format for the report (csv, json, toml, yaml)")
	testCmd.PersistentFlags().StringVarP(&reportFormat, "format", "f", "stdout", "set report output format (stdout,txt or html)")
	testCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "set quiet mode (won't print all response statuses)")

	err := testCmd.MarkPersistentFlagRequired("url")
	if err != nil {
		return
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
