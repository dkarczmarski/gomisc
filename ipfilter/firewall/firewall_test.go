package firewall_test

import (
	"github.com/dkarczmarski/gomisc/ipfilter/firewall"
	"reflect"
	"testing"
)

func TestService(t *testing.T) {
	for _, tt := range []struct {
		name         string
		initBefore   func(service *firewall.Service)
		ip           string
		expectedErr  func(err error) bool
		expectedList []firewall.IPEntry
	}{
		{
			name: "incorrect ip",
			ip:   "1.2.3,,4",
			expectedErr: func(err error) bool {
				return err != nil
			},
		},
		{
			name: "add ip when empty",
			ip:   "1.2.3.4",
			expectedErr: func(err error) bool {
				return err == nil
			},
			expectedList: []firewall.IPEntry{
				{
					IP: "1.2.3.4",
				},
			},
		},
		{
			name: "add another ip",
			ip:   "2.2.8.8",
			initBefore: func(service *firewall.Service) {
				_ = service.AddIP("1.2.3.4")
			},
			expectedErr: func(err error) bool {
				return err == nil
			},
			expectedList: []firewall.IPEntry{
				{IP: "1.2.3.4"},
				{IP: "2.2.8.8"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			service := &firewall.Service{}

			if tt.initBefore != nil {
				tt.initBefore(service)
			}

			err := service.AddIP(tt.ip)
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

func TestService_AddIP(t *testing.T) {
	service := firewall.Service{}

	if err := service.AddIP("123.2.3.7"); err != nil {
		t.Error(err)
	}

	if err := service.DeleteIP("123.2.3.7"); err != nil {
		t.Error(err)
	}

	if err := service.DeleteIP("123.2.3.7"); err == nil {
		t.Error("should be an error")
	}
}
