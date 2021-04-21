package handler

import (
	"sort"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keRin7/nexus-manager/nexusmanager"
	"github.com/sirupsen/logrus"
)

type TagsList struct {
	Data []string `json:"data"`
}

// Id is uniq, It contains one of tags
type Image struct {
	Id   string
	Tags string
	Size string
	Data string
	SHA  string
}

type Images struct {
	Images []Image
}
type DataStruct struct {
	Images
	Repos       nexusmanager.Repositories
	CurrentRepo string
	Username    string
}

// search in list the same SHA, if it found then union tags
func tagHasAliases(list Images, sha string, tag string) bool {
	for id, elem := range list.Images {
		if elem.SHA == sha {
			list.Images[id].Tags = list.Images[id].Tags + " <BR>" + tag
			return true
		}
	}
	return false
}

func (h *Handler) getRepos(c *gin.Context) {
	token, _ := h.GetToken(c)
	username := h.nexusmanager.Auth.GetUsername(token)
	logrus.Println("User: ", username, "get access to: ", h.nexusmanager.Config.Nexus_repo+c.Param("id"))

	var list Images
	tags := h.nexusmanager.ListTagsByImage(h.nexusmanager.Config.Nexus_repo + c.Param("id"))
	repo := h.nexusmanager.List()
	for _, v := range tags {
		data, sha := h.nexusmanager.GetDataAndSHAV1(h.nexusmanager.Config.Nexus_repo+c.Param("id"), v)
		if tagHasAliases(list, sha, v) {
			continue
		}
		size := h.nexusmanager.GetSize(h.nexusmanager.Config.Nexus_repo+c.Param("id"), v)
		list.Images = append(list.Images, Image{v, v, strconv.FormatInt(size/1024/1024, 10), data, sha})
	}

	sort.Slice(list.Images, func(i, j int) bool {
		return list.Images[i].Data < list.Images[j].Data
	})

	tmpl, _ := template.ParseFiles("template/repo.html")
	tmpl.Execute(c.Writer, &DataStruct{list, *repo, h.nexusmanager.Config.Nexus_repo + c.Param("id"), username})
}

type RepoTemplate struct {
	RepoName string
	Size     int64
}

var StructRepoTemplate struct {
	Data []RepoTemplate
	Time time.Time
}

func (h *Handler) getReposList(c *gin.Context) {

	if time.Since(StructRepoTemplate.Time) > 1*time.Minute {
		repos := h.nexusmanager.List()
		StructRepoTemplate.Data = nil
		for _, repo := range repos.Images {
			StructRepoTemplate.Data = append(StructRepoTemplate.Data, RepoTemplate{repo, h.nexusmanager.GetRepoSize(repo) / 1024 / 1024})
		}
		StructRepoTemplate.Time = time.Now()

		sort.Slice(StructRepoTemplate.Data, func(i, j int) bool {
			return StructRepoTemplate.Data[i].Size > StructRepoTemplate.Data[j].Size
		})
	}

	tmpl, _ := template.ParseFiles("template/index.html")
	tmpl.Execute(c.Writer, &StructRepoTemplate)
}
