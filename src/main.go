package main

import (
	"bufio"
	"fmt"
	"go2web/src/utils"
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
			fmt.Println("Fetching URL:", url)
			fmt.Println(utils.SendHTTPRequest("GET", url))
		case "-s":
			if len(args) < 2 {
				fmt.Println("Usage: -s <search-term>")
				continue
			}
			searchTerm := strings.Join(args[1:], " ")
			fmt.Println("Searching for:", searchTerm)
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
