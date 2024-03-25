package utils

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func SendHTTPRequest(method, rawURL string, redirectCount int8) (body string, err error) {
	if val, ok, _ := Get(rawURL); ok {
		return val, nil
	}
	if redirectCount == 5 {
		return "Cannot Parse : " + rawURL, nil
	}

	currentRedirectCount := redirectCount

	headers := ""
	rawURL = strings.ReplaceAll(rawURL, "\r", "")
	if !hasScheme(rawURL) {
		rawURL = "http://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)

	if err != nil {
		return "", fmt.Errorf("failed parsing URL: %s", err)
	}

	host := parsedURL.Host
	path := parsedURL.Path

	if path == "" {
		path = "/"
	}

	var conn net.Conn
	if strings.Contains(parsedURL.Scheme, "https") {
		if !strings.Contains(host, ":") {
			host += ":443"
		}
		conn, err = tls.Dial("tcp", host, &tls.Config{})
	} else {
		if !strings.Contains(host, ":") {
			host += ":80"
		}
		conn, err = net.Dial("tcp", host)
	}

	if err != nil {
		return "", fmt.Errorf("failed connecting to given host: %s", err)
	}
	defer conn.Close()

	request := fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\nUser-Agent: UTM\r\n\r\n", method, path, parsedURL.Host)
	_, err = conn.Write([]byte(request))
	if err != nil {
		return "", fmt.Errorf("failed requesting content: %s", err)
	}

	reader := bufio.NewReader(conn)
	var response strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if line == "\r\n" {
			break
		}
		response.WriteString(line)
	}

	headers = response.String()

	statusCode, err := getStatusCode(headers)
	contentType := getContentType(headers)

	if err != nil {
		return "", fmt.Errorf("failed parsing status code: %s", err)
	}

	switch statusCode {
	case 301:
		fmt.Println("statusCode: ", statusCode)
		for _, line := range strings.Split(headers, "\n") {
			if strings.HasPrefix(line, "Location:") {
				locationURL := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
				urlParsed, err := url.Parse(locationURL)
				if err != nil {
					return "", err
				}
				return SendHTTPRequest("GET", urlParsed.String(), currentRedirectCount+1)
			}
		}
	}

	response.Reset()

	for {
		line, err := reader.ReadString('\n')
		if strings.Contains(contentType, "application/json") {
			body = line
			break
		}
		if err != nil {
			break
		}
		response.WriteString(line)
	}

	if strings.Contains(contentType, "text/html") {
		body = response.String()
	}

	var parsedContent string

	if strings.Contains(contentType, "text/html") {
		parsedContent = ParseHTML(body)
	} else if strings.Contains(contentType, "application/json") {
		parsedContent = string(ParseJSON(body))
	}

	Set(rawURL, parsedContent)
	return parsedContent, nil
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

func getContentType(headers string) string {
	parts := strings.Split(headers, "\r\n")

	for _, header := range parts {
		if strings.Contains(header, "Content-Type") {
			contentTypeLine := strings.Split(header, ": ")[1]
			return contentTypeLine
		}

	}
	return ""
}

func hasScheme(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}
