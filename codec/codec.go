package codec

import (
	"errors"
	"math"

	"github.com/gogo/protobuf/proto"
	"github.com/pokt-network/pocket-core/codec/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

type Codec struct {
	protoCdc        *ProtoCodec
	legacyCdc       *LegacyAmino
	upgradeOverride int
}

func NewCodec(anyUnpacker types.AnyUnpacker) *Codec {
	return &Codec{
		protoCdc:        NewProtoCodec(anyUnpacker),
		legacyCdc:       NewLegacyAminoCodec(),
		upgradeOverride: -1,
	}
}

var (
	UpgradeHeight                    int64 = math.MaxInt64
	OldUpgradeHeight                 int64 = 0
	NotProtoCompatibleInterfaceError       = errors.New("the interface passed for encoding does not implement proto marshaller")
)

const UpgradeCodecHeight = int64(30024)

//aux func to bypass on lower heights for test and because we are passing to tendermint
//the first upgrade will be treated as the "codec upgrade" with all the features that unlocked at that point the
//max block for that upgrade is UpgradeCodecHeight(30024)
func GetCodecUpgradeHeight() int64 {
	//check if this is the first upgrade
	if OldUpgradeHeight == 0 {
		if UpgradeHeight >= UpgradeCodecHeight {
			return UpgradeCodecHeight
		} else {
			return UpgradeHeight
		}
	} else {
		//this is not the first upgrade as OldUpgradeHeight is non 0
		//we return OldUpgradeHeight as the CodecUpgradeHeight
		return OldUpgradeHeight
	}
}

func (cdc *Codec) RegisterStructure(o interface{}, name string) {
	cdc.legacyCdc.RegisterConcrete(o, name, nil)
}

func (cdc *Codec) SetUpgradeOverride(b bool) {
	if b {
		cdc.upgradeOverride = 1
	} else {
		cdc.upgradeOverride = 0
	}
}

func (cdc *Codec) DisableUpgradeOverride() {
	cdc.upgradeOverride = -1
}

func (cdc *Codec) RegisterInterface(name string, iface interface{}, impls ...proto.Message) {
	res, ok := cdc.protoCdc.anyUnpacker.(types.InterfaceRegistry)
	if !ok {
		panic("unable to convert protocodec.anyUnpacker into types.InterfaceRegistry")
	}
	res.RegisterInterface(name, iface, impls...)
	cdc.legacyCdc.Amino.RegisterInterface(iface, nil)
}

func (cdc *Codec) RegisterImplementation(iface interface{}, impls ...proto.Message) {
	res, ok := cdc.protoCdc.anyUnpacker.(types.InterfaceRegistry)
	if !ok {
		panic("unable to convert protocodec.anyUnpacker into types.InterfaceRegistry")
	}
	res.RegisterImplementations(iface, impls...)
}

func (cdc *Codec) MarshalBinaryBare(o interface{}, height int64) ([]byte, error) { // TODO take height as parameter, move upgrade height to this package, switch based on height not upgrade mod
	p, ok := o.(ProtoMarshaler)
	if !ok {
		if cdc.IsAfterUpgrade(height) {
			return nil, NotProtoCompatibleInterfaceError
		}
		return cdc.legacyCdc.MarshalBinaryBare(o)
	} else {
		if cdc.IsAfterUpgrade(height) {
			return cdc.protoCdc.MarshalBinaryBare(p)
		}
		return cdc.legacyCdc.MarshalBinaryBare(p)
	}
}

func (cdc *Codec) MarshalBinaryLengthPrefixed(o interface{}, height int64) ([]byte, error) {
	p, ok := o.(ProtoMarshaler)
	if !ok {
		if cdc.IsAfterUpgrade(height) {
			return nil, NotProtoCompatibleInterfaceError
		}
		return cdc.legacyCdc.MarshalBinaryLengthPrefixed(o)
	} else {
		if cdc.IsAfterUpgrade(height) {
			return cdc.protoCdc.MarshalBinaryLengthPrefixed(p)
		}
		return cdc.legacyCdc.MarshalBinaryLengthPrefixed(p)
	}
}

func (cdc *Codec) UnmarshalBinaryBare(bz []byte, ptr interface{}, height int64) error {
	p, ok := ptr.(ProtoMarshaler)
	if !ok {
		if cdc.IsAfterUpgrade(height) {
			return NotProtoCompatibleInterfaceError
		}
		return cdc.legacyCdc.UnmarshalBinaryBare(bz, ptr)
	}
	if cdc.IsAfterUpgrade(height) {
		return cdc.protoCdc.UnmarshalBinaryBare(bz, p)
	}
	e := cdc.legacyCdc.UnmarshalBinaryBare(bz, ptr)
	if e != nil {
		return cdc.protoCdc.UnmarshalBinaryBare(bz, p)
	}
	return e
}

func (cdc *Codec) UnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}, height int64) error {
	p, ok := ptr.(ProtoMarshaler)
	if !ok {
		if cdc.IsAfterUpgrade(height) {
			return NotProtoCompatibleInterfaceError
		}
		return cdc.legacyCdc.UnmarshalBinaryLengthPrefixed(bz, ptr)
	}
	if cdc.IsAfterUpgrade(height) {
		return cdc.protoCdc.UnmarshalBinaryLengthPrefixed(bz, p)
	}
	e := cdc.legacyCdc.UnmarshalBinaryLengthPrefixed(bz, ptr)
	if e != nil {
		return cdc.protoCdc.UnmarshalBinaryLengthPrefixed(bz, p)
	}
	return e
}

func (cdc *Codec) ProtoMarshalBinaryBare(o ProtoMarshaler) ([]byte, error) {
	return cdc.protoCdc.MarshalBinaryBare(o)
}

func (cdc *Codec) LegacyMarshalBinaryBare(o interface{}) ([]byte, error) {
	return cdc.legacyCdc.MarshalBinaryBare(o)
}

func (cdc *Codec) ProtoUnmarshalBinaryBare(bz []byte, ptr ProtoMarshaler) error {
	return cdc.protoCdc.UnmarshalBinaryBare(bz, ptr)
}

func (cdc *Codec) LegacyUnmarshalBinaryBare(bz []byte, ptr interface{}) error {
	return cdc.legacyCdc.UnmarshalBinaryBare(bz, ptr)
}

func (cdc *Codec) ProtoMarshalBinaryLengthPrefixed(o ProtoMarshaler) ([]byte, error) {
	return cdc.protoCdc.MarshalBinaryLengthPrefixed(o)
}

func (cdc *Codec) LegacyMarshalBinaryLengthPrefixed(o interface{}) ([]byte, error) {
	return cdc.legacyCdc.MarshalBinaryLengthPrefixed(o)
}

func (cdc *Codec) ProtoUnmarshalBinaryLengthPrefixed(bz []byte, ptr ProtoMarshaler) error {
	return cdc.protoCdc.UnmarshalBinaryLengthPrefixed(bz, ptr)
}

func (cdc *Codec) LegacyUnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}) error {
	return cdc.legacyCdc.UnmarshalBinaryLengthPrefixed(bz, ptr)
}

func (cdc *Codec) MarshalJSONIndent(o interface{}, prefix string, indent string) ([]byte, error) {
	return cdc.legacyCdc.MarshalJSONIndent(o, prefix, indent)
}

func (cdc *Codec) MarshalJSON(o interface{}) ([]byte, error) {
	return cdc.legacyCdc.MarshalJSON(o)
}

func (cdc *Codec) UnmarshalJSON(bz []byte, o interface{}) error {
	return cdc.legacyCdc.UnmarshalJSON(bz, o)
}

func (cdc *Codec) MustMarshalJSON(o interface{}) []byte {
	bz, err := cdc.MarshalJSON(o)
	if err != nil {
		panic(err)
	}
	return bz
}

func (cdc *Codec) MustUnmarshalJSON(bz []byte, ptr interface{}) {
	err := cdc.UnmarshalJSON(bz, ptr)
	if err != nil {
		panic(err)
	}
}

func RegisterEvidences(legacy *LegacyAmino, _ *ProtoCodec) {
	tmTypes.RegisterEvidences(legacy.Amino)
}

func (cdc *Codec) AminoCodec() *LegacyAmino {
	return cdc.legacyCdc
}

func (cdc *Codec) ProtoCodec() *ProtoCodec {
	return cdc.protoCdc
}

//Note: includes the actual upgrade height
//now with shortcut on mainnet so we dont go back
func (cdc *Codec) IsAfterUpgrade(height int64) bool {
	//quick exit for mainnet
	if height >= UpgradeCodecHeight {
		return true
	} else {
		//original
		if cdc.upgradeOverride != -1 {
			return cdc.upgradeOverride == 1
		}
		return UpgradeHeight <= height || height == -1
	}
}

func (cdc *Codec) IsAfterNextUpgrade(height int64) bool {

	if OldUpgradeHeight == 0 || height <= OldUpgradeHeight {
		return false
	} else {
		return UpgradeHeight <= height
	}
}
