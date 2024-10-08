// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build windows

// Package service provides a way to interact with os services
package service

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/DataDog/datadog-agent/pkg/fleet/installer/repository"
	"github.com/DataDog/datadog-agent/pkg/fleet/internal/cdn"
	"github.com/DataDog/datadog-agent/pkg/fleet/internal/paths"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func msiexec(target, operation string, args []string) (err error) {
	updaterPath := filepath.Join(paths.PackagesPath, "datadog-agent", target)
	msis, err := filepath.Glob(filepath.Join(updaterPath, "datadog-agent-*-1-x86_64.msi"))
	if err != nil {
		return err
	}
	if len(msis) > 1 {
		return fmt.Errorf("too many MSIs in package")
	} else if len(msis) == 0 {
		return fmt.Errorf("no MSIs in package")
	}

	cmd := exec.Command("msiexec", append([]string{operation, msis[0], "/qn", "MSIFASTINSTALL=7"}, args...)...)
	return cmd.Run()
}

// SetupAgent installs and starts the agent
func SetupAgent(ctx context.Context, args []string) (err error) {
	span, _ := tracer.StartSpanFromContext(ctx, "setup_agent")
	defer func() {
		if err != nil {
			log.Errorf("Failed to setup agent: %s", err)
		}
		span.Finish(tracer.WithError(err))
	}()
	return msiexec("stable", "/i", args)
}

// StartAgentExperiment starts the agent experiment
func StartAgentExperiment(ctx context.Context) (err error) {
	span, _ := tracer.StartSpanFromContext(ctx, "start_experiment")
	defer func() {
		if err != nil {
			log.Errorf("Failed to start agent experiment: %s", err)
		}
		span.Finish(tracer.WithError(err))
	}()
	return msiexec("experiment", "/i", nil)
}

// StopAgentExperiment stops the agent experiment
func StopAgentExperiment(ctx context.Context) (err error) {
	span, _ := tracer.StartSpanFromContext(ctx, "stop_experiment")
	defer func() {
		if err != nil {
			log.Errorf("Failed to stop agent experiment: %s", err)
		}
		span.Finish(tracer.WithError(err))
	}()
	err = msiexec("experiment", "/x", nil)
	if err != nil {
		return err
	}

	// TODO: Need args here to restore DDAGENTUSER
	return msiexec("stable", "/i", nil)
}

// PromoteAgentExperiment promotes the agent experiment
func PromoteAgentExperiment(_ context.Context) error {
	// noop
	return nil
}

// RemoveAgent stops and removes the agent
func RemoveAgent(ctx context.Context) (err error) {
	span, _ := tracer.StartSpanFromContext(ctx, "remove_agent")
	defer func() {
		if err != nil {
			log.Errorf("Failed to remove agent: %s", err)
		}
		span.Finish(tracer.WithError(err))
	}()
	return msiexec("stable", "/x", nil)
}

// ConfigureAgent noop
func ConfigureAgent(_ context.Context, _ *cdn.CDN, _ *repository.Repositories) error {
	return nil
}
