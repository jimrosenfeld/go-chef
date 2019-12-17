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
	assert.NotNil(t, policies)

	cookbookA := policies["cookbook_a"]
	assert.NotNil(t, cookbookA)
	assert.Equal(t, "cookbook_a", cookbookA.Name)
	assert.Equal(t, "https://localhost/organizations/org/policies/cookbook_a", cookbookA.Uri)
	assert.NotNil(t, cookbookA.Revisions)
	assert.Equal(t, "8701f8de2ed3c7bf01a9f2bec8b2f722be0806c6253b154cde63dcb3e436694a", cookbookA.Revisions[0])

	cookbookB := policies["cookbook_b"]
	assert.NotNil(t, cookbookB)
	assert.Equal(t, "cookbook_b", cookbookB.Name)
	assert.Equal(t, "https://localhost/organizations/org/policies/cookbook_b", cookbookB.Uri)
	assert.NotNil(t, cookbookB.Revisions)
	assert.Equal(t, "05fba2a9ec7c1d638899f1e3a0d91d364dd4725414b1f78f756f5f13c67fe474", cookbookB.Revisions[0])

}
