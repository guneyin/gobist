package gobist

import (
	"github.com/imroc/req/v3"
	"go.nhat.io/cookiejar"
)

var (
	crumb string

	jar = cookiejar.NewPersistentJar(
		cookiejar.WithFilePath("cookies.json"),
		cookiejar.WithFilePerm(0755),
		cookiejar.WithAutoSync(true),
	)
)

func newClient() *req.Client {
	client := req.NewClient().
		SetCookieJar(jar).
		SetBaseURL(baseURL).
		SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Windows; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	setClientCommonHeaders(client)

	return client
}

func setClientCommonHeaders(c *req.Client) {
	headers := make(map[string]string)
	headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8"
	headers["Host"] = "query2.finance.yahoo.com"

	c.SetCommonHeaders(headers)
}

func getYahooCrumb() string {
	c := newClient()

	res, err := c.R().Get(crumbPath)
	if err != nil {
		return ""
	}

	return res.String()
}
