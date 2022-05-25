package corehandlers

import (
	"encoding/json"
	"strings"

	"github.com/QN-zhangzhuo/go-sdk/qiniu/defs"
	"github.com/QN-zhangzhuo/go-sdk/qiniu/qerr"
	"github.com/QN-zhangzhuo/go-sdk/qiniu/request"
)

// UnmarshalHandler 反序列化http.Body到相应的结构体中
var UnmarshalHandler = request.NamedHandler{
	Name: "UnmarshalHandler",
	Fn: func(r *request.Request) {
		if r.DataFilled() {
			contentType := r.HTTPResponse.Header.Get("Content-Type")
			splits := strings.Split(contentType, ";")
			if len(splits) > 0 {
				contentType = splits[0]
			}
			switch contentType {
			case defs.CONTENT_TYPE_JSON:
				err := json.NewDecoder(r.HTTPResponse.Body).Decode(r.Data)
				if err != nil {
					r.Error = qerr.New(qerr.ErrCodeDeserialization, "failed to decode data with content-type: "+contentType, err)
					return
				}
			}
		}
	},
}
