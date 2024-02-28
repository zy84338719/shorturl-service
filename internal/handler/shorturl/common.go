package shorturl

import (
	"github.com/imroc/req/v3"
	"net/http"
)

func checkUrlEffective(url string) bool {
	code := req.C().Options(url).Do().StatusCode
	if code >= http.StatusBadRequest {
		return true
	}
	return false
}
