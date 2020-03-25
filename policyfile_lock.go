package chef

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type PolicyfileLock struct {
	RevisionID           string                  `json:"revision_id,omitempty"`
	Name                 string                  `json:"name,omitempty"`
	RunList              []string                `json:"run_list,omitempty"`
	IncludedPolicyLocks  []IncludedPolicyLock    `json:"included_policy_locks,omitempty"`
	CookbookLocks        map[string]CookbookLock `json:"cookbook_locks,omitempty"`
	DefaultAttributes    map[string]string       `json:"default_attributes,omitempty"`
	OverrideAttributes   map[string]string       `json:"override_attributes,omitempty"`
	SolutionDependencies SolutionDependencies    `json:"solution_dependencies,omitempty"`
}

type CookbookLock struct {
	Version                 string            `json:"version,omitempty"`
	Identifier              string            `json:"identifier,omitempty"`
	DottedDecimalIdentifier string            `json:"dotted_decimal_identifier,omitempty"`
	Source                  string            `json:"source,omitempty"`
	CacheKey                string            `json:"cache_key,omitempty"`
	Origin                  string            `json:"origin,omitempty"`
	SCMInfo                 SCMInfo           `json:"scm_info,omitempty"`
	SourceOptions           map[string]string `json:"source_options,omitempty"`
}

type SCMInfo struct {
	SCM                        string   `json:"scm,omitempty"`
	Remote                     string   `json:"remote,omitempty"`
	Revision                   string   `json:"revision,omitempty"`
	WorkingTreeClean           bool     `json:"working_tree_clean,omitempty"`
	Published                  bool     `json:"published,omitempty"`
	SynchronizedRemoteBranches []string `json:"synchronized_remote_branches,omitempty"`
}

type IncludedPolicyLock struct {
	Name          string            `json:"name,omitempty"`
	RevisionID    string            `json:"revision_id,omitempty"`
	SourceOptions map[string]string `json:"source_options,omitempty"`
}

type SolutionDependencies struct {
	Policyfile   [][]string            `json:"Policyfile,omitempty"`
	Dependencies map[string][][]string `json:"dependencies,omitempty"`
}

func (p *PolicyService) GetPolicyfileLock(policy string, policyGroup string) (result PolicyfileLock, err error) {
	requestPath := fmt.Sprintf("policy_groups/%v/policies/%v", policyGroup, policy)
	err = p.client.magicRequestDecoder("GET", requestPath, nil, &result)
	return
}

func (p *PolicyService) GetPolicyfileLockFromFile(path string) (result PolicyfileLock, err error) {
	lockFile, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(lockFile, &result)
	if err != nil {
		return
	}

	return
}
