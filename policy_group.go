package chef

type apiPolicyGroupPolicy struct {
	RevisionID string `json:"revision_id,omitempty"`
}

type apiPolicyGroup struct {
	Uri      string `json:"uri,omitempty"`
	Policies map[string]apiPolicyGroupPolicy
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
