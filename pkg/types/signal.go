package types

import (
	"fmt"
)

type SignalType string

const (
	SignalTypeMetric SignalType = "metric"
	SignalTypeLog    SignalType = "log"
	SignalTypeTrace  SignalType = "trace"
)

func (q SignalType) Validate() error {
	switch q {
	case SignalTypeMetric, SignalTypeLog, SignalTypeTrace:
		return nil
	}
	return fmt.Errorf("invalid signal type: %s", q)
}
