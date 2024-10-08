// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build !windows

package listener

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DataDog/datadog-agent/comp/logs/agent/config"
	"github.com/DataDog/datadog-agent/pkg/logs/message"
	"github.com/DataDog/datadog-agent/pkg/logs/pipeline/mock"
	"github.com/DataDog/datadog-agent/pkg/logs/sources"
)

func TestUDPShoulProperlyCollectLogSplitPerDatadgram(t *testing.T) {
	pp := mock.NewMockProvider()
	msgChan := pp.NextPipelineChan()
	frameSize := 100
	listener := NewUDPListener(pp, sources.NewLogSource("", &config.LogsConfig{Port: udpTestPort}), frameSize)
	listener.Start()

	conn, err := net.Dial("udp", listener.tailer.Conn.LocalAddr().String())
	assert.Nil(t, err)

	var msg *message.Message

	fmt.Fprint(conn, strings.Repeat("a", 10))
	msg = <-msgChan
	assert.Equal(t, strings.Repeat("a", 10), string(msg.GetContent()))

	fmt.Fprint(conn, strings.Repeat("a", 10)+"\n")
	msg = <-msgChan
	assert.Equal(t, strings.Repeat("a", 10), string(msg.GetContent()))

	fmt.Fprint(conn, strings.Repeat("a", 10)+"\n"+strings.Repeat("a", 10))
	msg = <-msgChan
	assert.Equal(t, strings.Repeat("a", 10), string(msg.GetContent()))
	msg = <-msgChan
	assert.Equal(t, strings.Repeat("a", 10), string(msg.GetContent()))

	fmt.Fprint(conn, strings.Repeat("a", 10)+"\n"+strings.Repeat("a", 10)+"\n")
	msg = <-msgChan
	assert.Equal(t, strings.Repeat("a", 10), string(msg.GetContent()))
	msg = <-msgChan
	assert.Equal(t, strings.Repeat("a", 10), string(msg.GetContent()))

	listener.Stop()
}

func TestUDPShouldProperlyTruncateBigMessages(t *testing.T) {
	pp := mock.NewMockProvider()
	msgChan := pp.NextPipelineChan()
	frameSize := 100
	listener := NewUDPListener(pp, sources.NewLogSource("", &config.LogsConfig{Port: udpTestPort}), frameSize)
	listener.Start()

	conn, err := net.Dial("udp", listener.tailer.Conn.LocalAddr().String())
	assert.Nil(t, err)

	var msg *message.Message

	fmt.Fprint(conn, strings.Repeat("a", frameSize-10))
	msg = <-msgChan
	assert.Equal(t, strings.Repeat("a", frameSize-10), string(msg.GetContent()))

	fmt.Fprint(conn, strings.Repeat("a", frameSize))
	msg = <-msgChan
	assert.Equal(t, strings.Repeat("a", frameSize), string(msg.GetContent()))

	fmt.Fprint(conn, strings.Repeat("a", frameSize+10))
	msg = <-msgChan
	assert.Equal(t, strings.Repeat("a", frameSize), string(msg.GetContent()))

	listener.Stop()
}

func TestUDPShoulDropTooBigMessages(t *testing.T) {
	// Skip if we can't detect the max UDP frame length (currently only linux/darwin).
	if maxUDPFrameLen <= 0 {
		return
	}

	pp := mock.NewMockProvider()
	msgChan := pp.NextPipelineChan()
	listener := NewUDPListener(pp, sources.NewLogSource("", &config.LogsConfig{Port: udpTestPort}), maxUDPFrameLen)
	listener.Start()

	conn, err := net.Dial("udp", listener.tailer.Conn.LocalAddr().String())
	assert.Nil(t, err)

	var msg *message.Message

	fmt.Fprint(conn, strings.Repeat("a", maxUDPFrameLen-100)+"\n")
	msg = <-msgChan
	assert.Equal(t, strings.Repeat("a", maxUDPFrameLen-100), string(msg.GetContent()))

	// the first frame should be dropped as it's too big compare to the limit.
	fmt.Fprint(conn, strings.Repeat("a", maxUDPFrameLen+100)+"\n")
	fmt.Fprint(conn, strings.Repeat("a", maxUDPFrameLen-200)+"\n")
	msg = <-msgChan
	assert.Equal(t, strings.Repeat("a", maxUDPFrameLen-200), string(msg.GetContent()))

	listener.Stop()
}
