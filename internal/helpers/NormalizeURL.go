package helpers

import (
	urlUtil "net/url"
	"strings"
)

func NormalizeURL(rawURL string) (string, error) {
	url, err := urlUtil.Parse(rawURL)
	if err != nil {
		return "", err
	}

	url.Scheme = strings.ToLower(url.Scheme)
	url.Host = strings.ToLower(url.Host)

	if url.Scheme == "http" && strings.HasSuffix(url.Host, ":80") ||
		url.Scheme == "https" && strings.HasSuffix(url.Host, ":443") {

		url.Host = strings.Split(url.Host, ":")[0]
	}

	if len(url.Path) > 0 && url.Path[len(url.Path)-1] == '/' {
		url.Path += "/"
	}

	if url.RawQuery != "" {
		values, _ := urlUtil.ParseQuery(url.RawQuery)
		url.RawQuery = values.Encode()
	}

	return url.String(), nil
}
