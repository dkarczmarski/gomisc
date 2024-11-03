package windsnsr_test

import (
	"context"
	"github.com/dkarczmarski/gomisc/windguru/env"
	"github.com/dkarczmarski/gomisc/windguru/windsnsr"
	"log"
	"os"
	"testing"
	"time"
)

func TestService_Get(t *testing.T) {
	if err := env.Load("./../.env"); err != nil {
		log.Fatal(err)
	}

	deviceID := os.Getenv("DEVICE_ID")
	token := os.Getenv("TOKEN")

	srv := windsnsr.Service{
		DeviceID: deviceID,
		Token:    token,
	}
	for i := 0; i < 6*30; i++ {
		resp, err := srv.Get(context.TODO())
		if err != nil {
			t.Error(err)
		}
		t.Logf("%+v", &resp)
		time.Sleep(10 * time.Second)
	}
}
