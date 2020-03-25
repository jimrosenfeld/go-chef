package chef

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const policyfileLockResponseFile = "test/policyfile.lock.json"

func TestGetPolicyfileLock(t *testing.T) {
	setup()
	defer teardown()

	testResponse, err := ioutil.ReadFile(policyfileLockResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/policy_groups/test-group/policies/jenkins", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(testResponse))
	})

	p, err := client.Policies.GetPolicyfileLock("jenkins", "test-group")
	assert.Nil(t, err)
	checkPolicyFile(t, p)
}

func TestGetPolicyfileLockFromFile(t *testing.T) {
	setup()
	defer teardown()

	p, err := client.Policies.GetPolicyfileLockFromFile(policyfileLockResponseFile)
	assert.Nil(t, err)
	checkPolicyFile(t, p)
}

func checkPolicyFile(t *testing.T, p PolicyfileLock) {
	assert.NotNil(t, p)

	// basic
	assert.Equal(t, "304566f86a620aae85797a3c491a51fb8c6ecf996407e77b8063aa3ee59672c5", p.RevisionID)
	assert.Equal(t, "jenkins", p.Name)

	// run list
	assert.Equal(t, "recipe[apt::default]", p.RunList[0])
	assert.Equal(t, 4, len(p.RunList))

	// an included policy lock
	ipl := p.IncludedPolicyLocks
	assert.NotNil(t, ipl)
	assert.Equal(t, 1, len(ipl))
	assert.Equal(t, "06da9a52b6f4ff79232a7dc77ade9766470ea51bf0b08d38bae701dbb0899dd1", ipl[0].RevisionID)
	assert.Equal(t, "included_policy", ipl[0].Name)
	assert.Equal(t, "included_policy.lock.json", ipl[0].SourceOptions["path"])

	// an internal cookbook
	assert.NotNil(t, p.CookbookLocks["policyfile_demo"])
	cl := p.CookbookLocks["policyfile_demo"]
	assert.Equal(t, "0.1.0", cl.Version)
	assert.Equal(t, "ea96c99da079db9ff3cb22601638fabd5df49599", cl.Identifier)
	assert.Equal(t, "66030937227426267.45022575077627448.275691232073113", cl.DottedDecimalIdentifier)
	assert.Equal(t, "cookbooks/policyfile_demo", cl.Source)
	assert.Equal(t, "", cl.CacheKey)
	assert.NotNil(t, cl.SCMInfo)
	si := cl.SCMInfo
	assert.Equal(t, "git", si.SCM)
	assert.Equal(t, "git@github.com:danielsdeleo/policyfile-jenkins-demo.git", si.Remote)
	assert.Equal(t, "cf0885f3f2f5edaa44bf8d5e5de4c4d0efa51411", si.Revision)
	assert.Equal(t, false, si.WorkingTreeClean)
	assert.True(t, si.Published)
	assert.Equal(t, "mine/master", si.SynchronizedRemoteBranches[0])
	assert.Equal(t, "cookbooks/policyfile_demo", cl.SourceOptions["path"])

	// an external cookbook
	assert.NotNil(t, p.CookbookLocks["nginx"])
	cl = p.CookbookLocks["nginx"]
	assert.Equal(t, "2.7.6", cl.Version)
	assert.Equal(t, "0dc33e2f3660b7865bad5878c0fcd401348612f9", cl.Identifier)
	assert.Equal(t, "3873846544720055.37818446951006460.233101641257721", cl.DottedDecimalIdentifier)
	assert.Equal(t, "nginx-2.7.6-supermarket.chef.io", cl.CacheKey)
	assert.Equal(t, "https://supermarket.chef.io/api/v1/cookbooks/nginx/versions/2.7.6/download", cl.SourceOptions["artifactserver"])
	assert.Equal(t, "2.7.6", cl.SourceOptions["version"])

	// attributes
	assert.Equal(t, "Attributes, f*** yeah", p.DefaultAttributes["greeting"])
	assert.Equal(t, "use -a", p.OverrideAttributes["attr_only_updating"])

	// dependencies
	policyfileNginxConstraint := []string{"nginx", "= 2.7.6"}
	assert.Contains(t, p.SolutionDependencies.Policyfile, policyfileNginxConstraint)
	policyfileDemoConstraint := []string{"policyfile_demo", ">= 0.0.0"}
	assert.Contains(t, p.SolutionDependencies.Policyfile, policyfileDemoConstraint)
	dependenciesRunitConstraint := []string{"runit", "~> 1.2"}
	assert.Contains(t, p.SolutionDependencies.Dependencies["nginx (2.7.6)"], dependenciesRunitConstraint)
}
