package chef

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const policiesListResponseFile = "test/policies_list.json"

func TestListPolicies(t *testing.T) {
	setup()
	defer teardown()

	testResponse, err := ioutil.ReadFile(policiesListResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/policies", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(testResponse))
	})

	policies, err := client.Policies.List()
	assert.Nil(t, err)
	if assert.NotNil(t, policies) {
		assert.NotNil(t, policies["cookbook_a"])
		assert.Equal(t, "https://localhost/organizations/org/policies/cookbook_a", policies["cookbook_a"].Uri)
		assert.NotNil(t, policies["cookbook_a"].Revisions)
		assert.Equal(t, "8701f8de2ed3c7bf01a9f2bec8b2f722be0806c6253b154cde63dcb3e436694a", policies["cookbook_a"].Revisions[0])

		assert.NotNil(t, policies["cookbook_b"])
		assert.Equal(t, "https://127.0.0.1/organizations/org/policies/cookbook_b", policies["cookbook_b"].Uri)
		assert.NotNil(t, policies["cookbook_b"].Revisions)
		assert.Equal(t, "05fba2a9ec7c1d638899f1e3a0d91d364dd4725414b1f78f756f5f13c67fe474", policies["cookbook_b"].Revisions[0])
	}
}
