package firewall

import (
	"fmt"
	"log"
	"time"
)

type FixedTime struct {
	time.Time
}

func (ft *FixedTime) SetDateTime(value string) {
	ft.Time = MustParseDateTime(value)
}

func (ft *FixedTime) TimeFunc() func() time.Time {
	return func() time.Time {
		return ft.Time
	}
}

func MustParseDateTime(value string) time.Time {
	t, err := time.Parse(time.DateTime, value)
	if err != nil {
		log.Panic(fmt.Errorf("time.Parse() %s: %w", value, err))
	}
	return t
}
