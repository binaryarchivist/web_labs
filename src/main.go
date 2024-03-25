package main

import (
	"bufio"
	"fmt"
	"go2web/src/utils"
	"net/url"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("go2web> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		args := strings.Split(input, " ")

		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "-u":
			if len(args) != 2 {
				fmt.Println("Usage: -u <URL>")
				continue
			}
			url := args[1]
			result, _ := utils.SendHTTPRequest("GET", url, 0)
			fmt.Println(result)

		case "-s":
			if len(args) < 2 {
				fmt.Println("Usage: -s <search-term>")
				continue
			}
			searchTerm := strings.Join(args[1:], " ")

			url, _ := url.Parse("https://google.com/")
			query := url.Query()
			query.Add("search", searchTerm)
			url.RawQuery = query.Encode()

			// result, _ := utils.GetLinksFromGoogle(url)
			// fmt.Println(result)

		case "-h":
			fmt.Println("Help message...")
		case "exit":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Unknown command:", args[0])
		}
	}
}
