/**
 * Copyright (C) 2011 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
	mesh "istio.io/api/mesh/v1alpha1"
	mixer "istio.io/api/mixer/v1"
	routing "istio.io/api/routing/v1alpha1"
	prometheus "istio.io/istio/mixer/adapter/prometheus/config"
	"github.com/snowdrop/istio-java-api/pkg/schemagen"
	"istio.io/istio/mixer/template/metric"
	"istio.io/istio/mixer/template/checknothing"
	"istio.io/istio/mixer/template/listentry"
	"istio.io/istio/mixer/template/logentry"
	"istio.io/istio/mixer/template/apikey"
	"istio.io/istio/mixer/template/authorization"
	"istio.io/istio/mixer/template/quota"
	"istio.io/istio/mixer/template/reportnothing"
	"istio.io/istio/mixer/template/tracespan"
	circonus "istio.io/istio/mixer/adapter/circonus/config"
	denier "istio.io/istio/mixer/adapter/denier/config"
	rbac "istio.io/api/rbac/v1alpha1"
)

type Schema struct {
	MeshConfig           mesh.MeshConfig
	ProxyConfig          mesh.ProxyConfig
	Attributes           mixer.Attributes
	AttributeValue       mixer.Attributes_AttributeValue
	CheckRequest         mixer.CheckRequest
	QuotaParams          mixer.CheckRequest_QuotaParams
	CheckResponse        mixer.CheckResponse
	QuotaResult          mixer.CheckResponse_QuotaResult
	CompressedAttributes mixer.CompressedAttributes
	ReferencedAttributes mixer.ReferencedAttributes
	ReportRequest        mixer.ReportRequest
	ReportResponse       mixer.ReportResponse
	StringMap            mixer.StringMap
	CircuitBreaker       routing.CircuitBreaker
	CorsPolicy           routing.CorsPolicy
	DestinationPolicy    routing.DestinationPolicy
	DestinationWeight    routing.DestinationWeight
	EgressRule           routing.EgressRule
	HTTPFaultInjection   routing.HTTPFaultInjection
	HTTPRedirect         routing.HTTPRedirect
	HTTPRetry            routing.HTTPRetry
	HTTPRewrite          routing.HTTPRewrite
	HTTPTimeout          routing.HTTPTimeout
	IngressRule          routing.IngressRule
	IstioService         routing.IstioService
	L4FaultInjection     routing.L4FaultInjection
	L4MatchAttributes    routing.L4MatchAttributes
	LoadBalancing        routing.LoadBalancing
	MatchCondition       routing.MatchCondition
	MatchRequest         routing.MatchRequest
	RouteRule            routing.RouteRule
	StringMatch          routing.StringMatch
	ServiceRole          rbac.ServiceRole
	ServiceRoleBinding   rbac.ServiceRoleBinding
	Prometheus           prometheus.Params
	Circonus             circonus.Params
	Denier               denier.Params
	APIKey               apikey.InstanceMsg
	Authorization        authorization.InstanceMsg
	CheckNothing         checknothing.InstanceMsg
	ListEntry            listentry.InstanceMsg
	LogEntry             logentry.InstanceMsg
	Metric               metric.InstanceMsg
	Quota                quota.InstanceMsg
	ReportNothing        reportnothing.InstanceMsg
	TraceSpan            tracespan.InstanceMsg
}

func main() {
	packages := []schemagen.PackageDescriptor{
		{"istio.io/api/mesh/v1alpha1", "me.snowdrop.istio.api.model.v1.mesh", "istio_mesh_"},
		{"istio.io/api/mixer/v1", "me.snowdrop.istio.api.model.v1.mixer", "istio_mixer_"},
		{"istio.io/api/routing/v1alpha1", "me.snowdrop.istio.api.model.v1.routing", "istio_routing_"},
		{"istio.io/api/rbac/v1alpha1", "me.snowdrop.istio.api.model.v1.rbac", "istio_rbac_"},
		{"istio.io/istio/mixer/adapter/circonus/config", "me.snowdrop.istio.adapter.circonus", "istio_adapter_circonus_"},
		{"istio.io/istio/mixer/adapter/denier/config", "me.snowdrop.istio.adapter.denier", "istio_adapter_denier_"},
		{"istio.io/istio/mixer/adapter/prometheus/config", "me.snowdrop.istio.adapter.prometheus", "istio_adapter_prometheus_"},
		{"istio.io/istio/mixer/template/apikey", "me.snowdrop.istio.api.model.v1.mixer.template", "istio_mixer_apikey_"},
		{"istio.io/istio/mixer/template/authorization", "me.snowdrop.istio.api.model.v1.mixer.template", "istio_mixer_authorization_"},
		{"istio.io/istio/mixer/template/checknothing", "me.snowdrop.istio.api.model.v1.mixer.template", "istio_mixer_checknothing_"},
		{"istio.io/istio/mixer/template/listentry", "me.snowdrop.istio.api.model.v1.mixer.template", "istio_mixer_listentry_"},
		{"istio.io/istio/mixer/template/logentry", "me.snowdrop.istio.api.model.v1.mixer.template", "istio_mixer_logentry_"},
		{"istio.io/istio/mixer/template/metric", "me.snowdrop.istio.api.model.v1.mixer.template", "istio_mixer_metric_"},
		{"istio.io/istio/mixer/template/quota", "me.snowdrop.istio.api.model.v1.mixer.template", "istio_mixer_quota_"},
		{"istio.io/istio/mixer/template/reportnothing", "me.snowdrop.istio.api.model.v1.mixer.template", "istio_mixer_reportnothing_"},
		{"istio.io/istio/mixer/template/tracespan", "me.snowdrop.istio.api.model.v1.mixer.template", "istio_mixer_tracespan_"},
		{"github.com/golang/protobuf/ptypes/duration", "me.snowdrop.istio.api.model", "protobuf_duration_"},
		{"github.com/gogo/protobuf/types", "me.snowdrop.istio.api.model", "protobuf_types_"},
		{"github.com/golang/protobuf/ptypes/any", "me.snowdrop.istio.api.model", "protobuf_any_"},
		{"github.com/gogo/googleapis/google/rpc", "me.snowdrop.istio.api.model", "google_rpc_"},
	}

	typeMap := map[reflect.Type]reflect.Type{
		reflect.TypeOf(time.Time{}): reflect.TypeOf(""),
		reflect.TypeOf(struct{}{}):  reflect.TypeOf(""),
	}

	enumMap := map[string]string{
		"istio.mesh.v1alpha1.MeshConfig_IngressControllerMode":      "me.snowdrop.istio.api.model.v1.mesh.IngressControllerMode",
		"istio.mesh.v1alpha1.MeshConfig_AuthPolicy":                 "me.snowdrop.istio.api.model.v1.mesh.AuthenticationPolicy",
		"istio.mesh.v1alpha1.AuthenticationPolicy":                  "me.snowdrop.istio.api.model.v1.mesh.AuthenticationPolicy",
		"istio.mesh.v1alpha1.ProxyConfig_InboundInterceptionMode":   "me.snowdrop.istio.api.model.v1.mesh.InboundInterceptionMode",
		"istio.mesh.v1alpha1.MeshConfig_OutboundTrafficPolicy_Mode": "me.snowdrop.istio.api.model.v1.mesh.Mode",
		"istio.mixer.v1.ReferencedAttributes_Condition":             "me.snowdrop.istio.api.model.v1.mixer.Condition",
		"istio.mixer.v1.config.descriptor.ValueType":                "me.snowdrop.istio.api.model.v1.mixer.config.descriptor.ValueType",
		"adapter.circonus.config.Params_MetricInfo_Type":            "me.snowdrop.istio.adapter.circonus.Type",
		"adapter.prometheus.config.Params_MetricInfo_Kind":          "me.snowdrop.istio.adapter.prometheus.Kind",
	}

	schema, err := schemagen.GenerateSchema(reflect.TypeOf(Schema{}), packages, typeMap, enumMap)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	args := os.Args[1:]
	if len(args) < 1 || args[0] != "validation" {
		schema.Resources = nil
	}

	b, err := json.Marshal(&schema)
	if err != nil {
		log.Fatal(err)
	}
	result := string(b)
	result = strings.Replace(result, "\"additionalProperty\":", "\"additionalProperties\":", -1)
	var out bytes.Buffer
	err = json.Indent(&out, []byte(result), "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out.String())
}
