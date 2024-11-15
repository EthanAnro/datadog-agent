// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build kubeapiserver

package kubernetesapiserver

import (
	v1 "k8s.io/api/core/v1"

	tagger "github.com/DataDog/datadog-agent/comp/core/tagger/def"
	"github.com/DataDog/datadog-agent/pkg/metrics/event"
)

func newBundledTransformer(clusterName string, taggerInstance tagger.Component, collectedTypes []collectedEventType, filteringEnabled bool) eventTransformer {
	return &bundledTransformer{
		clusterName:      clusterName,
		taggerInstance:   taggerInstance,
		collectedTypes:   collectedTypes,
		filteringEnabled: filteringEnabled,
	}
}

type bundledTransformer struct {
	clusterName      string
	taggerInstance   tagger.Component
	collectedTypes   []collectedEventType
	filteringEnabled bool
}

func (c *bundledTransformer) Transform(events []*v1.Event) ([]event.Event, []error) {
	var errors []error

	bundlesByObject := make(map[bundleID]*kubernetesEventBundle)

	for _, event := range events {
		if event.InvolvedObject.Kind == "" ||
			event.InvolvedObject.Name == "" ||
			event.Reason == "" ||
			event.Message == "" {
			continue
		}

		kubeEvents.Inc(
			event.InvolvedObject.Kind,
			event.Source.Component,
			event.Type,
			event.Reason,
			getEventSource(event.ReportingController, event.Source.Component),
		)

		if c.filteringEnabled {
			if !(shouldCollectByDefault(event) || shouldCollect(event, c.collectedTypes)) {
				continue
			}
		}

		id := buildBundleID(event)

		bundle, found := bundlesByObject[id]
		if !found {
			bundle = newKubernetesEventBundler(c.clusterName, event)
			bundlesByObject[id] = bundle
		}

		err := bundle.addEvent(event)
		if err != nil {
			errors = append(errors, err)
			continue
		}
	}

	datadogEvs := make([]event.Event, 0, len(bundlesByObject))

	for id, bundle := range bundlesByObject {
		datadogEv, err := bundle.formatEvents(c.taggerInstance)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		emittedEvents.Inc(
			id.kind,
			id.evType,
			getEventSource(bundle.reportingController, bundle.component),
			"true",
		)

		datadogEvs = append(datadogEvs, datadogEv)
	}

	return datadogEvs, errors
}

type bundleID struct {
	kind   string
	uid    string
	evType string
}

// buildBundleID generates a unique ID to separate k8s events
// based on their InvolvedObject UIDs and event Types
func buildBundleID(e *v1.Event) bundleID {
	return bundleID{
		kind:   e.InvolvedObject.Kind,
		uid:    string(e.InvolvedObject.UID),
		evType: e.Type,
	}
}
