// Copyright 2018 Istio Authors
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

package v1alpha3

import (
	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/plugin"
)

type ConfigGeneratorImpl struct {
	// List of plugins that modify code generated by this config generator
	Plugins []plugin.Plugin
	// List of outbound clusters for sidecars
	// Must be rebuilt for each push epoch
	// TODO: Get rid of the two types once the old code for global MUTUAL TLS is gone
	OutboundClustersForSidecars []*xdsapi.Cluster
	// List of outbound clusters for gateways
	OutboundClustersForGateways []*xdsapi.Cluster
	// TODO: add others in future
}

func NewConfigGenerator(plugins []plugin.Plugin) *ConfigGeneratorImpl {
	return &ConfigGeneratorImpl{
		Plugins: plugins,
	}
}

func (configgen *ConfigGeneratorImpl) BuildSharedPushState(env *model.Environment, push *model.PushContext) error {
	configgen.OutboundClustersForSidecars = configgen.buildOutboundClusters(env, model.Sidecar, push)
	configgen.OutboundClustersForGateways = configgen.buildOutboundClusters(env, model.Router, push)
	return nil
}