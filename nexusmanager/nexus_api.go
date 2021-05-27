package nexusmanager

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

type ImageManifestV2 struct {
	SchemaVersion int64  `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int64  `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []struct {
		MediaType string `json:"mediaType"`
		Size      int64  `json:"size"`
		Digest    string `json:"digest"`
	} `json:"layers"`
}

type ImageManifestV1 struct {
	SchemaVersion int64 `json:"schemaVersion"`
	FsLayers      []struct {
		BlobSum string `json:"blobSum"`
	} `json:"fsLayers"`
	History []struct {
		V1Compatibility string `json:"v1Compatibility"`
	} `json:"history"`
}

func (c *NexusManager) GetManifestV1(image string, tag string) ImageManifestV1 {

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", c.Config.Nexus_url, c.Config.Nexus_repo, image, tag)
	out, _ := c.rest.DoGet(url, nil, c.Config.Nexus_username, c.Config.Nexus_password)

	var imageManifest ImageManifestV1
	err := json.Unmarshal(out, &imageManifest)
	if err != nil {
		logrus.Fatal(err)
	}
	return imageManifest
}

func (c *NexusManager) GetManifestV2(image string, tag string) (ImageManifestV2, map[string][]string) {

	headers := map[string]string{
		"Accept": ACCEPT_HEADER,
	}
	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", c.Config.Nexus_url, c.Config.Nexus_repo, image, tag)
	out, h := c.rest.DoGet(url, headers, c.Config.Nexus_username, c.Config.Nexus_password)

	var imageManifest ImageManifestV2
	err := json.Unmarshal(out, &imageManifest)
	if err != nil {
		logrus.Fatal(err)
	}
	return imageManifest, h
}
