package wrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func MakeRequest[Req any, Res any](ctx context.Context, url string, req Req) (Res, error) {
	var res Res

	rawJSON, err := json.Marshal(req)
	if err != nil {
		return res, fmt.Errorf("marshaling json: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(rawJSON))
	if err != nil {
		return res, fmt.Errorf("creating http request: %w", err)
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return res, fmt.Errorf("MakeRequesting http: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return res, fmt.Errorf("wrong status code %d", httpResponse.StatusCode)
	}

	err = json.NewDecoder(httpResponse.Body).Decode(&res)
	if err != nil {
		return res, fmt.Errorf("decoding json: %w", err)
	}

	return res, nil
}
