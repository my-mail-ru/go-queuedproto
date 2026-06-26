package queuedproto

import (
	"encoding/binary"
	"fmt"
	"slices"

	iproto "github.com/my-mail-ru/go-iproto"
)

// SkipTimeUpdate - для запроса [ReqFullUpdate] - magic value для UnixTime, в случае
// передачи этого значения время не изменяется.
const SkipTimeUpdate = 0xFFFFFFFF

// ReqFullUpdate - формат tuple запроса на обновление нескольких событий в очереди
//
//adv:iproto:
type ReqFullUpdate struct {
	StorageType uint16         // номер очереди
	UnixTime    uint32         // новое время активации обновлённых событий, если [SkipTimeUpdate] (0xFFFFFFFF), время не обновляется
	Items       []UpdQueueItem `iproto:"u32"`
}

var (
	_ Request     = ReqFullUpdate{}
	_ withQueueID = ReqFullUpdate{}
)

// UpdQueueItem - аналог структуры [QueueItem] для запроса [ReqFullUpdate].
// Если Data == nil, данные не обновляются - только время активации.
// Чтобы передать пустой массив данных, надо указать []byte{}.
type UpdQueueItem struct {
	EventID EventID
	Data    []byte `iproto:"u16"` // Данные (до 4кБ)
}

// Cmd - команда queued: CmdFullUpdate (30)
func (ReqFullUpdate) Cmd() Cmd {
	return CmdFullUpdate
}

func (req ReqFullUpdate) QueueID() uint16 {
	return req.StorageType
}

// MarshalIProto кодирует структуру UpdQueueItem.
// Если Data == nil, передаётся длина 0xFFFF.
func (uqi UpdQueueItem) MarshalIProto(buf []byte) ([]byte, error) {
	buf, err := uqi.EventID.MarshalIProto(buf)
	if err != nil {
		return nil, err
	}

	if uqi.Data == nil {
		return append(buf, 0xFF, 0xFF), nil
	}

	l := len(uqi.Data)

	return append(append(buf, byte(l&0xFF), byte(l>>8)), uqi.Data...), nil
}

// UnmarshalIProto декодирует структуру UpdQueueItem.
// Если переданæ длина данных 0xFFFF, поле Data будет равняться nil.
func (uqi *UpdQueueItem) UnmarshalIProto(buf []byte) ([]byte, error) {
	buf, err := uqi.EventID.UnmarshalIProto(buf)
	if err != nil {
		return nil, err
	}

	if len(buf) < 2 {
		return nil, fmt.Errorf("UpdQueueItem.UnmarshalIProto: %w", iproto.ErrOverflow)
	}

	l := int(binary.LittleEndian.Uint16(buf))
	buf = buf[2:]

	if l == 0xFFFF {
		uqi.Data = nil
		return buf, nil
	}

	if len(buf) < l {
		return nil, fmt.Errorf("UpdQueueItem.UnmarshalIProto: %w: %d < %d", iproto.ErrOverflow, len(buf), l)
	}

	uqi.Data = slices.Clone(buf[:l])

	return buf[l:], nil
}
