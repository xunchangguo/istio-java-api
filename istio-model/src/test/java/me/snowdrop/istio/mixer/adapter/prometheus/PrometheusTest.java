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
package me.snowdrop.istio.mixer.adapter.prometheus;

import java.io.InputStream;

import io.fabric8.kubernetes.api.model.HasMetadata;
import io.fabric8.kubernetes.api.model.KubernetesResource;
import me.snowdrop.istio.tests.BaseIstioTest;
import org.junit.Test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;

/**
 * @author <a href="claprun@redhat.com">Christophe Laprun</a>
 */
public class PrometheusTest extends BaseIstioTest {
    @Test
    public void metricRoundtripShouldWork() throws Exception {
        /*
    apiVersion: "config.istio.io/v1alpha2"
kind: prometheus
metadata:
  name: recommendationrequestcounthandler
  namespace: istio-system
spec:
  metrics:
  - name: recommendation_request_count
    instance_name: recommendationrequestcount.metric.istio-system
    kind: COUNTER
    label_names:
    - source
    - destination
    - user_agent
    - version
     */

        Prometheus prometheus = new PrometheusBuilder()
                .withNewMetadata()
                .withName("recommendationrequestcounthandler")
                .withNamespace("istio-system")
                .endMetadata()
                .withNewSpec()
                    .addNewMetric()
                        .withKind(Kind.COUNTER)
                        .withInstanceName("recommendationrequestcount.metric.istio-system")
                        .addToLabelNames("source", "destination", "user_agent", "version")
                    .endMetric()
                .endSpec()
                .build();


        final String output = mapper.writeValueAsString(prometheus);
        KubernetesResource reloaded = mapper.readValue(output, KubernetesResource.class);
        assertEquals(prometheus, reloaded);
    }

    @Test
    public void loadingFromYAMLShouldWork() throws Exception {
        ClassLoader classloader = Thread.currentThread().getContextClassLoader();
        InputStream is = classloader.getResourceAsStream("prometheus.yaml");
        final HasMetadata resource = mapper.readValue(is, HasMetadata.class);

        assertEquals(resource.getKind(), "prometheus");

        assertTrue(resource instanceof me.snowdrop.istio.mixer.adapter.prometheus.Prometheus);

        final Prometheus prometheus = (Prometheus) resource;
        assertEquals(1, prometheus.getSpec().getMetrics().size());
        final MetricInfo metricInfo = prometheus.getSpec().getMetrics().get(0);
        assertEquals("recommendation_request_count", metricInfo.getName());
        assertEquals("recommendationrequestcount.metric.istio-system", metricInfo.getInstanceName());
        assertEquals(Kind.COUNTER, metricInfo.getKind());
        assertEquals(4, metricInfo.getLabelNames().size());
    }
}
