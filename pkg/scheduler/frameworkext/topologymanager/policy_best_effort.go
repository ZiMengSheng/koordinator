/*
Copyright 2022 The Koordinator Authors.
Copyright 2019 The Kubernetes Authors.

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

package topologymanager

import apiext "github.com/koordinator-sh/koordinator/apis/extension"

type bestEffortPolicy struct {
	//List of NUMA Nodes available on the underlying machine
	numaNodes []int
}

var _ Policy = &bestEffortPolicy{}

// PolicyBestEffort policy name.
const PolicyBestEffort string = "best-effort"

// NewBestEffortPolicy returns best-effort policy.
func NewBestEffortPolicy(numaNodes []int) Policy {
	return &bestEffortPolicy{numaNodes: numaNodes}
}

func (p *bestEffortPolicy) Name() string {
	return PolicyBestEffort
}

func (p *bestEffortPolicy) canAdmitPodResult(hint *NUMATopologyHint) bool {
	return true
}

/*
1. 优先选择 preferred hint，即所有资源都 prefer 且所有资源的 NUMA Node 是对齐的；
2. 如果都是 prefer 或者都不是 prefer，选择 bitwise AND 出的 NUMANode 更少的
3. 如果 prefer 一致且 merge 出的 NUMANode 一致，那么选择 score 更高的
*/

func (p *bestEffortPolicy) Merge(providersHints []map[string][]NUMATopologyHint, exclusivePolicy apiext.NumaTopologyExclusive, allNUMANodeStatus []apiext.NumaNodeStatus) (NUMATopologyHint, bool) {
	filteredProvidersHints := filterProvidersHints(providersHints)
	bestHint := mergeFilteredHints(p.numaNodes, filteredProvidersHints, exclusivePolicy, allNUMANodeStatus)
	admit := p.canAdmitPodResult(&bestHint)
	return bestHint, admit
}
