// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package launchgui implements 'agent launch-gui'.
package launchgui

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/DataDog/datadog-agent/cmd/agent/command"
	"github.com/DataDog/datadog-agent/comp/core"
	"github.com/DataDog/datadog-agent/comp/core/config"
	log "github.com/DataDog/datadog-agent/comp/core/log/def"
	apiutil "github.com/DataDog/datadog-agent/pkg/api/util"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

// cliParams are the command-line arguments for this subcommand
type cliParams struct {
	*command.GlobalParams
}

// Commands returns a slice of subcommands for the 'agent' command.
func Commands(globalParams *command.GlobalParams) []*cobra.Command {
	cliParams := &cliParams{
		GlobalParams: globalParams,
	}
	launchCmd := &cobra.Command{
		Use:   "launch-gui",
		Short: "starts the Datadog Agent GUI",
		Long:  ``,
		RunE: func(_ *cobra.Command, _ []string) error {
			return fxutil.OneShot(launchGui,
				fx.Supply(cliParams),
				fx.Supply(command.GetDefaultCoreBundleParams(cliParams.GlobalParams)),
				core.Bundle(),
			)
		},
		SilenceUsage: true,
	}

	return []*cobra.Command{launchCmd}
}

func launchGui(config config.Component, _ *cliParams, _ log.Component) error {
	guiPort := config.GetString("GUI_port")
	if guiPort == "-1" {
		return fmt.Errorf("GUI not enabled: to enable, please set an appropriate port in your datadog.yaml file")
	}

	endpoint, err := apiutil.NewIPCEndpoint(config, "/agent/gui/intent")
	if err != nil {
		return err
	}

	intentToken, err := endpoint.DoGet()
	if err != nil {
		return err
	}

	// Open the GUI in a browser, passing the authorization tokens as parameters
	err = open("http://127.0.0.1:" + guiPort + "/auth?intent=" + string(intentToken))
	if err != nil {
		return fmt.Errorf("error opening GUI: %s", err.Error())
	}

	fmt.Printf("GUI opened at 127.0.0.1:%s\n", guiPort)
	return nil
}
