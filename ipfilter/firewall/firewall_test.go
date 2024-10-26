package firewall_test

import (
	"github.com/dkarczmarski/gomisc/ipfilter/firewall"
	"reflect"
	"testing"
)

func TestService_AddDeleteIP(t *testing.T) {
	for _, tt := range []struct {
		name         string
		initBefore   func(service *firewall.Service)
		testFunc     func(service *firewall.Service) error
		expectedErr  func(err error) bool
		expectedList []firewall.IPEntry
	}{
		{
			name: "add incorrect ip",
			testFunc: func(service *firewall.Service) error {
				return service.AddIP("1.2.3,,4")
			},
			expectedErr: func(err error) bool {
				return err != nil
			},
		},
		{
			name: "add first ip",
			testFunc: func(service *firewall.Service) error {
				return service.AddIP("1.2.3.4")
			},
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
			initBefore: func(service *firewall.Service) {
				_ = service.AddIP("1.2.3.4")
			},
			testFunc: func(service *firewall.Service) error {
				return service.AddIP("2.2.8.8")
			},
			expectedErr: func(err error) bool {
				return err == nil
			},
			expectedList: []firewall.IPEntry{
				{IP: "1.2.3.4"},
				{IP: "2.2.8.8"},
			},
		},
		{
			name: "delete incorrect ip",
			testFunc: func(service *firewall.Service) error {
				return service.DeleteIP("1.2.3,,4")
			},
			expectedErr: func(err error) bool {
				return err != nil
			},
		},
		{
			name: "delete nonexistent ip",
			testFunc: func(service *firewall.Service) error {
				return service.DeleteIP("1.2.3.4")
			},
			expectedErr: func(err error) bool {
				return err != nil
			},
		},
		{
			name: "add existed ip",
			initBefore: func(service *firewall.Service) {
				_ = service.AddIP("1.2.3.4")
				_ = service.AddIP("2.2.8.8")
			},
			testFunc: func(service *firewall.Service) error {
				return service.DeleteIP("1.2.3.4")
			},
			expectedErr: func(err error) bool {
				return err == nil
			},
			expectedList: []firewall.IPEntry{
				{IP: "2.2.8.8"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			service := &firewall.Service{
				WrapperCmd: "echo",
			}

			if tt.initBefore != nil {
				tt.initBefore(service)
			}

			err := tt.testFunc(service)
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
