package main

import (
	"flag"
	"fmt"
	"go2web/src/utils"
	"net/url"
	"os"
)

var (
	urlFlag    = flag.String("u", "", "URL to fetch")
	searchTerm = flag.String("s", "", "Search term for Google search")
	showHelp   = flag.Bool("h", false, "Show help message")
)

func main() {
	flag.Parse()

	if *showHelp || len(os.Args) <= 1 {
		fmt.Println("Usage of go2web:")
		flag.PrintDefaults()
		return
	}

	if *urlFlag != "" {
		result, err := utils.SendHTTPRequest("GET", *urlFlag, 0)
		if err != nil {
			fmt.Println("Error fetching URL:", err)
			return
		}
		fmt.Println(result)
	} else if *searchTerm != "" {
		fmt.Println("Search functionality not implemented.")
		url := "https://google.com/search?q=" + url.QueryEscape(*searchTerm)
		fmt.Println("Search URL would be:", url)
	} else {
		fmt.Println("No operation specified. Use -h for help.")
	}
}
