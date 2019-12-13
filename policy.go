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

type listApiPoliciesResult map[string]apiPolicy

type Policy struct {
	Uri       string
	Revisions []string
}

type ListPoliciesResult map[string]Policy

func (p *PolicyService) List() (result ListPoliciesResult, err error) {
	var listApiPoliciesResult listApiPoliciesResult
	err = p.client.magicRequestDecoder("GET", "policies", nil, &listApiPoliciesResult)

	result = make(ListPoliciesResult, len(listApiPoliciesResult))
	for p := range listApiPoliciesResult {
		revisions := make([]string, len(listApiPoliciesResult[p].Revisions))
		i := 0
		for r := range listApiPoliciesResult[p].Revisions {
			revisions[i] = r
			i++
		}
		policy := Policy{
			Uri:       listApiPoliciesResult[p].Uri,
			Revisions: revisions,
		}
		result[p] = policy
	}
	return
}
