package windguru_test

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dkarczmarski/gomisc/windguru/env"
	"github.com/dkarczmarski/gomisc/windguru/windguru"
	"log"
	"os"
	"testing"
	"time"
)

func TestMd5(t *testing.T) {
	s := "abc12301234"

	hashf := md5.New()
	hashf.Write([]byte(s))
	sum1 := hashf.Sum(nil)
	sum2 := hex.EncodeToString(sum1)

	fmt.Printf("MD5 hash: %x\n", sum1)
	fmt.Printf("MD5 hash: %s\n", sum2)
}

func TestService_Send(t *testing.T) {
	if err := env.Load("./../.env"); err != nil {
		log.Fatal(err)
	}

	windguruUID := os.Getenv("WINDGURU_UID")
	windguruPass := os.Getenv("WINDGURU_PASS")

	for i := 0; i < 5; i++ {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
			defer cancel()

			srv := windguru.New(windguruUID, windguruPass)
			resp, err := srv.Send(
				ctx,
				10.3,
				11.1+0.4*float64(i),
				13.2+0.4*float64(i),
				10,
				28,
				1001,
				80,
			)
			if err != nil {
				t.Error(err)
			}
			t.Log(resp)
		}()

		time.Sleep(10 * time.Second)
	}
}
