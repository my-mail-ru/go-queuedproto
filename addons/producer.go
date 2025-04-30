package addons

import (
	"github.com/my-mail-ru/go-iproto"
	iprototypes "github.com/my-mail-ru/go-iproto/types"
)

/*
Producer - имя источника события (для статы/метрик, заполняется автоматически значением "TP-GO").

Для отключения автозаполнения Producer, либо задания своего значения, проще всего определить свой тип:

	type MyProducer struct {
		addons.Producer
	}

	func (mp *MyProducer) Build() error {
		mp.Producer = "GO-REST"
		return nil
	}

	type MyEvent struct {
		MyProducer
		...
	}
*/
type Producer string

// DefaultProducer - Producer по умолчанию (для перлового формата, если подключён этот аддон)
//
// TODO добавить стандартные константы для других источников событий.
const DefaultProducer = "TP-GO"

var (
	_ PerlAddon                   = (*Producer)(nil)
	_ Builder                     = (*Producer)(nil)
	_ iproto.MarshalerUnmarshaler = (*Producer)(nil)
)

// AddonID - возвращает ProducerID (4)
func (Producer) AddonID() uint16 {
	return ProducerID
}

// MarshalIProto - no op (чтобы данные аддона не попали в стандартный payload).
//
// TODO выкинуть эти методы, когда библиотека iproto сможет парзить директивы встроенных в структуру типов.
func (Producer) MarshalIProto(buf []byte) ([]byte, error) {
	return buf, nil
}

// UnmarshalIProto - no op (не пытаемся читать данные из стандартного payload).
func (Producer) UnmarshalIProto(buf []byte) ([]byte, error) {
	return buf, nil
}

// MarshalAddon кодирует данные аддона
func (p Producer) MarshalAddon() ([]byte, error) {
	return iprototypes.String(p).MarshalIProto(make([]byte, 0, len(p)+4))
}

// UnmarshalAddon декодирует данные аддона
func (p *Producer) UnmarshalAddon(data []byte) error {
	var s iprototypes.String

	if _, err := s.UnmarshalIProto(data); err != nil {
		return err
	}

	*p = Producer(s)

	return nil
}

// Build записывает в *p значение DefaultProducer, если в *p дефолтное значение (пустая строка).
func (p *Producer) Build() error {
	if *p == "" {
		*p = DefaultProducer
	}

	return nil
}
