package flora

import (
	"encoding/json"
	"net/http"
	"sort"

	version "github.com/hashicorp/go-version"
)

const releasesIndexURL string = "https://releases.hashicorp.com/terraform/index.json"

func GetLatestVersion() (string, error) {
	type CheckResponse struct {
		CurrentVersion string `json:"current_version"`
	}

	checkResponse := CheckResponse{}

	r, err := http.Get(tfCheckpointURL)
	if err != nil {
		return "", err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&checkResponse)

	if err != nil {
		return "", err
	}

	return checkResponse.CurrentVersion, nil
}

func ListAllVersions() ([]*version.Version, error) {
	var versions []*version.Version
	versionsWrapper := struct {
		Versions map[string]interface{}
	}{}

	r, err := http.Get(releasesIndexURL)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&versionsWrapper)

	if err != nil {
		return nil, err
	}

	versions = make([]*version.Version, len(versionsWrapper.Versions))

	i := 0
	for ver := range versionsWrapper.Versions {
		versions[i], _ = version.NewVersion(ver)
		i++
	}

	sort.Sort(version.Collection(versions))

	return versions, nil
}
