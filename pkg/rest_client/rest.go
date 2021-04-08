package rest_client

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Rest_client struct {
	httpClient *http.Client
}

func NewRestClient() *Rest_client {
	return &Rest_client{
		httpClient: &http.Client{},
	}
}

func (r *Rest_client) DoGet(path string, headers map[string]string, username string, password string) ([]byte, map[string][]string) {

	//client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, path, nil)
	req.SetBasicAuth(username, password)
	for k, v := range headers {
		//fmt.Printf("%s %s", k, v)
		req.Header.Add(k, v)
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		log.Printf("DoGet: %s", path)
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if (resp.Status != "200 OK") && (resp.Status != "201 Created") {
		log.Println("GET:", resp.Status, " path: ", path)
	} else {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("DoGet:ReadAll: %s", path)
			log.Fatal(err)
		}
		return body, resp.Header
	}
	return nil, nil
	/* for true {

		bs := make([]byte, 1014)
		n, err := resp.Body.Read(bs)
		fmt.Println(string(bs[:n]))

		if n == 0 || err != nil {
			break
		}
	} */
}

func (r *Rest_client) DoPost(path string, body []byte, headers map[string]string) []byte {
	p := bytes.NewReader(body)
	/* resp, err := http.Post(r.config.url+path, "application/json", p)
	if err != nil {
		log.Printf("DoPost: %s", path)
		log.Fatal(err)
	} */

	//fmt.Println(r.config.url + path)
	req, err := http.NewRequest(http.MethodPost, path, p)
	if err != nil {
		log.Printf("DoPost: %s", path)
		log.Fatal(err)
	}

	for k, v := range headers {
		//fmt.Printf("%s %s", k, v)
		req.Header.Add(k, v)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		log.Printf("DoPost: %s", path)
		log.Fatal(err)
	}
	defer resp.Body.Close()
	//fmt.Println(resp.Status)

	if (resp.Status != "200 OK") && (resp.Status != "201 Created") {
		log.Println("POST:", resp.Status, " path: ", path)
	} else {

		bodyResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("DoPost:ReadAll: %s", path)
			log.Fatal(err)
		}
		return bodyResp
	}
	return nil

	//fmt.Println(string(resp.Status))

}

func (r *Rest_client) DoPatch(path string, body []byte, headers map[string]string) []byte {
	p := bytes.NewReader(body)
	//resp, err := http.

	req, err := http.NewRequest(http.MethodPatch, path, p)

	for k, v := range headers {
		//fmt.Printf("%s %s", k, v)
		req.Header.Add(k, v)
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		log.Printf("DoPatch: %s", path)
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if (resp.Status != "200 OK") && (resp.Status != "201 Created") {
		log.Println("PATCH:", resp.Status, " path: ", path)
	} else {

		bodyResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("DoPatch:ReadAll: %s", path)
			log.Fatal(err)
		}
		return bodyResp
	}
	return nil
	//fmt.Println(string(body))

}

/* data := []byte(`{"foo":"bar"}`)
r := bytes.NewReader(data)
resp, err := http.Post("http://example.com/upload", "application/json", r)
if err != nil {
    return err
} */

func (r *Rest_client) DoPut(path string, body io.Reader, headers map[string]string) []byte {
	//p := bytes.NewReader(body)
	//resp, err := http.

	req, err := http.NewRequest(http.MethodPut, path, body)

	for k, v := range headers {
		//fmt.Printf("%s %s", k, v)
		req.Header.Add(k, v)
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		log.Printf("DoPut: %s", path)
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if (resp.Status != "200 OK") && (resp.Status != "201 Created") {
		log.Println("PATCH:", resp.Status, " path: ", path)
	} else {

		bodyResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("DoPut:ReadAll: %s", path)
			log.Fatal(err)
		}
		return bodyResp
	}
	//fmt.Println(string(body))
	return nil

}
