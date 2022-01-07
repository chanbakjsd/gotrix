package encrypt

import "github.com/chanbakjsd/gotrix/api"

var (
	EndpointKeys        = api.EndpointBase + "/keys"
	EndpointKeysChanges = EndpointKeys + "/changes"
	EndpointKeysClaim   = EndpointKeys + "/claim"
	EndpointKeysQuery   = EndpointKeys + "/query"
	EndpointKeysUpload  = EndpointKeys + "/upload"
)
