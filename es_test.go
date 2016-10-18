package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestElasticSearch(t *testing.T) {
	esURL := "http://127.0.0.1:9200"
	indexName := "k8s_cd_test"
	err := cleanupIndex(esURL, indexName)
	if err != nil {
		t.Fatal(err)
	}
	payload := `
	{"index":{"_id":1,"_type":"record","_index":"` + indexName + `"}}
	{"email": "tobstarr@gmail.com"}
	`
	rsp, err := http.Post(esURL+"/_bulk?refresh=true", "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()
	if rsp.Status[0] != '2' {
		b, _ := ioutil.ReadAll(rsp.Body)
		t.Fatalf("got status %s but expected 2x. body=%s", rsp.Status, string(b))
	}

	rsp, err = http.Get(esURL + "/" + indexName + "/_search")
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if rsp.Status[0] != '2' {
		t.Fatalf("got status %s but expected 2x. body=%s", rsp.Status, string(b))
	}
	_ = b
	var res struct {
		Hits struct {
			Total int `json:"total"`
			Hits  []struct {
				Type   string                 `json:"_type"`
				ID     string                 `json:"_id"`
				Index  string                 `json:"_index"`
				Source map[string]interface{} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct{ Has, Want interface{} }{
		{res.Hits.Total, 1},
		{len(res.Hits.Hits), 1},
		{res.Hits.Hits[0].Type, "record"},
		{res.Hits.Hits[0].Source["email"], "tobstarr@gmail.com"},
	}
	for i, tc := range tests {
		if tc.Has != tc.Want {
			t.Errorf("%d: want=%#v has=%#v", i+1, tc.Want, tc.Has)
		}
	}
}

func cleanupIndex(url, index string) error {
	req, err := http.NewRequest("DELETE", url+"/"+index, nil)
	if err != nil {
		return err
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.Status[0] != '2' && rsp.StatusCode != 404 {
		b, _ := ioutil.ReadAll(rsp.Body)
		return fmt.Errorf("got status %s but expected 2x. body=%s", rsp.Status, string(b))
	}
	return nil

}
