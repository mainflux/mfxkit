package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pingEndpoint = "ping"
)

func (sdk mfSDK) Ping(secret string) (string, error) {
	req := pingReq{Secret: secret}
	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s", sdk.mfxkitURL, pingEndpoint)

	_, body, sdkerr := sdk.processRequest(http.MethodPost, url, "", string(CTJSON), data, http.StatusOK)
	if sdkerr != nil {
		return "", sdkerr
	}

	res := pingRes{}
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	return res.Greeting, nil
}
