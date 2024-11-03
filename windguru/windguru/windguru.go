// Package windguru provides method to send data to windguru station graph
package windguru

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Service struct {
	uid, pass string
}

func New(uid, pass string) *Service {
	return &Service{
		uid:  uid,
		pass: pass,
	}
}

func (srv *Service) Send(
	ctx context.Context, windMin, windAvg, windMax, direction, temperature, pressure, rh float64,
) (string, error) {
	var salt = fmt.Sprintf("%v", time.Now().UnixMilli())

	var authStr = fmt.Sprintf("%v%v%v", salt, srv.uid, srv.pass)
	hashf := md5.New()
	hashf.Write([]byte(authStr))
	sumBytes := hashf.Sum(nil)
	sum := hex.EncodeToString(sumBytes)

	u := url.Values{}
	u.Add("uid", srv.uid)
	u.Add("salt", salt)
	u.Add("hash", sum)
	u.Add("wind_min", fmt.Sprintf("%v", windMin))
	u.Add("wind_avg", fmt.Sprintf("%v", windAvg))
	u.Add("wind_max", fmt.Sprintf("%v", windMax))
	u.Add("wind_direction", fmt.Sprintf("%v", direction))
	u.Add("temperature", fmt.Sprintf("%v", temperature))
	u.Add("mslp", fmt.Sprintf("%v", pressure))
	u.Add("rh", fmt.Sprintf("%v", rh))

	reqURL := "http://www.windguru.cz/upload/api.php?" + u.Encode()
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return "", err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
