package gobist

import (
	"fmt"
	"github.com/imroc/req/v3"
	"go.nhat.io/cookiejar"
)

func newClient() *req.Client {
	jar := cookiejar.NewPersistentJar(
		cookiejar.WithFilePath("cookies.json"),
		cookiejar.WithFilePerm(0755),
		cookiejar.WithAutoSync(true),
	)

	headers := make(map[string]string)
	headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8"
	headers["Host"] = "query2.finance.yahoo.com"

	client := req.NewClient().
		SetCookieJar(jar).
		SetBaseURL(baseURL).
		SetCommonHeaders(headers).
		SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Windows; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	return client
}

func setYahooCrumb() string {
	res, err := req.C().R().Get(crumbPath)
	if err != nil {
		fmt.Printf("crumb error: %v\n", err)
	}

	crumb := res.String()

	fmt.Printf("crumb has been set: %v\n", crumb)
	return crumb
}
