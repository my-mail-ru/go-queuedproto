package queuedproto

import (
	"encoding"
	"fmt"

	"github.com/my-mail-ru/go-iproto"
)

//go:generate go tool iprotogen -r -tests

type (
	// Request - структуры команд протокола должны поддерживать этот интерфейс (возвращать код команды методом Cmd)
	Request interface {
		iproto.Marshaler
		Cmd() uint32
	}

	// EventID - составной ID с номером шарда. Нужен, т.к. протокол (и сам queued) не поддерживает шардинг.
	// Связка Shard:ID уникальна, тогда как сам по себе ID в списке событий может быть неуникальным, если события
	// пришли из разных шардов. Такое возможно только при автогенерации ID событий на стороне queued.
	//
	// В структурах запросах этот тип используется только для слайсов (чтобы не перелопачивать слайсы целиком из параметров метода,
	// которому номер шарда должен как-то поступать). Скалярные айдишники надо будет скопировать (только ID, передающийся по прококолу).
	//adv:iproto:
	EventID struct {
		Shard uint `iproto:"-"`
		ID    uint64
	}

	// QueueItem - событие в очереди
	QueueItem struct {
		EventID EventID
		Data    []byte `iproto:"u16"` // Данные (до 4кБ)
	}

	// ItemList - список событий
	//adv:iproto:
	ItemList struct {
		Items []QueueItem `iproto:"u32"`
	}
)

var (
	_ fmt.Stringer           = EventID{}
	_ encoding.TextMarshaler = EventID{}
)

// String возвращает EventID в перловом текстовом формате
func (eid EventID) String() string {
	return fmt.Sprintf("%016X:%d", eid.ID, eid.Shard)
}

// MarshalText - для упрощения логгирования списков событий при помощи zerolog.Event.Interface (не нужно переваливать слайс
// айдишников событий в слайс строк/стрингеров).
func (eid EventID) MarshalText() ([]byte, error) {
	return []byte(eid.String()), nil
}
