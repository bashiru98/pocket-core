package types

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/x/params"
)

// POS params default values
const (
	DefaultSessionNodeCount = 5
	MaxSessionNodes         = 25
)

// nolint - Keys for parameter access
var (
	KeySessionNodeCount = []byte("SessionNodeCount")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for pos module
type Params struct {
	SessionNodeCount uint
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeySessionNodeCount, Value: &p.SessionNodeCount},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		SessionNodeCount: DefaultSessionNodeCount,
	}
}

// validate a set of params
func (p Params) Validate() error {
	if p.SessionNodeCount > 25 {
		return errors.New("too many nodes per session, maximum is 5")
	}
	return nil
}

// Checks the equality of two param objects
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  SessionNodeCount:          %s`,
		p.SessionNodeCount)
}

// unmarshal the current pos params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	p, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return p
}

// unmarshal the current pos params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}