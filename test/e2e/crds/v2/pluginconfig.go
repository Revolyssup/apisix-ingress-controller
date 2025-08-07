// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package v2

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	apiv2 "github.com/apache/apisix-ingress-controller/api/v2"
	"github.com/apache/apisix-ingress-controller/test/e2e/framework"
	"github.com/apache/apisix-ingress-controller/test/e2e/scaffold"
)

var _ = FDescribe("Test ApisixPluginConfig", Label("apisix.apache.org", "v2", "apisixpluginconfig"), func() {
	var (
		s = scaffold.NewScaffold(&scaffold.Options{
			ControllerName: fmt.Sprintf("apisix.apache.org/apisix-ingress-controller-%d", time.Now().Unix()),
		})
		applier = framework.NewApplier(s.GinkgoT, s.K8sClient, s.CreateResourceFromString)
	)

	Context("Test ApisixPluginConfig", func() {
		BeforeEach(func() {
			By("create GatewayProxy")
			gatewayProxyName := s.UniqueNameRegistry.Get(scaffold.GatewayProxy)
			gatewayProxy := getGatewayProxyYaml(gatewayProxyName, s.Namespace(), s.Deployer.GetAdminEndpoint(), s.AdminKey())
			err := s.CreateResourceFromStringWithNamespace(gatewayProxy, s.Namespace())
			Expect(err).NotTo(HaveOccurred(), "creating GatewayProxy")
			time.Sleep(5 * time.Second)

			By("create IngressClass")
			ingressClassName := s.UniqueNameRegistry.Get(scaffold.IngressClass)
			err = s.CreateResourceFromStringWithNamespace(getIngressClassYaml(ingressClassName, s.GetControllerName(), gatewayProxyName, s.Namespace()), "")
			Expect(err).NotTo(HaveOccurred(), "creating IngressClass")
			time.Sleep(5 * time.Second)
		})

		It("Basic ApisixPluginConfig test", func() {
			pluginConfigName := s.UniqueNameRegistry.Get("plugin-config")
			routeName := s.UniqueNameRegistry.Get("test-route")
			ingressClassName := s.UniqueNameRegistry.Get(scaffold.IngressClass)
			const apisixPluginConfigSpec = `
apiVersion: apisix.apache.org/v2
kind: ApisixPluginConfig
metadata:
  name: %s
spec:
  ingressClassName: %s
  plugins:
  - name: response-rewrite
    enable: true
    config:
      headers:
        X-Plugin-Config: "test-response-rewrite"
        X-Plugin-Test: "enabled"
`

			const apisixRouteSpec = `
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: %s
spec:
  ingressClassName: %s
  http:
  - name: rule0
    match:
      paths:
      - /*
    backends:
    - serviceName: httpbin-service-e2e-test
      servicePort: 80
    plugin_config_name: %s
`

			By("apply ApisixPluginConfig")
			var apisixPluginConfig apiv2.ApisixPluginConfig
			applier.MustApplyAPIv2(types.NamespacedName{Namespace: s.Namespace(), Name: pluginConfigName},
				&apisixPluginConfig, fmt.Sprintf(apisixPluginConfigSpec, pluginConfigName, ingressClassName))

			By("apply ApisixRoute that references ApisixPluginConfig")
			var apisixRoute apiv2.ApisixRoute
			applier.MustApplyAPIv2(types.NamespacedName{Namespace: s.Namespace(), Name: routeName},
				&apisixRoute, fmt.Sprintf(apisixRouteSpec, routeName, ingressClassName, pluginConfigName))

			By("verify ApisixRoute works with plugin config")
			request := func() int {
				return s.NewAPISIXClient().GET("/get").Expect().Raw().StatusCode
			}
			Eventually(request).WithTimeout(30 * time.Second).ProbeEvery(1 * time.Second).Should(Equal(http.StatusOK))

			By("verify plugin from ApisixPluginConfig works")
			resp := s.NewAPISIXClient().GET("/get").Expect().Status(http.StatusOK)
			resp.Header("X-Plugin-Config").IsEqual("test-response-rewrite")
			resp.Header("X-Plugin-Test").IsEqual("enabled")

			By("delete ApisixRoute")
			err := s.DeleteResource("ApisixRoute", routeName)
			Expect(err).ShouldNot(HaveOccurred(), "deleting ApisixRoute")

			By("delete ApisixPluginConfig")
			err = s.DeleteResource("ApisixPluginConfig", pluginConfigName)
			Expect(err).ShouldNot(HaveOccurred(), "deleting ApisixPluginConfig")

			Eventually(request).WithTimeout(30 * time.Second).ProbeEvery(1 * time.Second).Should(Equal(http.StatusNotFound))
		})

		It("Test ApisixPluginConfig update", func() {
			pluginConfigName := s.UniqueNameRegistry.Get("plugin-config-update")
			ingressClassName := s.UniqueNameRegistry.Get(scaffold.IngressClass)
			const apisixPluginConfigSpecV1 = `
apiVersion: apisix.apache.org/v2
kind: ApisixPluginConfig
metadata:
  name: %s
spec:
  ingressClassName: %s
  plugins:
  - name: response-rewrite
    enable: true
    config:
      headers:
        X-Version: "v1"
`

			const apisixPluginConfigSpecV2 = `
apiVersion: apisix.apache.org/v2
kind: ApisixPluginConfig
metadata:
  name: %s
spec:
  ingressClassName: %s
  plugins:
  - name: response-rewrite
    enable: true
    config:
      headers:
        X-Version: "v2"
        X-Updated: "true"
`

			const apisixRouteSpec = `
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: %s
spec:
  ingressClassName: %s
  http:
  - name: rule0
    match:
      paths:
      - /*
    backends:
    - serviceName: httpbin-service-e2e-test
      servicePort: 80
    plugin_config_name: %s
`

			By("apply initial ApisixPluginConfig")
			routeName := s.UniqueNameRegistry.Get("test-route-update")
			var apisixPluginConfig apiv2.ApisixPluginConfig
			applier.MustApplyAPIv2(types.NamespacedName{Namespace: s.Namespace(), Name: pluginConfigName},
				&apisixPluginConfig, fmt.Sprintf(apisixPluginConfigSpecV1, pluginConfigName, ingressClassName))

			By("apply ApisixRoute that references ApisixPluginConfig")
			var apisixRoute apiv2.ApisixRoute
			applier.MustApplyAPIv2(types.NamespacedName{Namespace: s.Namespace(), Name: routeName},
				&apisixRoute, fmt.Sprintf(apisixRouteSpec, routeName, ingressClassName, pluginConfigName))

			By("verify initial plugin config works")
			request := func() int {
				return s.NewAPISIXClient().GET("/get").Expect().Raw().StatusCode
			}
			Eventually(request).WithTimeout(30 * time.Second).ProbeEvery(1 * time.Second).Should(Equal(http.StatusOK))

			resp := s.NewAPISIXClient().GET("/get").Expect().Status(http.StatusOK)
			resp.Header("X-Version").IsEqual("v1")
			resp.Header("X-Updated").IsEmpty()

			By("update ApisixPluginConfig")
			applier.MustApplyAPIv2(types.NamespacedName{Namespace: s.Namespace(), Name: pluginConfigName},
				&apisixPluginConfig, fmt.Sprintf(apisixPluginConfigSpecV2, pluginConfigName, ingressClassName))
			time.Sleep(5 * time.Second)

			By("verify updated plugin config works")
			resp = s.NewAPISIXClient().GET("/get").Expect().Status(http.StatusOK)
			resp.Header("X-Version").IsEqual("v2")
			resp.Header("X-Updated").IsEqual("true")

			By("delete resources")
			err := s.DeleteResource("ApisixRoute", routeName)
			Expect(err).ShouldNot(HaveOccurred(), "deleting ApisixRoute")
			err = s.DeleteResource("ApisixPluginConfig", pluginConfigName)
			Expect(err).ShouldNot(HaveOccurred(), "deleting ApisixPluginConfig")
		})

		It("Test ApisixPluginConfig with disabled plugin", func() {
			pluginConfigName := s.UniqueNameRegistry.Get("plugin-config-disabled")
			ingressClassName := s.UniqueNameRegistry.Get(scaffold.IngressClass)
			routeName := s.UniqueNameRegistry.Get("test-route-disabled")
			const apisixPluginConfigSpec = `
apiVersion: apisix.apache.org/v2
kind: ApisixPluginConfig
metadata:
  name: %s
spec:
  ingressClassName: %s
  plugins:
  - name: response-rewrite
    enable: false
    config:
      headers:
        X-Should-Not-Exist: "disabled"
  - name: cors
    enable: true
    config:
      allow_origins: "*"
      allow_methods: "GET,POST"
`

			const apisixRouteSpec = `
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: %s
spec:
  ingressClassName: %s
  http:
  - name: rule0
    match:
      paths:
      - /*
    backends:
    - serviceName: httpbin-service-e2e-test
      servicePort: 80
    plugin_config_name: %s
`

			By("apply ApisixPluginConfig with disabled plugin")
			var apisixPluginConfig apiv2.ApisixPluginConfig
			applier.MustApplyAPIv2(types.NamespacedName{Namespace: s.Namespace(), Name: pluginConfigName},
				&apisixPluginConfig, fmt.Sprintf(apisixPluginConfigSpec, pluginConfigName, ingressClassName))

			By("apply ApisixRoute that references ApisixPluginConfig")
			var apisixRoute apiv2.ApisixRoute
			applier.MustApplyAPIv2(types.NamespacedName{Namespace: s.Namespace(), Name: routeName},
				&apisixRoute, fmt.Sprintf(apisixRouteSpec, routeName, ingressClassName, pluginConfigName))

			By("verify ApisixRoute works")
			request := func() int {
				return s.NewAPISIXClient().GET("/get").Expect().Raw().StatusCode
			}
			Eventually(request).WithTimeout(30 * time.Second).ProbeEvery(1 * time.Second).Should(Equal(http.StatusOK))

			By("verify disabled plugin is not applied")
			resp := s.NewAPISIXClient().GET("/get").Expect().Status(http.StatusOK)
			resp.Header("X-Should-Not-Exist").IsEmpty()

			By("verify enabled plugin is applied")
			resp.Header("Access-Control-Allow-Origin").IsEqual("*")

			By("delete resources")
			err := s.DeleteResource("ApisixRoute", routeName)
			Expect(err).ShouldNot(HaveOccurred(), "deleting ApisixRoute")
			err = s.DeleteResource("ApisixPluginConfig", pluginConfigName)
			Expect(err).ShouldNot(HaveOccurred(), "deleting ApisixPluginConfig")
		})

		It("Test ApisixPluginConfig overridden by route plugins", func() {
			pluginConfigName := s.UniqueNameRegistry.Get("plugin-config-override")
			routeName := s.UniqueNameRegistry.Get("test-route-override")
			ingressClassName := s.UniqueNameRegistry.Get(scaffold.IngressClass)
			const apisixPluginConfigSpec = `
apiVersion: apisix.apache.org/v2
kind: ApisixPluginConfig
metadata:
  name: %s
spec:
  ingressClassName: %s
  plugins:
  - name: response-rewrite
    enable: true
    config:
      headers:
        X-From-Config: "plugin-config"
        X-Shared: "from-config"
`

			const apisixRouteSpec = `
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: %s
spec:
  ingressClassName: %s
  http:
  - name: rule0
    match:
      paths:
      - /*
    backends:
    - serviceName: httpbin-service-e2e-test
      servicePort: 80
    plugin_config_name: %s
    plugins:
    - name: response-rewrite
      enable: true
      config:
        headers:
          X-From-Route: "route"
          X-Shared: "from-route"
`

			By("apply ApisixPluginConfig")
			var apisixPluginConfig apiv2.ApisixPluginConfig
			applier.MustApplyAPIv2(types.NamespacedName{Namespace: s.Namespace(), Name: pluginConfigName},
				&apisixPluginConfig, fmt.Sprintf(apisixPluginConfigSpec, pluginConfigName, ingressClassName))

			By("apply ApisixRoute with overriding plugins")
			var apisixRoute apiv2.ApisixRoute
			applier.MustApplyAPIv2(types.NamespacedName{Namespace: s.Namespace(), Name: routeName},
				&apisixRoute, fmt.Sprintf(apisixRouteSpec, routeName, ingressClassName, pluginConfigName))

			By("verify ApisixRoute works")
			request := func() int {
				return s.NewAPISIXClient().GET("/get").Expect().Raw().StatusCode
			}
			Eventually(request).WithTimeout(30 * time.Second).ProbeEvery(1 * time.Second).Should(Equal(http.StatusOK))

			By("verify route plugins override plugin config")
			resp := s.NewAPISIXClient().GET("/get").Expect().Status(http.StatusOK)
			resp.Header("X-From-Config").IsEmpty()
			resp.Header("X-From-Route").IsEqual("route")
			resp.Header("X-Shared").IsEqual("from-route")

			By("delete resources")
			err := s.DeleteResource("ApisixRoute", routeName)
			Expect(err).ShouldNot(HaveOccurred(), "deleting ApisixRoute")
			err = s.DeleteResource("ApisixPluginConfig", pluginConfigName)
			Expect(err).ShouldNot(HaveOccurred(), "deleting ApisixPluginConfig")
		})

		FIt("Test cross-namespace ApisixPluginConfig reference", Serial, func() {
			ingressClassName := s.UniqueNameRegistry.Get(scaffold.IngressClass)
			const crossNamespaceApisixPluginConfigSpec = `
apiVersion: apisix.apache.org/v2
kind: ApisixPluginConfig
metadata:
  name: cross-ns-plugin-config
  namespace: default
spec:
  ingressClassName: %s
  plugins:
  - name: response-rewrite
    enable: true
    config:
      headers:
        X-Cross-Namespace: "true"
        X-Namespace: "default"
`

			const apisixRouteSpec = `
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: test-route-cross-ns
spec:
  ingressClassName: %s
  http:
  - name: rule0
    match:
      paths:
      - /*
    backends:
    - serviceName: httpbin-service-e2e-test
      servicePort: 80
    plugin_config_name: cross-ns-plugin-config
    plugin_config_namespace: default
`

			By("apply ApisixPluginConfig in default namespace")
			err := s.CreateResourceFromStringWithNamespace(fmt.Sprintf(crossNamespaceApisixPluginConfigSpec, ingressClassName), "default")
			Expect(err).NotTo(HaveOccurred(), "creating default/cross-ns-plugin-config")
			time.Sleep(5 * time.Second)

			By("apply ApisixRoute in test namespace that references ApisixPluginConfig in default namespace")
			var apisixRoute apiv2.ApisixRoute
			applier.MustApplyAPIv2(types.NamespacedName{Namespace: s.Namespace(), Name: "test-route-cross-ns"},
				&apisixRoute, fmt.Sprintf(apisixRouteSpec, ingressClassName))

			By("verify cross-namespace reference works")
			request := func() int {
				return s.NewAPISIXClient().GET("/get").Expect().Raw().StatusCode
			}
			Eventually(request).WithTimeout(30 * time.Second).ProbeEvery(1 * time.Second).Should(Equal(http.StatusOK))

			resp := s.NewAPISIXClient().GET("/get").Expect().Status(http.StatusOK)
			resp.Header("X-Cross-Namespace").IsEqual("true")
			resp.Header("X-Namespace").IsEqual("default")

			By("delete resources")
			err = s.DeleteResource("ApisixRoute", "test-route-cross-ns")
			Expect(err).ShouldNot(HaveOccurred(), "deleting ApisixRoute")
			err = s.DeleteResourceFromStringWithNamespace(fmt.Sprintf(crossNamespaceApisixPluginConfigSpec, ingressClassName), "default")
			Expect(err).ShouldNot(HaveOccurred(), "deleting ApisixPluginConfig")
		})

		It("Test ApisixPluginConfig with SecretRef", func() {
			pluginConfigName := s.UniqueNameRegistry.Get("plugin-config-secret")
			routeName := s.UniqueNameRegistry.Get("test-route-secret")
			ingressClassName := s.UniqueNameRegistry.Get(scaffold.IngressClass)
			const secretSpec = `
apiVersion: v1
kind: Secret
metadata:
  name: plugin-secret
type: Opaque
data:
  key: dGVzdC1rZXk=
  username: dGVzdC11c2Vy
  password: dGVzdC1wYXNzd29yZA==
`

			const apisixPluginConfigSpec = `
apiVersion: apisix.apache.org/v2
kind: ApisixPluginConfig
metadata:
  name: %s
spec:
  ingressClassName: %s
  plugins:
  - name: response-rewrite
    enable: true
    secretRef: plugin-secret
    config:
      headers:
        X-Secret-Ref: "true"
`

			const apisixRouteSpec = `
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: %s
spec:
  ingressClassName: %s
  http:
  - name: rule0
    match:
      paths:
      - /*
    backends:
    - serviceName: httpbin-service-e2e-test
      servicePort: 80
    plugin_config_name: %s
`

			By("apply Secret")
			err := s.CreateResourceFromStringWithNamespace(secretSpec, s.Namespace())
			Expect(err).NotTo(HaveOccurred(), "creating Secret")

			By("apply ApisixPluginConfig with SecretRef")
			var apisixPluginConfig apiv2.ApisixPluginConfig
			applier.MustApplyAPIv2(types.NamespacedName{Namespace: s.Namespace(), Name: pluginConfigName},
				&apisixPluginConfig, fmt.Sprintf(apisixPluginConfigSpec, pluginConfigName, ingressClassName))

			By("apply ApisixRoute that references ApisixPluginConfig")
			var apisixRoute apiv2.ApisixRoute
			applier.MustApplyAPIv2(types.NamespacedName{Namespace: s.Namespace(), Name: routeName},
				&apisixRoute, fmt.Sprintf(apisixRouteSpec, routeName, ingressClassName, pluginConfigName))

			By("verify ApisixRoute works with SecretRef")
			request := func() int {
				return s.NewAPISIXClient().GET("/get").Expect().Raw().StatusCode
			}
			Eventually(request).WithTimeout(30 * time.Second).ProbeEvery(1 * time.Second).Should(Equal(http.StatusOK))

			resp := s.NewAPISIXClient().GET("/get").Expect().Status(http.StatusOK)
			resp.Header("X-Secret-Ref").IsEqual("true")

			By("delete resources")
			err = s.DeleteResource("ApisixRoute", routeName)
			Expect(err).ShouldNot(HaveOccurred(), "deleting ApisixRoute")
			err = s.DeleteResource("ApisixPluginConfig", pluginConfigName)
			Expect(err).ShouldNot(HaveOccurred(), "deleting ApisixPluginConfig")
			err = s.DeleteResource("Secret", "plugin-secret")
			Expect(err).ShouldNot(HaveOccurred(), "deleting Secret")
		})
	})
})
