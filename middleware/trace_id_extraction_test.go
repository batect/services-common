// Copyright 2019-2021 Charles Korn.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/batect/service-observability/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Based on https://cloud.google.com/run/docs/logging#writing_structured_logs and
// https://cloud.google.com/trace/docs/troubleshooting#force-trace
var _ = Describe("Trace ID extraction middleware", func() {
	var ctx context.Context
	var m http.Handler

	BeforeEach(func() {
		m = middleware.TraceIDExtractionMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			ctx = r.Context()
		}))
	})

	Context("when the request contains an established trace", func() {
		var traceID string

		BeforeEach(func() {
			req := httptest.NewRequest("GET", "/blah", nil)
			req, traceID = addTraceToRequest(req)
			m.ServeHTTP(nil, req)
		})

		It("extracts the trace ID from the established trace", func() {
			Expect(middleware.TraceIDFromContext(ctx)).To(Equal(traceID))
		})
	})

	Context("when the request does not contain an established trace", func() {
		BeforeEach(func() {
			req := httptest.NewRequest("GET", "/blah", nil)
			m.ServeHTTP(nil, req)
		})

		It("returns a generated trace ID", func() {
			Expect(middleware.TraceIDFromContext(ctx)).To(HavePrefix("autogenerated-"))
		})
	})
})

func addTraceToRequest(req *http.Request) (*http.Request, string) {
	tracer := sdktrace.NewTracerProvider().Tracer("Tracer")
	ctx, span := tracer.Start(context.Background(), "My test span")
	traceID := span.SpanContext().TraceID().String()
	req = req.WithContext(ctx)

	return req, traceID
}
