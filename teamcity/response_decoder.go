package teamcity

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type responseDecoder struct{}

func (responseDecoder) Decode(resp *http.Response, v interface{}) error {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	b, err:= json.Marshal(bodyBytes)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, v)
	return err
}