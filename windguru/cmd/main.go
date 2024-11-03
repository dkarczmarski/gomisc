package main

import (
	"context"
	"github.com/dkarczmarski/gomisc/windguru/env"
	"github.com/dkarczmarski/gomisc/windguru/windguru"
	"github.com/dkarczmarski/gomisc/windguru/windsnsr"
	"log"
	"time"
)

func MsToKnot(v float64) float64 {
	return v * 1.9438
}

func main() {
	if err := env.Load(".env"); err != nil {
		log.Fatal(err)
	}

	windguruUID := env.Getenv("WINDGURU_UID")
	windguruPass := env.Getenv("WINDGURU_PASS")
	deviceID := env.Getenv("DEVICE_ID")
	token := env.Getenv("TOKEN")

	srv := windsnsr.Service{
		DeviceID: deviceID,
		Token:    token,
	}

	for {
		wsResp, err := func() (windsnsr.Response, error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			return srv.Get(ctx)
		}()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("%+v", &wsResp)

		wgResp, err := func() (string, error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			c := windguru.New(windguruUID, windguruPass)
			return c.Send(
				ctx,
				MsToKnot(wsResp.WindMin),
				MsToKnot(wsResp.WindAvg),
				MsToKnot(wsResp.WindGust),
				wsResp.WindDirection,
				wsResp.Temperature,
				wsResp.Pressure,
				wsResp.Humidity,
			)
		}()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(wgResp)

		time.Sleep(30 * time.Second)
	}
}
