package nexusmanager

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/keRin7/nexus-manager/pkg/appcache"
	"github.com/keRin7/nexus-manager/pkg/auth"
	"github.com/keRin7/nexus-manager/pkg/ldapclient"
	"github.com/keRin7/nexus-manager/pkg/rest_client"
	"github.com/sirupsen/logrus"
)

type NexusManager struct {
	Config     *Config
	Cache      *appcache.AppCache
	rest       *rest_client.Rest_client
	Ldapclient *ldapclient.LdapClient
	Auth       *auth.Auth
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

// func GetDataV1
type LayersHistory2 struct {
	Created string `json:"created"`
	ID      string `json:"id"`
}

// func GetImageLayers
type LayersHistory3 struct {
	Created string `json:"created"`
	ID      string `json:"id"`
	Cmd     struct {
		Container_config string `json:"container_config"`
	} `json:"Cmd,omitempty"`
	Config struct {
		User       string   `json:"User,omitempty"`
		WorkingDir string   `json:"WorkingDir,omitempty"`
		Entrypoint []string `json:"Entrypoint,omitempty"`
		Env        []string `json:"Env,omitempty"`
	} `json:"config,omitempty"`
}

func New(config *Config) *NexusManager {
	return &NexusManager{
		Config:     config,
		Cache:      appcache.NewCache(),
		rest:       rest_client.NewRestClient(),
		Ldapclient: ldapclient.New(config.Ldap),
		Auth:       &auth.Auth{Admin_users: config.Admin_users},
	}
}

/*
List()
nexus.com/repository/ROOT-REPO/v2/_catalog    			-get list repos in ROOT-REPO
ListTagsByImage()
nexus.com/repository/ROOT-REPO/v2/IMAGE-NAME/tags/list	-get list tags if i know IMAGE-NAME
*/

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

func (c *NexusManager) GetSize(image string, tag string) int64 {
	imageManifest, _ := c.GetManifestV2(image, tag)
	var imageSize int64
	for _, v := range imageManifest.Layers {
		imageSize += v.Size
	}
	return imageSize

}

// Get top layer data created and layersSHA
func (c *NexusManager) GetDataAndSHAV1(image string, tag string) (data string, sha string) {
	if data, sha, ok := c.Cache.GetData(image + tag); ok {
		return data, sha
	}

	imageManifest := c.GetManifestV1(image, tag)

	var layersHistory2 LayersHistory2
	err := json.Unmarshal([]byte(imageManifest.History[0].V1Compatibility), &layersHistory2)
	if err != nil {
		logrus.Fatal(err)
	}
	// Join fslayers
	layersSHA := ""
	for _, l := range imageManifest.FsLayers {
		layersSHA = layersSHA + " " + l.BlobSum
	}

	c.Cache.SetData(image+tag, layersSHA, layersHistory2.Created)
	return layersHistory2.Created, layersSHA
}

type ImageLayerTemplate struct {
	Date     string
	Cmd      string
	LayerSHA string
	Size     int64
}

// Get layer-docker commands and layer sizes
func (c *NexusManager) GetLayersInfoV1(image string, tag string) (Layers []ImageLayerTemplate, ImagePropertyUser string, ImagePropertyWorkingDir string, ImagePropertyEntrypoint []string, ImagePropertyEnv []string) {
	var ReturnLayers []ImageLayerTemplate

	imageManifestV1 := c.GetManifestV1(image, tag)

	var Layer LayersHistory3
	DateLength := 16 //Show only 16 symbols in template [YYYY-MM-DDTHH:MM]:SS
	for _, el := range imageManifestV1.History {
		err := json.Unmarshal([]byte(el.V1Compatibility), &Layer)
		if err != nil {
			logrus.Fatal(err)
		}

		replaceBashLogicalOperation := strings.ReplaceAll(Layer.Cmd.Container_config, "&&", "&&<BR>")
		replaceBashCommandSeparator := strings.ReplaceAll(replaceBashLogicalOperation, ";", ";<BR>")
		ReturnLayers = append(ReturnLayers, ImageLayerTemplate{Layer.Created[:DateLength], replaceBashCommandSeparator, "", 0})
	}

	imageManifestV2, _ := c.GetManifestV2(image, tag)

	LayersSHAandSizes := make(map[string]int64)
	for _, el := range imageManifestV2.Layers {
		LayersSHAandSizes[el.Digest] = el.Size
	}

	// Lets insert sizes in our struct
	for i, _ := range ReturnLayers {
		ReturnLayers[i].LayerSHA = imageManifestV1.FsLayers[i].BlobSum
		if val, ok := LayersSHAandSizes[imageManifestV1.FsLayers[i].BlobSum]; ok {
			ReturnLayers[i].Size = val / 1024 / 1024
		}
	}

	return ReturnLayers, Layer.Config.User, Layer.Config.WorkingDir, Layer.Config.Entrypoint, Layer.Config.Env
}

func (c *NexusManager) GetImageSHA(image string, tag string) (string, error) {
	_, h := c.GetManifestV2(image, tag)
	return strings.Join(h["Docker-Content-Digest"], ""), nil
}

func (c *NexusManager) DeleteImageByTag(imageNameWithRepoPath string, tag string) error {
	sha, err := c.GetImageSHA(imageNameWithRepoPath, tag)

	if err != nil {
		logrus.Printf("Error delete image %s", err.Error())
		return err
	}

	headers := map[string]string{
		"Accept": ACCEPT_HEADER,
	}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", c.Config.Nexus_url, c.Config.Nexus_repo, imageNameWithRepoPath, sha)
	out, err := c.rest.DoDelete(url, headers, c.Config.Nexus_username, c.Config.Nexus_password)
	//logrus.Println(url)
	if err != nil {
		logrus.Printf("Error delete image %s", err.Error())
		return err
	}

	fmt.Printf("%v", out)
	return nil
}

func (c *NexusManager) GetRepoSize(image string) int64 {
	layers := make(map[string]int64)
	tags := c.ListTagsByImage(image)
	var imageManifest ImageManifestV2

	for _, tag := range tags {
		imageManifest, _ = c.GetManifestV2(image, tag)

		for _, v := range imageManifest.Layers {
			layers[v.Digest] = v.Size
		}

	}

	var projectSize int64

	for _, v := range layers {
		projectSize = projectSize + v
	}

	return projectSize

}
