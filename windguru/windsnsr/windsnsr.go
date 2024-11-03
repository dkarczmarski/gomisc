// Package windsnsr provides method to get observation data from Tempest wind sensor device
package windsnsr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/guregu/null/v5"
	"io"
	"net/http"
)

type Service struct {
	DeviceID string
	Token    string
}

type apiResponse struct {
	Status apiStatus `json:"status"`
	Obs    [][]null.Float
}

type apiStatus struct {
	StatusCode    null.Int `json:"status_code"`
	StatusMessage string   `json:"status_message"`
}

const (
	//1 - Wind Lull (m/s)
	obsWindLullMs = 1

	//2 - Wind Avg (m/s)
	obsWindAvgMs = 2

	//3 - Wind Gust (m/s)
	obsWindGustMs = 3

	//4 - Wind Direction (degrees)
	obsWindDirectionDegrees = 4

	//6 - Pressure (MB)
	obsPressureMb = 6

	//7 - Air Temperature (C)
	obsAirTemperatureC = 7

	//8 - Relative Humidity (%)
	absRelativeHumidityPrc = 8
)

type Response struct {
	Status        int
	ErrorMessage  string
	WindMin       float64
	WindAvg       float64
	WindGust      float64
	WindDirection float64
	Pressure      float64
	Temperature   float64
	Humidity      float64
}

func (srv *Service) Get(ctx context.Context) (Response, error) {
	var response Response

	url := fmt.Sprintf("https://swd.weatherflow.com/swd/rest/observations/?device_id=%v&token=%v",
		srv.DeviceID, srv.Token,
	)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, err
	}
	req = req.WithContext(ctx)

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return response, err
	}

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("status not ok: %v", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	var rawResp apiResponse
	if err := json.Unmarshal(b, &rawResp); err != nil {
		return response, err
	}

	response.Status = int(rawResp.Status.StatusCode.ValueOrZero())
	response.ErrorMessage = rawResp.Status.StatusMessage

	if ok, v := rawResp.Status.StatusCode.Valid, rawResp.Status.StatusCode.ValueOrZero(); !ok || v != 0 {
		return response, fmt.Errorf("bad response status: %v", rawResp.Status.StatusCode)
	}

	if len(rawResp.Obs) != 1 {
		return response, errors.New("data error")
	}
	obs := rawResp.Obs[0]

	response.WindMin = obs[obsWindLullMs].ValueOrZero()
	response.WindAvg = obs[obsWindAvgMs].ValueOrZero()
	response.WindGust = obs[obsWindGustMs].ValueOrZero()
	response.WindDirection = obs[obsWindDirectionDegrees].ValueOrZero()
	response.Pressure = obs[obsPressureMb].ValueOrZero()
	response.Temperature = obs[obsAirTemperatureC].ValueOrZero()
	response.Humidity = obs[absRelativeHumidityPrc].ValueOrZero()

	return response, nil
}
