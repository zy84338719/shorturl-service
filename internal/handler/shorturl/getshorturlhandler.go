package shorturl

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shorturl/internal/logic/shorturl"
	"shorturl/internal/svc"
	"shorturl/internal/types"
)

func GetShorturlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetShortRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := shorturl.NewGetShorturlLogic(r.Context(), svcCtx)
		resp, err := l.GetShorturl(&req)
		if err != nil {
			l.Errorf("GetShorturl error: %v", err)
		}
		http.Redirect(w, r, resp, http.StatusFound)
	}
}

// Redirect replies to the request with a redirect to url,
// which may be a path relative to the request path.
//
// The provided code should be in the 3xx range and is usually
// StatusMovedPermanently, StatusFound or StatusSeeOther.
//
// If the Content-Type header has not been set, Redirect sets it
// to "text/html; charset=utf-8" and writes a small HTML body.
// Setting the Content-Type header to any value, including nil,
// disables that behavior.

// 翻译下上面的注释
// Redirect回复请求，重定向到url，该url可能是相对于请求路径的路径。
// 提供的code应该在3xx范围内，通常是StatusMovedPermanently，StatusFound或StatusSeeOther。
// 如果未设置Content-Type标头，则Redirect将其设置为"text/html; charset=utf-8"并写入一个小的HTML正文。
// 将Content-Type标头设置为任何值（包括nil）都会禁用该行为。
//
// Redirect函数的作用是重定向到指定的url，code指定了重定向的类型，如下：
// 301：永久重定向，表示请求的资源已经被分配了新的URI，以后应使用资源现在所指的URI。
// 302：临时重定向，表示请求的资源临时被分配了新的URI，希望本次能使用新的URI访问。
// 303：与302有相同的功能，但是期望客户端在请求新的URI时，使用GET方法请求。
// 307：与302有相同的功能，但是期望客户端在请求新的URI时，保持请求方法不变向新的URI发出请求。
// 308：永久重定向，表示请求的资源已经被分配了新的URI，以后应使用资源现在所指的URI。
// 以上的code都是3xx，表示重定向，但是有区别，301和308表示永久重定向，302、303、307表示临时重定向。
// 301和302是HTTP1.0的标准，303和307是HTTP1.1的标准，而308是HTTP3.0的标准。
// 301和308的区别是，301是HTTP1.0的标准，308是HTTP3.0的标准，所以308是301的升级版。
