package airtime

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type InformationEndpoint struct {
	Req          InformationPrincipalRequest
	Res          InformationPrincipalResponse
	resellerID   string
	resellerType string
	endpoint     string
}

func (info *InformationEndpoint) SetID(id string) *InformationEndpoint {
	info.resellerID = id
	return info
}

func (info *InformationEndpoint) SetType(resellertype string) *InformationEndpoint {
	info.resellerType = resellertype
	return info
}

func (info *InformationEndpoint) SetEdpoint(endpoint string) *InformationEndpoint {
	info.endpoint = endpoint
	return info
}
func (info *InformationEndpoint) GetEndpoint() string {
	return info.endpoint
}

func (info *InformationEndpoint) PrepareRequest() ([]byte, error) {
	timeout := 5 * time.Microsecond
	user := User{ID: "RES0000059747", Type: "RESELLERUSER", RequestType: "DMark"}
	user.XMLName = xml.Name{Local: "initiatorPrincipalId"}
	ref := randomString(8)
	context := Context{"WEBSERVICE", "DMARK", timeout, user, "DM@rk321", ref, "dmark topup subscriber account"}
	acc_owner := User{ID: info.resellerID, Type: info.resellerType}
	acc_owner.XMLName = xml.Name{Local: "principalId"}

	infoRequest := new(InformationPrincipalRequest)
	infoRequest.Context = context
	infoRequest.Reseller = acc_owner

	body := toXML(infoRequest)

	return []byte(body), nil
}

func (info *InformationEndpoint) PrepareResponse(r http.Response) (interface{}, error) {
	var resp InformationPrincipalResponse = InformationPrincipalResponse{}
	//footer := "</soap:Body></soap:Envelope>"
	//lastindex := strings.Index(b, footer)
	//b = b[80:lastindex] // Hard coded xml header length. How did I find this? never mind

	var b bytes.Buffer
	if _, err := io.Copy(&b, r.Body); err != nil {
		return nil, err
	}

	if err := xml.Unmarshal([]byte(b.String()), &resp); err != nil {
		//c.log("Error:", err)
		return nil, err
	}
	if resp.ResultCode != 0 {
		var _ ErrorResponse = (*InformationPrincipalResponse)(nil)
		return &resp, nil
	}
	return &resp, nil
}
