package encrypt

import "github.com/chanbakjsd/gotrix/api"

type Endpoints struct {
	api.Endpoints
}

func (e Endpoints) Keys() string        { return e.Base() + "/keys" }
func (e Endpoints) KeysChanges() string { return e.Base() + "/changes" }
func (e Endpoints) KeysClaim() string   { return e.Base() + "/claim" }
func (e Endpoints) KeysQuery() string   { return e.Base() + "/query" }
func (e Endpoints) KeysUpload() string  { return e.Base() + "/upload" }
