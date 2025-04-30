package queuedproto

import (
	"slices"

	iproto "github.com/my-mail-ru/go-iproto"
	iprototypes "github.com/my-mail-ru/go-iproto/types"
)

// PerlData - событие очереди в перловом формате.
//
// Аддоны декодируются, только если установлен флаг FlagEnableAddons, иначе Addons.List устанавливается в nil.
// При кодировании, если len(Addons.List) != 0, они кодируются, и автоматически устанавливается флаг FlagEnableAddons.
//
// Extensions - аналогично, однако их внутренности (пока) не обрабатываются (просто декодируются, как слайс байтов).
//
// При кодировании флаг FlagEncodeInUTF8 устанавливается автоматически всегда.
//
//adv:iproto:
type PerlData struct {
	Version    uint8          // версия (содержит флаги FlagEnableExtensions, FlagEncodeInUTF8, FlagEnableAddons)
	Addons     PerlAddons     // список данных аддонов  (если установлен FlagEnableAddons - иначе поле не передаётся!)
	Extensions PerlExtensions // данные расширения (если установлен FlagEnableExtensions - иначе поле не передаётся!)
	Data       []byte         // данные в формате perl pack
}

// PerlAddonData - данные аддона. Константы ID и структуры формата данных оперделены в пакете addons
type PerlAddonData struct {
	ID   uint16 // addons.CreationTimeID, addons.RetryID, addons.Retry2ID или addons.ProducerID
	Data []byte `iproto:"u8"`
}

// PerlAddons - список аддонов
//
//adv:iproto:
type PerlAddons struct {
	List []PerlAddonData `iproto:"u8"`
}

// PerlExtensions - данные расширений (пока не обрабатываются, де/кодируются как слайс байтов)
//
//adv:iproto:
type PerlExtensions struct {
	Flags uint32
	Data  []byte `iproto:"u32"`
}

var (
	_ iproto.Marshaler   = PerlData{}
	_ iproto.Unmarshaler = &PerlData{}
)

// MarshalIProto кодирует данные в перловом формате, передавая необязательные поля, только если они заданы. В этом случае устанавливаются соответствующие им флаги (FlagEnableAddons - для Addons, FlagEnableExtensions - для Extensions).
func (pd PerlData) MarshalIProto(buf []byte) ([]byte, error) {
	version := pd.Version | FlagEncodeInUTF8

	if len(pd.Addons.List) != 0 {
		version |= FlagEnableAddons
	}

	if len(pd.Extensions.Data) != 0 {
		version |= FlagEnableExtensions
	}

	buf = append(buf, version)

	var err error
	if len(pd.Addons.List) != 0 {
		if buf, err = pd.Addons.MarshalIProto(buf); err != nil {
			return nil, err
		}
	}

	if len(pd.Extensions.Data) != 0 {
		if buf, err = pd.Extensions.MarshalIProto(buf); err != nil {
			return nil, err
		}
	}

	return append(buf, pd.Data...), nil
}

// UnmarshalIProto декодирует данные в перловом формате, обрабатывая необязательные поля, только если установлены соответствующие им флаги (FlagEnableAddons - для Addons, FlagEnableExtensions - для Extensions).
//
// Если флаг для необязательного поля не установлен, оно сбрасывается в дефлотное значение, так сделано, чтобы не прилетели данные от предыдущей записи в случае переиспользования объекта.
func (pd *PerlData) UnmarshalIProto(buf []byte) ([]byte, error) {
	var err error

	var version iprototypes.Uint8
	if buf, err = version.UnmarshalIProto(buf); err != nil {
		return nil, err
	}

	pd.Version = uint8(version)

	if pd.Version&FlagEnableAddons != 0 {
		if buf, err = pd.Addons.UnmarshalIProto(buf); err != nil {
			return nil, err
		}
	} else {
		pd.Addons.List = nil
	}

	if pd.Version&FlagEnableExtensions != 0 {
		if buf, err = pd.Extensions.UnmarshalIProto(buf); err != nil {
			return nil, err
		}
	} else {
		pd.Extensions = PerlExtensions{}
	}

	pd.Version &= VersionMask

	pd.Data = slices.Clone(buf)

	return buf[:0], nil
}
