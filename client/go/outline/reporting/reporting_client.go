/**
 * Copyright 2025 The Outline Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package reporting

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Jigsaw-Code/outline-apps/client/go/outline/connectivity"
	"github.com/Jigsaw-Code/outline-sdk/transport"
)

const (
	testTCPWebsite = "example.com:443"
)

func CheckTCPConnectivity(tcp transport.StreamDialer) error {
	tcpErr := connectivity.CheckTCPConnectivityWithHTTP(tcp, testTCPWebsite)
	if tcpErr != nil {
		return tcpErr
	}
	return nil
}

// Start starts reporting.
func Report(tcp transport.StreamDialer) (err error) {
	// Perform a TCP connectivity check.
	if err := CheckTCPConnectivity(tcp); err != nil {
		return fmt.Errorf("TCP connectivity check failed: %w", err)
	}
	// Create a context with a timeout to avoid indefinite hangs.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Dial the stream using the StreamDialer.
	conn, err := tcp.DialStream(ctx, testTCPWebsite)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", testTCPWebsite, err)
	}
	defer conn.Close()

	// Write an HTTP request to the connection.
	request := "GET / HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Connection: close\r\n\r\n"

	_, err = conn.Write([]byte(request))
	if err != nil {
		return fmt.Errorf("failed to write request: %w", err)
	}

	// Read the response from the server.
	var response []byte
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read response: %w", err)
		}
		response = append(response, buffer[:n]...)
	}

	// Parse the HTTP response to extract cookies.
	headersEnd := bytes.Index(response, []byte("\r\n\r\n"))
	if headersEnd == -1 {
		return fmt.Errorf("failed to find end of headers in response")
	}

	headers := string(response[:headersEnd])
	fmt.Println("HTTP Headers:")
	fmt.Println(headers)

	// Extract cookies from the headers.
	lines := strings.Split(headers, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Set-Cookie:") {
			cookie := strings.TrimPrefix(line, "Set-Cookie: ")
			fmt.Println("Cookie found:", cookie)
		}
	}

	return nil
}

// Time duration constant
const reportInterval = 10 * time.Second

// StartReporting calls the Report function every 10 seconds
func StartReporting(tcp transport.StreamDialer) {
	ticker := time.NewTicker(reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := Report(tcp)
			if err != nil {
				// Handle error (e.g., log it)
				fmt.Printf("Report failed: %v\n", err)
			}
		}
	}
}
