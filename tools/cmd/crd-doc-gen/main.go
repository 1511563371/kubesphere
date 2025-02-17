/*
Copyright 2020 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/meta"
	urlruntime "k8s.io/apimachinery/pkg/util/runtime"

	applicationv1alpha1 "kubesphere.io/api/application/v1alpha1"
	clusterv1alpha1 "kubesphere.io/api/cluster/v1alpha1"
	devopsv1alpha3 "kubesphere.io/api/devops/v1alpha3"

	"kubesphere.io/kubesphere/pkg/version"
	"kubesphere.io/kubesphere/tools/lib"

	"github.com/go-openapi/spec"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kube-openapi/pkg/common"

	applicationinstall "kubesphere.io/api/application/crdinstall"
	devopsinstall "kubesphere.io/api/devops/crdinstall"
	devopsv1alpha1 "kubesphere.io/api/devops/v1alpha1"
	networkinstall "kubesphere.io/api/network/crdinstall"
	networkv1alpha1 "kubesphere.io/api/network/v1alpha1"
	servicemeshinstall "kubesphere.io/api/servicemesh/crdinstall"
	servicemeshv1alpha2 "kubesphere.io/api/servicemesh/v1alpha2"
	tenantinstall "kubesphere.io/api/tenant/crdinstall"
	tenantv1alpha1 "kubesphere.io/api/tenant/v1alpha1"
)

var output string

func init() {
	flag.StringVar(&output, "output", "./api/openapi-spec/swagger.json", "--output=./api/openapi-spec/swagger.json")
}

func main() {

	var (
		Scheme = runtime.NewScheme()
		Codecs = serializer.NewCodecFactory(Scheme)
	)

	servicemeshinstall.Install(Scheme)
	tenantinstall.Install(Scheme)
	networkinstall.Install(Scheme)
	devopsinstall.Install(Scheme)
	applicationinstall.Install(Scheme)

	urlruntime.Must(clusterv1alpha1.AddToScheme(Scheme))
	urlruntime.Must(Scheme.SetVersionPriority(clusterv1alpha1.SchemeGroupVersion))

	mapper := meta.NewDefaultRESTMapper(nil)

	mapper.AddSpecific(servicemeshv1alpha2.SchemeGroupVersion.WithKind(servicemeshv1alpha2.ResourceKindServicePolicy),
		servicemeshv1alpha2.SchemeGroupVersion.WithResource(servicemeshv1alpha2.ResourcePluralServicePolicy),
		servicemeshv1alpha2.SchemeGroupVersion.WithResource(servicemeshv1alpha2.ResourceSingularServicePolicy), meta.RESTScopeRoot)

	mapper.AddSpecific(servicemeshv1alpha2.SchemeGroupVersion.WithKind(servicemeshv1alpha2.ResourceKindStrategy),
		servicemeshv1alpha2.SchemeGroupVersion.WithResource(servicemeshv1alpha2.ResourcePluralStrategy),
		servicemeshv1alpha2.SchemeGroupVersion.WithResource(servicemeshv1alpha2.ResourceSingularStrategy), meta.RESTScopeRoot)

	mapper.AddSpecific(tenantv1alpha1.SchemeGroupVersion.WithKind(tenantv1alpha1.ResourceKindWorkspace),
		tenantv1alpha1.SchemeGroupVersion.WithResource(tenantv1alpha1.ResourcePluralWorkspace),
		tenantv1alpha1.SchemeGroupVersion.WithResource(tenantv1alpha1.ResourceSingularWorkspace), meta.RESTScopeRoot)

	mapper.AddSpecific(devopsv1alpha1.SchemeGroupVersion.WithKind(devopsv1alpha1.ResourceKindS2iBuilder),
		devopsv1alpha1.SchemeGroupVersion.WithResource(devopsv1alpha1.ResourcePluralS2iBuilder),
		devopsv1alpha1.SchemeGroupVersion.WithResource(devopsv1alpha1.ResourceSingularS2iBuilder), meta.RESTScopeRoot)

	mapper.AddSpecific(devopsv1alpha1.SchemeGroupVersion.WithKind(devopsv1alpha1.ResourceKindS2iBuilderTemplate),
		devopsv1alpha1.SchemeGroupVersion.WithResource(devopsv1alpha1.ResourcePluralS2iBuilderTemplate),
		devopsv1alpha1.SchemeGroupVersion.WithResource(devopsv1alpha1.ResourceSingularS2iBuilderTemplate), meta.RESTScopeRoot)

	mapper.AddSpecific(devopsv1alpha1.SchemeGroupVersion.WithKind(devopsv1alpha1.ResourceKindS2iRun),
		devopsv1alpha1.SchemeGroupVersion.WithResource(devopsv1alpha1.ResourcePluralS2iRun),
		devopsv1alpha1.SchemeGroupVersion.WithResource(devopsv1alpha1.ResourceSingularS2iRun), meta.RESTScopeRoot)
	mapper.AddSpecific(devopsv1alpha1.SchemeGroupVersion.WithKind(devopsv1alpha1.ResourceKindS2iBinary),
		devopsv1alpha1.SchemeGroupVersion.WithResource(devopsv1alpha1.ResourceSingularS2iBinary),
		devopsv1alpha1.SchemeGroupVersion.WithResource(devopsv1alpha1.ResourcePluralS2iBinary), meta.RESTScopeRoot)

	mapper.AddSpecific(networkv1alpha1.SchemeGroupVersion.WithKind(networkv1alpha1.ResourceKindNamespaceNetworkPolicy),
		networkv1alpha1.SchemeGroupVersion.WithResource(networkv1alpha1.ResourcePluralNamespaceNetworkPolicy),
		networkv1alpha1.SchemeGroupVersion.WithResource(networkv1alpha1.ResourceSingularNamespaceNetworkPolicy), meta.RESTScopeRoot)
	mapper.AddSpecific(networkv1alpha1.SchemeGroupVersion.WithKind(networkv1alpha1.ResourceKindIPPool),
		networkv1alpha1.SchemeGroupVersion.WithResource(networkv1alpha1.ResourcePluralIPPool),
		networkv1alpha1.SchemeGroupVersion.WithResource(networkv1alpha1.ResourceSingularIPPool), meta.RESTScopeRoot)

	mapper.AddSpecific(devopsv1alpha3.SchemeGroupVersion.WithKind(devopsv1alpha3.ResourceKindDevOpsProject),
		devopsv1alpha3.SchemeGroupVersion.WithResource(devopsv1alpha3.ResourcePluralDevOpsProject),
		devopsv1alpha3.SchemeGroupVersion.WithResource(devopsv1alpha3.ResourceSingularDevOpsProject), meta.RESTScopeRoot)
	mapper.AddSpecific(devopsv1alpha3.SchemeGroupVersion.WithKind(devopsv1alpha3.ResourceKindPipeline),
		devopsv1alpha3.SchemeGroupVersion.WithResource(devopsv1alpha3.ResourcePluralPipeline),
		devopsv1alpha3.SchemeGroupVersion.WithResource(devopsv1alpha3.ResourceSingularPipeline), meta.RESTScopeRoot)

	mapper.AddSpecific(clusterv1alpha1.SchemeGroupVersion.WithKind(clusterv1alpha1.ResourceKindCluster),
		clusterv1alpha1.SchemeGroupVersion.WithResource(clusterv1alpha1.ResourcesPluralCluster),
		clusterv1alpha1.SchemeGroupVersion.WithResource(clusterv1alpha1.ResourcesSingularCluster), meta.RESTScopeRoot)

	mapper.AddSpecific(clusterv1alpha1.SchemeGroupVersion.WithKind(clusterv1alpha1.ResourceKindCluster),
		clusterv1alpha1.SchemeGroupVersion.WithResource(clusterv1alpha1.ResourcesPluralCluster),
		clusterv1alpha1.SchemeGroupVersion.WithResource(clusterv1alpha1.ResourcesSingularCluster), meta.RESTScopeRoot)

	mapper.AddSpecific(applicationv1alpha1.SchemeGroupVersion.WithKind(applicationv1alpha1.ResourceKindHelmApplication),
		applicationv1alpha1.SchemeGroupVersion.WithResource(applicationv1alpha1.ResourcePluralHelmApplication),
		applicationv1alpha1.SchemeGroupVersion.WithResource(applicationv1alpha1.ResourceSingularHelmApplication), meta.RESTScopeRoot)

	mapper.AddSpecific(applicationv1alpha1.SchemeGroupVersion.WithKind(applicationv1alpha1.ResourceKindHelmApplicationVersion),
		applicationv1alpha1.SchemeGroupVersion.WithResource(applicationv1alpha1.ResourcePluralHelmApplicationVersion),
		applicationv1alpha1.SchemeGroupVersion.WithResource(applicationv1alpha1.ResourceSingularHelmApplicationVersion), meta.RESTScopeRoot)

	mapper.AddSpecific(applicationv1alpha1.SchemeGroupVersion.WithKind(applicationv1alpha1.ResourceKindHelmRelease),
		applicationv1alpha1.SchemeGroupVersion.WithResource(applicationv1alpha1.ResourcePluralHelmRelease),
		applicationv1alpha1.SchemeGroupVersion.WithResource(applicationv1alpha1.ResourceSingularHelmRelease), meta.RESTScopeRoot)

	mapper.AddSpecific(applicationv1alpha1.SchemeGroupVersion.WithKind(applicationv1alpha1.ResourceKindHelmRepo),
		applicationv1alpha1.SchemeGroupVersion.WithResource(applicationv1alpha1.ResourcePluralHelmRepo),
		applicationv1alpha1.SchemeGroupVersion.WithResource(applicationv1alpha1.ResourceSingularHelmRepo), meta.RESTScopeRoot)

	mapper.AddSpecific(applicationv1alpha1.SchemeGroupVersion.WithKind(applicationv1alpha1.ResourceKindHelmCategory),
		applicationv1alpha1.SchemeGroupVersion.WithResource(applicationv1alpha1.ResourcePluralHelmCategory),
		applicationv1alpha1.SchemeGroupVersion.WithResource(applicationv1alpha1.ResourceSingularHelmCategory), meta.RESTScopeRoot)

	spec, err := lib.RenderOpenAPISpec(lib.Config{
		Scheme: Scheme,
		Codecs: Codecs,
		Info: spec.InfoProps{
			Title:   "KubeSphere",
			Version: version.Get().GitVersion,
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "KubeSphere",
					URL:   "https://kubesphere.io/",
					Email: "kubesphere@yunify.com",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "Apache 2.0",
					URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
				},
			},
		},
		OpenAPIDefinitions: []common.GetOpenAPIDefinitions{
			servicemeshv1alpha2.GetOpenAPIDefinitions,
			tenantv1alpha1.GetOpenAPIDefinitions,
			networkv1alpha1.GetOpenAPIDefinitions,
			devopsv1alpha1.GetOpenAPIDefinitions,
			devopsv1alpha3.GetOpenAPIDefinitions,
			clusterv1alpha1.GetOpenAPIDefinitions,
		},
		Resources: []schema.GroupVersionResource{
			//TODO（runzexia） At present, the document generation requires the openapi structure of the go language,
			// but there is no +k8s:openapi-gen=true in the repository of https://github.com/knative/pkg,
			// and the api document cannot be generated temporarily.
			//servicemeshv1alpha2.SchemeGroupVersion.WithResource(servicemeshv1alpha2.ResourcePluralStrategy),
			//servicemeshv1alpha2.SchemeGroupVersion.WithResource(servicemeshv1alpha2.ResourcePluralServicePolicy),
			tenantv1alpha1.SchemeGroupVersion.WithResource(tenantv1alpha1.ResourcePluralWorkspace),
			devopsv1alpha1.SchemeGroupVersion.WithResource(devopsv1alpha1.ResourcePluralS2iBinary),
			devopsv1alpha1.SchemeGroupVersion.WithResource(devopsv1alpha1.ResourcePluralS2iRun),
			devopsv1alpha1.SchemeGroupVersion.WithResource(devopsv1alpha1.ResourcePluralS2iBuilderTemplate),
			devopsv1alpha1.SchemeGroupVersion.WithResource(devopsv1alpha1.ResourcePluralS2iBuilder),
			networkv1alpha1.SchemeGroupVersion.WithResource(networkv1alpha1.ResourcePluralNamespaceNetworkPolicy),
			networkv1alpha1.SchemeGroupVersion.WithResource(networkv1alpha1.ResourcePluralIPPool),
			devopsv1alpha3.SchemeGroupVersion.WithResource(devopsv1alpha3.ResourcePluralDevOpsProject),
			devopsv1alpha3.SchemeGroupVersion.WithResource(devopsv1alpha3.ResourcePluralPipeline),
			clusterv1alpha1.SchemeGroupVersion.WithResource(clusterv1alpha1.ResourcesPluralCluster),
		},
		Mapper: mapper,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(filepath.Dir(output), 0755)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(output, []byte(spec), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
