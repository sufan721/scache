package peer

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpGetter struct {
	baseUrl string
}

func (g *HttpGetter) Get(key string) (string, error) {
	url := fmt.Sprintf("%s?key=%s", g.baseUrl, key)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
