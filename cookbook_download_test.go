//
//  Author:: Salim Afiune <afiune@chef.io>
//

package chef

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const emptyCookbookResponseFile = "test/empty_cookbook.json"

func TestCookbooksDownloadThatDoesNotExist(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/cookbooks/foo/2.1.0", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", 404)
	})

	err := client.Cookbooks.Download("foo", "2.1.0")
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "404")
	}
}

func TestCookbooksDownloadCorrectsLatestVersion(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/cookbooks/foo/_latest", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", 404)
	})

	err := client.Cookbooks.Download("foo", "")
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "404")
	}

	err = client.Cookbooks.Download("foo", "latest")
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "404")
	}

	err = client.Cookbooks.Download("foo", "_latest")
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "404")
	}
}

func TestCookbooksDownloadEmptyWithVersion(t *testing.T) {
	setup()
	defer teardown()

	cbookResp, err := ioutil.ReadFile(emptyCookbookResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/cookbooks/foo/0.2.0", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(cbookResp))
	})

	err = client.Cookbooks.Download("foo", "0.2.0")
	assert.Nil(t, err)
}

func TestCookbooksDownloadTo(t *testing.T) {
	setup()
	defer teardown()

	mockedCookbookResponseFile := `
{
  "version": "0.2.1",
  "name": "foo-0.2.1",
  "cookbook_name": "foo",
  "frozen?": false,
  "chef_type": "cookbook_version",
  "json_class": "Chef::CookbookVersion",
  "attributes": [],
  "definitions": [],
  "files": [],
  "libraries": [],
  "providers": [],
  "recipes": [
    {
      "name": "default.rb",
      "path": "recipes/default.rb",
      "checksum": "8e751ed8663cb9b97499956b6a20b0de",
      "specificity": "default",
      "url": "` + server.URL + `/bookshelf/foo/default_rb"
    }
  ],
  "resources": [],
  "root_files": [
    {
      "name": "metadata.rb",
      "path": "metadata.rb",
      "checksum": "6607f3131919e82dc4ba4c026fcfee9f",
      "specificity": "default",
      "url": "` + server.URL + `/bookshelf/foo/metadata_rb"
    }
  ],
  "templates": [],
  "metadata": {},
  "access": {}
}
`

	tempDir, err := ioutil.TempDir("", "foo-cookbook")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tempDir) // clean up

	mux.HandleFunc("/cookbooks/foo/0.2.1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(mockedCookbookResponseFile))
	})
	mux.HandleFunc("/bookshelf/foo/metadata_rb", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "name 'foo'")
	})
	mux.HandleFunc("/bookshelf/foo/default_rb", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "log 'this is a resource'")
	})

	err = client.Cookbooks.DownloadTo("foo", "0.2.1", tempDir)
	assert.Nil(t, err)

	var (
		cookbookPath = path.Join(tempDir, "foo-0.2.1")
		metadataPath = path.Join(cookbookPath, "metadata.rb")
		recipesPath  = path.Join(cookbookPath, "recipes")
		defaultPath  = path.Join(recipesPath, "default.rb")
	)
	assert.DirExistsf(t, cookbookPath, "the cookbook directory should exist")
	assert.DirExistsf(t, recipesPath, "the recipes directory should exist")
	if assert.FileExistsf(t, metadataPath, "a metadata.rb file should exist") {
		metadataBytes, err := ioutil.ReadFile(metadataPath)
		assert.Nil(t, err)
		assert.Equal(t, "name 'foo'", string(metadataBytes))
	}
	if assert.FileExistsf(t, defaultPath, "the default.rb recipes should exist") {
		recipeBytes, err := ioutil.ReadFile(defaultPath)
		assert.Nil(t, err)
		assert.Equal(t, "log 'this is a resource'", string(recipeBytes))
	}
}

func TestCookbooksDownloadAt(t *testing.T) {
	setup()
	defer teardown()

	mockedCookbookResponseFile := `
{
  "version": "0.2.1",
  "name": "foo-0.2.1",
  "cookbook_name": "foo",
  "frozen?": false,
  "chef_type": "cookbook_version",
  "json_class": "Chef::CookbookVersion",
  "attributes": [],
  "definitions": [],
  "files": [],
  "libraries": [],
  "providers": [],
  "recipes": [
    {
      "name": "default.rb",
      "path": "recipes/default.rb",
      "checksum": "8e751ed8663cb9b97499956b6a20b0de",
      "specificity": "default",
      "url": "` + server.URL + `/bookshelf/foo/default_rb"
    }
  ],
  "resources": [],
  "root_files": [
    {
      "name": "metadata.rb",
      "path": "metadata.rb",
      "checksum": "6607f3131919e82dc4ba4c026fcfee9f",
      "specificity": "default",
      "url": "` + server.URL + `/bookshelf/foo/metadata_rb"
    }
  ],
  "templates": [],
  "metadata": {},
  "access": {}
}
`

        tempDir, err := ioutil.TempDir("", "foo-cookbook")
        if err != nil {
                t.Error(err)
        }
        defer os.RemoveAll(tempDir) // clean up

        mux.HandleFunc("/cookbooks/foo/0.2.1", func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintf(w, string(mockedCookbookResponseFile))
        })
        mux.HandleFunc("/bookshelf/foo/metadata_rb", func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintf(w, "name 'foo'")
        })
        mux.HandleFunc("/bookshelf/foo/default_rb", func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintf(w, "log 'this is a resource'")
        })

        err = client.Cookbooks.DownloadAt("foo", "0.2.1", tempDir)
        assert.Nil(t, err)

        var (
                cookbookPath = path.Join(tempDir, "foo-0.2.1")
                metadataPath = path.Join(cookbookPath, "metadata.rb")
                recipesPath  = path.Join(cookbookPath, "recipes")
                defaultPath  = path.Join(recipesPath, "default.rb")
        )
        assert.DirExistsf(t, cookbookPath, "the cookbook directory should exist")
        assert.DirExistsf(t, recipesPath, "the recipes directory should exist")
        if assert.FileExistsf(t, metadataPath, "a metadata.rb file should exist") {
                metadataBytes, err := ioutil.ReadFile(metadataPath)
                assert.Nil(t, err)
                assert.Equal(t, "name 'foo'", string(metadataBytes))
        }
        if assert.FileExistsf(t, defaultPath, "the default.rb recipes should exist") {
                recipeBytes, err := ioutil.ReadFile(defaultPath)
                assert.Nil(t, err)
                assert.Equal(t, "log 'this is a resource'", string(recipeBytes))
        }
}

func TestVerifyMD5Checksum(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "md5-test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tempDir) // clean up

	var (
		// if someone changes the test data,
		// you have to also update the below md5 sum
		testData = []byte("hello\nchef\n")
		filePath = path.Join(tempDir, "dat")
	)
	err = ioutil.WriteFile(filePath, testData, 0644)
	assert.Nil(t, err)
	assert.True(t, verifyMD5Checksum(filePath, "70bda176ac4db06f1f66f96ae0693be1"))
}
