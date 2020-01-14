package chef

import "fmt"

type apiPolicyGroupPolicy struct {
	RevisionID string `json:"revision_id,omitempty"`
}

type apiPolicyGroup struct {
	Uri      string                          `json:"uri,omitempty"`
	Policies map[string]apiPolicyGroupPolicy `json:"policies,omitempty"`
}

type apiListPolicyGroupsResult map[string]apiPolicyGroup

type CurrentPolicy struct {
	Name       string
	RevisionID string
}

type PolicyGroup struct {
	Name            string
	Uri             string
	CurrentPolicies map[string]CurrentPolicy
}

type ListPolicyGroupsResult map[string]PolicyGroup

func (p *PolicyService) ListGroups() (result ListPolicyGroupsResult, err error) {
	var alpgr apiListPolicyGroupsResult
	err = p.client.magicRequestDecoder("GET", "policy_groups", nil, &alpgr)
	if err != nil {
		return
	}

	result = make(ListPolicyGroupsResult, len(alpgr))
	for apiPolicyGroup := range alpgr {
		currentPolicies := make(map[string]CurrentPolicy, len(alpgr[apiPolicyGroup].Policies))
		for apiPolicyGroupPolicy := range alpgr[apiPolicyGroup].Policies {
			currentPolicy := CurrentPolicy{
				Name:       apiPolicyGroupPolicy,
				RevisionID: alpgr[apiPolicyGroup].Policies[apiPolicyGroupPolicy].RevisionID,
			}
			currentPolicies[apiPolicyGroupPolicy] = currentPolicy
		}
		policyGroup := PolicyGroup{
			Name:            apiPolicyGroup,
			Uri:             alpgr[apiPolicyGroup].Uri,
			CurrentPolicies: currentPolicies,
		}
		result[apiPolicyGroup] = policyGroup
	}
	return
}

func (p *PolicyService) GetGroup(name string) (group PolicyGroup, err error) {
	path := fmt.Sprintf("policy_groups/%v", name)
	var apg apiPolicyGroup
	err = p.client.magicRequestDecoder("GET", path, nil, &apg)
	if err != nil {
		return
	}

	group = PolicyGroup{Name: name}
	group.CurrentPolicies = make(map[string]CurrentPolicy, len(apg.Policies))
	for apiPolicyGroupPolicy := range apg.Policies {
		currentPolicy := CurrentPolicy{
			Name:       apiPolicyGroupPolicy,
			RevisionID: apg.Policies[apiPolicyGroupPolicy].RevisionID,
		}
		group.CurrentPolicies[apiPolicyGroupPolicy] = currentPolicy
	}

	return
}
