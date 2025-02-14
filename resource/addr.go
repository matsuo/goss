package resource

import (
	"fmt"
	"time"

	"github.com/goss-org/goss/system"
	"github.com/goss-org/goss/util"
)

type Addr struct {
	Title        string  `json:"title,omitempty" yaml:"title,omitempty"`
	Meta         meta    `json:"meta,omitempty" yaml:"meta,omitempty"`
	id           string  `json:"-" yaml:"-"`
	Address      string  `json:"address,omitempty" yaml:"address,omitempty"`
	LocalAddress string  `json:"local-address,omitempty" yaml:"local-address,omitempty"`
	Reachable    matcher `json:"reachable" yaml:"reachable"`
	Timeout      int     `json:"timeout" yaml:"timeout"`
	Skip         bool    `json:"skip,omitempty" yaml:"skip,omitempty"`
}

const (
	AddrResourceKey = "addr"
	AddResourceName = "Addr"
)

func init() {
	registerResource(AddrResourceKey, &Addr{})
}

func (a *Addr) ID() string {
	if a.Address != "" && a.Address != a.id {
		return fmt.Sprintf("%s: %s", a.id, a.Address)
	}
	return a.id
}
func (a *Addr) SetID(id string)  { a.id = id }
func (a *Addr) SetSkip()         { a.Skip = true }
func (a *Addr) TypeKey() string  { return AddrResourceKey }
func (a *Addr) TypeName() string { return AddResourceName }

// FIXME: Can this be refactored?
func (a *Addr) GetTitle() string { return a.Title }
func (a *Addr) GetMeta() meta    { return a.Meta }
func (a *Addr) GetAddress() string {
	if a.Address != "" {
		return a.Address
	}
	return a.id
}

func (a *Addr) Validate(sys *system.System) []TestResult {
	skip := a.Skip

	if a.Timeout == 0 {
		a.Timeout = 500
	}

	sysAddr := sys.NewAddr(a.GetAddress(), sys, util.Config{Timeout: time.Duration(a.Timeout) * time.Millisecond, LocalAddress: a.LocalAddress})

	var results []TestResult
	results = append(results, ValidateValue(a, "reachable", a.Reachable, sysAddr.Reachable, skip))
	return results
}

func NewAddr(sysAddr system.Addr, config util.Config) (*Addr, error) {
	address := sysAddr.Address()
	reachable, err := sysAddr.Reachable()
	a := &Addr{
		id:           address,
		Reachable:    reachable,
		Timeout:      config.TimeOutMilliSeconds(),
		LocalAddress: config.LocalAddress,
	}
	return a, err
}
