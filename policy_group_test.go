package chef

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const policyGroupsListResponseFile = "test/policy_groups_list.json"

func TestListPolicyGroups(t *testing.T) {
	setup()
	defer teardown()

	testResponse, err := ioutil.ReadFile(policyGroupsListResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/policy_groups", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(testResponse))
	})

	pgs, err := client.Policies.ListGroups()
	assert.Nil(t, err)
	assert.NotNil(t, pgs)

	dev := pgs["dev"]
	assert.NotNil(t, dev)
	assert.Equal(t, "dev", dev.Name)
	assert.Equal(t, "https://localhost/organizations/org/policy_groups/dev", dev.Uri)
	assert.Equal(t, "cookbook_a", dev.CurrentPolicies["cookbook_a"].Name)
	assert.Equal(t, "50d924a17c82ec6aa7826c4f3035c75b7a4b7c548a6dab7030df142b22d0f7b2", dev.CurrentPolicies["cookbook_a"].RevisionID)
	assert.Equal(t, "cookbook_b", dev.CurrentPolicies["cookbook_b"].Name)
	assert.Equal(t, "28c1bdcc56355a0a3c94182882a404f9e8b0fb3e1fd8d32a65bfb870b089e476", dev.CurrentPolicies["cookbook_b"].RevisionID)

	prod := pgs["prod"]
	assert.NotNil(t, prod)
	assert.Equal(t, "prod", prod.Name)
	assert.Equal(t, "https://localhost/organizations/org/policy_groups/prod", prod.Uri)
	assert.Equal(t, "cookbook_a", prod.CurrentPolicies["cookbook_a"].Name)
	assert.Equal(t, "69a8a638cf6dd7ba6be3c5fd53beef7bfa6e1c3de2f9bcf036c32ea49ab1bd86", prod.CurrentPolicies["cookbook_a"].RevisionID)
	assert.Equal(t, "cookbook_b", prod.CurrentPolicies["cookbook_b"].Name)
	assert.Equal(t, "2d54a5b1bb89e29aeb3689036a8e823fae907d5523c6725bfc0f8c32b888a348", prod.CurrentPolicies["cookbook_b"].RevisionID)
}
