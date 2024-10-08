// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build serverless

package agentimpl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/DataDog/datadog-agent/comp/core/config"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

func TestBuildServerlessEndpoints(t *testing.T) {
	config := fxutil.Test[config.Component](t, fx.Options(
		config.MockModule(),
	))

	endpoints, err := buildEndpoints()
	assert.Nil(t, err)
	assert.Equal(t, "http-intake.logs.datadoghq.com", endpoints.Main.Host)
	assert.Equal(t, "lambda-extension", string(endpoints.Main.Origin))
	assert.True(t, endpoints.Main.BatchWait > config.BatchWait*time.Second)
}
