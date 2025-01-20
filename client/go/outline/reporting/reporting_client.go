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
	"github.com/Jigsaw-Code/outline-apps/client/go/outline/connectivity"
	"github.com/Jigsaw-Code/outline-sdk/transport"
)

const (
	testTCPWebsite = "http://example.com"
)

// Start starts reporting.
func Start(tcp transport.StreamDialer, udp transport.PacketListener) (err error) {
	// Perform an initial check.
	tcpErr := connectivity.CheckTCPConnectivityWithHTTP(tcp, testTCPWebsite)
	if tcpErr != nil {
		return tcpErr
	}
	return
}
