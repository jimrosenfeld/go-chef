package chef

type PolicyService struct {
	client *Client
}

type apiPolicyRevisionDetail interface{}

type apiPolicyRevisions map[string]apiPolicyRevisionDetail

type apiPolicy struct {
	Uri       string             `json:"uri"`
	Revisions apiPolicyRevisions `json:"revisions"`
}

type apiListPoliciesResult map[string]apiPolicy

type Policy struct {
	Name      string
	Uri       string
	Revisions []string
}

type ListPoliciesResult map[string]Policy

func (p *PolicyService) List() (result ListPoliciesResult, err error) {
	var alpr apiListPoliciesResult
	err = p.client.magicRequestDecoder("GET", "policies", nil, &alpr)
	if err != nil {
		return
	}

	result = make(ListPoliciesResult, len(alpr))
	for apiPolicy := range alpr {
		revisions := make([]string, len(alpr[apiPolicy].Revisions))
		i := 0
		for revision := range alpr[apiPolicy].Revisions {
			revisions[i] = revision
			i++
		}
		policy := Policy{
			Name:      apiPolicy,
			Uri:       alpr[apiPolicy].Uri,
			Revisions: revisions,
		}
		result[apiPolicy] = policy
	}
	return
}
