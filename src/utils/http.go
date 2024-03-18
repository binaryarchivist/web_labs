package utils

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func SendHTTPRequest(method, rawURL string) (headers string, body string, err error) {
	fmt.Println("*rawURL: ", rawURL)

	rawURL = strings.ReplaceAll(rawURL, "\r", "")
	if !hasScheme(rawURL) {
		rawURL = "http://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)

	if err != nil {
		return "", "", fmt.Errorf("failed parsing URL: %s", err)
	}

	host := parsedURL.Host
	path := parsedURL.Path

	if path == "" {
		path = "/"
	}

	if !strings.Contains(host, ":") {
		host += ":80"
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		return "", "", fmt.Errorf("failed connecting to given host: %s", err)
	}
	defer conn.Close()

	request := fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", method, path, parsedURL.Host)
	_, err = conn.Write([]byte(request))
	if err != nil {
		return "", "", fmt.Errorf("failed requesting content: %s", err)
	}

	reader := bufio.NewReader(conn)
	var response strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", "", fmt.Errorf("failed parsing headers: %s", err)
		}

		// Headers are terminated by a blank line (\r\n).
		if line == "\r\n" {
			break
		}
		response.WriteString(line)
	}

	headers = response.String()

	statusCodeLine := strings.Split(headers, "\n")[0]
	statusCode, err := getStatusCode(statusCodeLine)

	if err != nil {
		return "", "", fmt.Errorf("failed parsing status code: %s", err)
	}

	fmt.Println("statusCode: ", statusCode)
	switch statusCode {
	case 301:
		redirectedHostLine := strings.Split(headers, "\n")[1]
		hostname := strings.Split(redirectedHostLine, ": ")[1]
		// fmt.Println("Redirected to hostname: ", hostname)
		return SendHTTPRequest("GET", hostname)

	}

	response.Reset()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		response.WriteString(line)
	}

	body = response.String()

	return headers, body, nil
}

func getStatusCode(statusLine string) (int, error) {
	parts := strings.Split(statusLine, " ")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid status line")
	}
	statusCode, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("failed parsing status code to int: %s", err)
	}
	return statusCode, nil
}

func hasScheme(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}
