package handler

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TagsList struct {
	Data []string `json:"data"`
}

type SizeList struct {
	id   string
	size string
	data string
}

func (h *Handler) getRoot(c *gin.Context) {
	//var view map[string]interface{}
	var list []SizeList
	var out string
	//fmt.Println("%v", )
	start := time.Now()
	tags := h.nexusmanager.ListTagsByImage("coolrocket/" + c.Param("id"))

	for _, v := range tags {
		//headers := h.nexusmanager.GetImageSHA("coolrocket/"+c.Param("id"), v)

		//	view[strconv.Itoa(k)] = v
		data := h.nexusmanager.GetDataV1("coolrocket/"+c.Param("id"), v)
		//data := "123"
		size := h.nexusmanager.GetSize("coolrocket/"+c.Param("id"), v)
		//size := 1024 * 1024 * 1024
		list = append(list, SizeList{v, strconv.FormatInt(size/1024/1024, 10), data})
		//out = out + "id: " + v + " size: " + strconv.FormatInt(size/1024/1024, 10) + "Mb date-created-last-layer:" + data + "<BR>"
		//	fmt.Println(v, size/1024/1024)
		//		fmt.Println(v, headers["Last-Modified"], headers["Date"])
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].data < list[j].data
	})

	for _, v := range list {
		out = out + "id: " + v.id + " size: " + v.size + "Mb_________ " + v.data + "<BR>"
	}

	//c.JSON(http.StatusOK, tags)
	//	c.JSON(http.StatusOK, TagsList{
	//		Data: tags,
	//	})
	//var s1 string
	//s1 = "123"
	//c.Header("")
	//c.String(http.StatusOK, "<html><br>OK<br>OK</html>")
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(out))
	//c.HTML(http.StatusOK, "<h1>%s</h1><div>%s</div>", s1)
	//c.
}
