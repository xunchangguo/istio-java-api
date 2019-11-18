/*
 * *
 *  * Copyright (C) 2018 Red Hat, Inc.
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *         http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 *
 */
package me.snowdrop.istio.template.metric;

import java.io.InputStream;
import java.util.Map;

import io.fabric8.kubernetes.api.model.HasMetadata;
import me.snowdrop.istio.api.cexl.TypedValue;
import me.snowdrop.istio.api.policy.v1beta1.Instance;
import me.snowdrop.istio.api.policy.v1beta1.InstanceBuilder;
import me.snowdrop.istio.api.policy.v1beta1.InstanceSpec;
import me.snowdrop.istio.mixer.template.metric.Metric;
import me.snowdrop.istio.tests.BaseIstioTest;
import org.junit.Test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertNull;
import static org.junit.Assert.assertTrue;

/**
 * @author <a href="claprun@redhat.com">Christophe Laprun</a>
 */
public class MetricTest extends BaseIstioTest {
    @Test
    public void metricRoundtripShouldWork() throws Exception {
        /*
        apiVersion: "config.istio.io/v1alpha2"
        kind: instance
        metadata:
          name: requestsize
          namespace: istio-system
        spec:
          compiledTemplate: metric
          params:
            value: request.size | 0
            dimensions:
              source_version: source.labels["version"] | "unknown"
              destination_service: destination.service.host | "unknown"
              destination_version: destination.labels["version"] | "unknown"
              response_code: response.code | 200
            monitored_resource_type: '"UNSPECIFIED"'
         */
    
        Instance resource = new InstanceBuilder()
            .withNewMetadata().withName("requestsize").withNamespace("istio-system").endMetadata()
            .withNewSpec().withCompiledTemplate("metric")
            .withNewMetric()
            .withValue(TypedValue.from("request.size | 0"))
            .addToDimensions("sourceService", TypedValue.from("source.service | \"unknown\""))
            .addToDimensions("sourceVersion", TypedValue.from("source.labels[\"version\"] | \"unknown\""))
            .addToDimensions("destinationService", TypedValue.from("destination.service | \"unknown\""))
            .addToDimensions("destinationVersion", TypedValue.from("destination.labels[\"version\"] | \"unknown\""))
            .addToDimensions("responseCode", TypedValue.from("response.code | 200"))
            .withMonitoredResourceType("UNSPECIFIED")
            .endMetric()
            .endSpec().build();
    
        final String output = mapper.writeValueAsString(resource);
    
        HasMetadata reloaded = mapper.readValue(output, HasMetadata.class);
    
        assertEquals(resource, reloaded);
    }
    
    @Test
    public void loadingFromYAMLShouldWork() throws Exception {
        ClassLoader classloader = Thread.currentThread().getContextClassLoader();
        InputStream is = classloader.getResourceAsStream("metric.yaml");
        final HasMetadata resource = mapper.readValue(is, HasMetadata.class);
        
        assertEquals(resource.getKind(), "instance");
        
        assertTrue(resource instanceof Instance);
        final InstanceSpec spec = ((Instance) resource).getSpec();
        assertEquals("metric", spec.getCompiledTemplate());
        
        final Metric metric = spec.getMetric();
        assertNotNull(metric);
        assertNull(spec.getApikey());
        assertNull(spec.getAuthorization());
        assertNull(spec.getChecknothing());
        assertNull(spec.getEdge());
        assertNull(spec.getKubernetes());
        assertNull(spec.getListentry());
        assertNull(spec.getLogentry());
        assertNull(spec.getQuota());
        assertNull(spec.getReportnothing());
        assertNull(spec.getTracespan());
        
        assertEquals("1", metric.getValue().getExpression());
        final Map<String, TypedValue> dimensions = metric.getDimensions();
        assertEquals(4, dimensions.size());
        assertTrue(dimensions.containsKey("source"));
        assertEquals("source.service | \"unknown\"", dimensions.get("source").getExpression());
        assertTrue(dimensions.containsKey("destination"));
        assertTrue(dimensions.containsKey("version"));
        assertTrue(dimensions.containsKey("user_agent"));
    }
}
