// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jaegertracing/jaeger/model"
	jaegerM "github.com/uber/jaeger-lib/metrics"
)

func TestProcessorMetrics(t *testing.T) {
	baseMetrics := jaegerM.NewLocalFactory(time.Hour)
	serviceMetrics := baseMetrics.Namespace("service", nil)
	hostMetrics := baseMetrics.Namespace("host", nil)
	spm := NewSpanProcessorMetrics(serviceMetrics, hostMetrics, []string{"scruffy"})
	benderFormatMetrics := spm.GetCountsForFormat("bender")
	assert.NotNil(t, benderFormatMetrics)
	jFormat := spm.GetCountsForFormat(JaegerFormatType)
	assert.NotNil(t, jFormat)
	jFormat.ReceivedBySvc.ReportServiceNameForSpan(&model.Span{
		Process: &model.Process{},
	})
	mSpan := model.Span{
		Process: &model.Process{
			ServiceName: "fry",
		},
	}
	jFormat.ReceivedBySvc.ReportServiceNameForSpan(&mSpan)
	mSpan.Flags.SetDebug()
	mSpan.ParentSpanID = model.SpanID(1234)
	jFormat.ReceivedBySvc.ReportServiceNameForSpan(&mSpan)
	counters, gauges := baseMetrics.LocalBackend.Snapshot()

	assert.EqualValues(t, 2, counters["service.jaeger.spans.by-svc.fry"])
	assert.EqualValues(t, 1, counters["service.jaeger.traces.by-svc.fry"])
	assert.EqualValues(t, 1, counters["service.jaeger.debug-spans.by-svc.fry"])
	assert.Empty(t, gauges)
}
