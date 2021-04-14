package nexusmanager

import (
	"encoding/json"
	"fmt"

	"github.com/keRin7/nexus-manager/pkg/appcache"
	"github.com/keRin7/nexus-manager/pkg/auth"
	"github.com/keRin7/nexus-manager/pkg/ldapcli"
	"github.com/keRin7/nexus-manager/pkg/rest_client"
	"github.com/sirupsen/logrus"
)

type NexusManager struct {
	Config  *Config
	cache   *appcache.AppCache
	rest    *rest_client.Rest_client
	Ldapcli *ldapcli.LdapCli
	Auth    *auth.Auth
}

// specific header for nexus API
const ACCEPT_HEADER = "application/vnd.docker.distribution.manifest.v2+json"

// func List
type Repositories struct {
	Images []string `json:"repositories"`
}

// func ListTagsByImage
type ImageTags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// func ImageManifest / GetSize
type ImageManifest struct {
	SchemaVersion int64       `json:"schemaVersion"`
	MediaType     string      `json:"mediaType"`
	Config        LayerInfo   `json:"config"`
	Layers        []LayerInfo `json:"layers"`
}

// including in ImageManifest
type LayerInfo struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

// func GetDataV1
type ImageManifestV1 struct {
	SchemaVersion int64            `json:"schemaVersion"`
	History       []LayersHistory1 `json:"history"`
}

// func GetDataV1
type LayersHistory1 struct {
	V1Compatibility string `json:"v1Compatibility"`
}

// func GetDataV1
type LayersHistory2 struct {
	Created string `json:"created"`
	ID      string `json:"id"`
}

func New(config *Config) *NexusManager {
	return &NexusManager{
		Config:  config,
		cache:   appcache.NewCache(),
		rest:    rest_client.NewRestClient(),
		Ldapcli: ldapcli.New(config.Ldap),
		Auth:    &auth.Auth{},
	}
}

func (c *NexusManager) List() *Repositories {
	//c.Config.Nexus_repo
	headers := map[string]string{
		"Accept": ACCEPT_HEADER,
	}

	url := fmt.Sprintf("%s/repository/%s/v2/_catalog", c.Config.Nexus_url, c.Config.Nexus_repo)
	out, _ := c.rest.DoGet(url, headers, c.Config.Nexus_username, c.Config.Nexus_password)

	var repositories Repositories
	err := json.Unmarshal(out, &repositories)
	if err != nil {
		logrus.Fatal(err)
	}
	//fmt.Printf("%v", repositories)
	return &repositories
}

func (c *NexusManager) ListTagsByImage(image string) []string {
	headers := map[string]string{
		"Accept": ACCEPT_HEADER,
	}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/tags/list", c.Config.Nexus_url, c.Config.Nexus_repo, image)

	out, _ := c.rest.DoGet(url, headers, c.Config.Nexus_username, c.Config.Nexus_password)

	var imageTags ImageTags
	err := json.Unmarshal(out, &imageTags)
	if err != nil {
		logrus.Fatal(err)
	}
	return imageTags.Tags

}

func (c *NexusManager) ImageManifest(image string, tag string) {
	headers := map[string]string{
		"Accept": ACCEPT_HEADER,
	}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", c.Config.Nexus_url, c.Config.Nexus_repo, image, tag)
	out, _ := c.rest.DoGet(url, headers, c.Config.Nexus_username, c.Config.Nexus_password)

	var imageManifest ImageManifest
	err := json.Unmarshal(out, &imageManifest)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("%#v", imageManifest)

}

func (c *NexusManager) GetSize(image string, tag string) int64 {
	headers := map[string]string{
		"Accept": ACCEPT_HEADER,
	}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", c.Config.Nexus_url, c.Config.Nexus_repo, image, tag)
	out, _ := c.rest.DoGet(url, headers, c.Config.Nexus_username, c.Config.Nexus_password)

	var imageManifest ImageManifest
	err := json.Unmarshal(out, &imageManifest)
	if err != nil {
		logrus.Fatal(err)
	}
	var imageSize int64
	for _, v := range imageManifest.Layers {
		imageSize += v.Size
	}
	return imageSize

}

func (c *NexusManager) GetImageSHA(image string, tag string) map[string][]string {
	headers := map[string]string{
		"Accept": ACCEPT_HEADER,
	}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", c.Config.Nexus_url, c.Config.Nexus_repo, image, tag)
	out, h := c.rest.DoGet(url, headers, c.Config.Nexus_username, c.Config.Nexus_password)

	var imageManifest ImageManifest
	err := json.Unmarshal(out, &imageManifest)
	if err != nil {
		logrus.Fatal(err)
	}
	return h

}

func (c *NexusManager) GetDataV1(image string, tag string) string {
	if v, ok := c.cache.Get(image + tag); ok {
		return v
	}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", c.Config.Nexus_url, c.Config.Nexus_repo, image, tag)
	out, _ := c.rest.DoGet(url, nil, c.Config.Nexus_username, c.Config.Nexus_password)

	var imageManifestV1 ImageManifestV1
	err := json.Unmarshal(out, &imageManifestV1)
	if err != nil {
		logrus.Fatal(err)
	}

	var layersHistory2 LayersHistory2
	err = json.Unmarshal([]byte(imageManifestV1.History[0].V1Compatibility), &layersHistory2)
	if err != nil {
		logrus.Fatal(err)
	}
	c.cache.Set(image+tag, layersHistory2.Created)
	return layersHistory2.Created

}
