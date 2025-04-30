package queuedproto

import (
	"github.com/my-mail-ru/go-iproto"
	iprototypes "github.com/my-mail-ru/go-iproto/types"
)

// ReqAddItem - формат tuple запроса на создание новой записи в очереди.
type ReqAddItem struct {
	Pid         uint32 // pid клиента
	ReqID       uint32 // номер запроса, в случае таймаута на сервере при перепосылке запроса ReqID 2й попытки должен совпадать с ReqID первой попытки
	StorageType uint16 // номер очереди
	ID          uint64 // ID события, 0 - queued сам назначает ID (автоинкремент)
	UnixTime    uint32 // время активации события
	Data        []byte `iproto:"u16"` // данные события, до 4096 байт
	Flags       uint32 `iproto:"-"`   // FlagIgnoreWrongItem. Поле опциональное (передаётся, только если не 0)
}

//adv:iproto:
type reqAddItemDumb ReqAddItem // отрезаем у типа все объявленные методы, чтобы не было рекурсии маршалеров

var (
	_ Request            = ReqAddItem{}
	_ iproto.Unmarshaler = &ReqAddItem{}
)

// Cmd - команда queued: CmdAddItem (20)
func (ReqAddItem) Cmd() uint32 {
	return CmdAddItem
}

// MarshalIProto кодирует запрос на создание новой записи в очереди.
//
// Кастомный маршалер необходим из-за опциональности поля Flags.
func (req ReqAddItem) MarshalIProto(buf []byte) ([]byte, error) {
	buf, err := reqAddItemDumb(req).MarshalIProto(buf)
	if err != nil {
		return nil, err
	}

	if req.Flags != 0 {
		buf, err = iprototypes.Uint32(req.Flags).MarshalIProto(buf)
	}

	return buf, err
}

// UnmarshalIProto декодирует запрос на создание новой записи в очереди.
//
// Метод объявлен чисто для симметрии маршалеров, для проектов клиентов/обработчиков очередей он не нужен.
// Пригодится для тестов, прокси iproto/queued->grpc, ну и для гошной версии queued :)
func (req *ReqAddItem) UnmarshalIProto(buf []byte) ([]byte, error) {
	buf, err := (*reqAddItemDumb)(req).UnmarshalIProto(buf)
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return buf, nil
	}

	var flags iprototypes.Uint32
	buf, err = flags.UnmarshalIProto(buf)
	req.Flags = uint32(flags)

	return buf, err
}
