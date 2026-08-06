// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package vsock

import (
	"bufio"
	"context"
	"errors"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemporaryNetErr(t *testing.T) {
	assert.True(t, isTemporaryNetErr(&ackError{cause: errors.New("ack")}))

	assert.False(t, isTemporaryNetErr(&connectMsgError{cause: errors.New("connect")}))
	assert.False(t, isTemporaryNetErr(errors.New("something else")))
	assert.False(t, isTemporaryNetErr(nil))
}

// TestDialFirstAttemptImmediate verifies dial connects on the first attempt
// without waiting for a RetryInterval tick. The interval is set far larger
// than the test deadline, so a ticker-gated first attempt would time out.
func TestDialFirstAttemptImmediate(t *testing.T) {
	udsPath := filepath.Join(t.TempDir(), "fc.sock")
	ln, err := net.Listen("unix", udsPath)
	require.NoError(t, err)
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		// Drain the "CONNECT <port>\n" message and reply with an OK ack.
		_, _ = bufio.NewReaderSize(conn, 32).ReadString('\n')
		_, _ = conn.Write([]byte("OK 12345\n"))
		// Keep the connection open briefly so the caller's read succeeds.
		time.Sleep(50 * time.Millisecond)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	c := defaultConfig()
	c.RetryInterval = 30 * time.Second // would dominate if the first attempt waited for a tick

	start := time.Now()
	conn, err := dial(ctx, udsPath, 2049, c)
	require.NoError(t, err)
	defer conn.Close()

	assert.Less(t, time.Since(start), 5*time.Second,
		"dial should attempt immediately, not wait for the first RetryInterval tick")
}
