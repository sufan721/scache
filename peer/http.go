package peer

import (
	"fmt"
	"io"
	"net/http"
)

type HttpGetter struct {
	BaseUrl string
}

func (g *HttpGetter) Get(group string, key string) (string, error) {

	url := fmt.Sprintf(
		"%s%s/%s",
		g.BaseUrl,
		group,
		key,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
