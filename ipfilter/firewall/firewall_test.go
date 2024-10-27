package firewall_test

import (
	"github.com/dkarczmarski/gomisc/ipfilter/firewall"
	"reflect"
	"testing"
)

func TestService_AddDeleteIP(t *testing.T) {
	for _, tt := range []struct {
		name         string
		initBefore   func(service *firewall.Service, fixedTime *firewall.FixedTime)
		testFunc     func(service *firewall.Service, fixedTime *firewall.FixedTime) error
		expectedErr  func(err error) bool
		expectedList []firewall.IPEntry
	}{
		{
			name: "add incorrect ip",
			testFunc: func(service *firewall.Service, fixedTime *firewall.FixedTime) error {
				return service.AddIP("1.2.3,,4")
			},
			expectedErr: func(err error) bool {
				return err != nil
			},
		},
		{
			name: "add first ip",
			testFunc: func(service *firewall.Service, fixedTime *firewall.FixedTime) error {
				fixedTime.SetDateTime("2001-01-01 10:00:00")
				return service.AddIP("1.2.3.4")
			},
			expectedErr: func(err error) bool {
				return err == nil
			},
			expectedList: []firewall.IPEntry{
				{
					IP:        "1.2.3.4",
					CreatedAt: firewall.MustParseDateTime("2001-01-01 10:00:00"),
					UpdatedAt: firewall.MustParseDateTime("2001-01-01 10:00:00"),
				},
			},
		},
		{
			name: "add another ip",
			initBefore: func(service *firewall.Service, fixedTime *firewall.FixedTime) {
				fixedTime.SetDateTime("2001-01-01 10:00:00")
				_ = service.AddIP("1.2.3.4")
			},
			testFunc: func(service *firewall.Service, fixedTime *firewall.FixedTime) error {
				fixedTime.SetDateTime("2001-01-01 10:01:00")
				return service.AddIP("2.2.8.8")
			},
			expectedErr: func(err error) bool {
				return err == nil
			},
			expectedList: []firewall.IPEntry{
				{
					IP:        "1.2.3.4",
					CreatedAt: firewall.MustParseDateTime("2001-01-01 10:00:00"),
					UpdatedAt: firewall.MustParseDateTime("2001-01-01 10:00:00"),
				},
				{
					IP:        "2.2.8.8",
					CreatedAt: firewall.MustParseDateTime("2001-01-01 10:01:00"),
					UpdatedAt: firewall.MustParseDateTime("2001-01-01 10:01:00"),
				},
			},
		},
		{
			name: "add the same ip",
			initBefore: func(service *firewall.Service, fixedTime *firewall.FixedTime) {
				fixedTime.SetDateTime("2001-01-01 10:00:00")
				_ = service.AddIP("1.2.3.4")
			},
			testFunc: func(service *firewall.Service, fixedTime *firewall.FixedTime) error {
				fixedTime.SetDateTime("2001-01-01 10:01:00")
				return service.AddIP("1.2.3.4")
			},
			expectedErr: func(err error) bool {
				return err == nil
			},
			expectedList: []firewall.IPEntry{
				{
					IP:        "1.2.3.4",
					CreatedAt: firewall.MustParseDateTime("2001-01-01 10:00:00"),
					UpdatedAt: firewall.MustParseDateTime("2001-01-01 10:01:00"),
				},
			},
		},
		{
			name: "delete incorrect ip",
			testFunc: func(service *firewall.Service, fixedTime *firewall.FixedTime) error {
				return service.DeleteIP("1.2.3,,4")
			},
			expectedErr: func(err error) bool {
				return err != nil
			},
		},
		{
			name: "delete nonexistent ip",
			testFunc: func(service *firewall.Service, fixedTime *firewall.FixedTime) error {
				return service.DeleteIP("1.2.3.4")
			},
			expectedErr: func(err error) bool {
				return err != nil
			},
		},
		{
			name: "add existed ip",
			initBefore: func(service *firewall.Service, fixedTime *firewall.FixedTime) {
				fixedTime.SetDateTime("2001-01-01 10:00:00")
				_ = service.AddIP("1.2.3.4")

				fixedTime.SetDateTime("2001-01-01 10:01:00")
				_ = service.AddIP("2.2.8.8")
			},
			testFunc: func(service *firewall.Service, fixedTime *firewall.FixedTime) error {
				return service.DeleteIP("1.2.3.4")
			},
			expectedErr: func(err error) bool {
				return err == nil
			},
			expectedList: []firewall.IPEntry{
				{
					IP:        "2.2.8.8",
					CreatedAt: firewall.MustParseDateTime("2001-01-01 10:01:00"),
					UpdatedAt: firewall.MustParseDateTime("2001-01-01 10:01:00"),
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var fixedTime firewall.FixedTime

			service := firewall.NewService(
				firewall.WithTimeFunc(fixedTime.TimeFunc()),
				firewall.WithEchoWrapper(),
			)

			if tt.initBefore != nil {
				tt.initBefore(service, &fixedTime)
			}

			err := tt.testFunc(service, &fixedTime)
			if !tt.expectedErr(err) {
				t.Error()
				return
			}
			if err == nil {
				if !reflect.DeepEqual(service.List(), tt.expectedList) {
					t.Error()
				}
			}
		})
	}
}
