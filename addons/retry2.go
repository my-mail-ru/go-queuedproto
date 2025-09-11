package addons

import (
	"encoding/binary"
	"fmt"
)

// Retryable2 - два независимых счётчика повторных попыток обработки события.
type Retryable2 struct {
	data *retry2Data `iproto:"-"`
}

type retry2Data struct {
	retry1Count uint32
	retry2Count uint32
	needRetry1  bool
	needRetry2  bool
}

var _ PerlAddon = &Retryable2{}

// AddonID - возвращает Retry2ID (3)
func (Retryable2) AddonID() uint16 {
	return Retry2ID
}

// IncRetry1Count увеличивает первый счётчик повторных попыток обработки.
// Для использования из tp. Не следует вызывать этот метод из кода обработчиков.
// Для объектов, не полученных при помощи [Retryable2.UnmarshalAddon], не делает ничего.
func (r Retryable2) IncRetry1Count() {
	if r.data != nil {
		r.data.retry1Count++
	}
}

// IncRetry2Count увеличивает второй счётчик повторных попыток обработки.
// Для использования из tp. Не следует вызывать этот метод из кода обработчиков.
// Для объектов, не полученных при помощи [Retryable2.UnmarshalAddon], не делает ничего.
func (r Retryable2) IncRetry2Count() {
	if r.data != nil {
		r.data.retry2Count++
	}
}

// Retry1 помечает событие подлежащим повторной обработке с учётом первого счётчика.
// Вызывать из обработчиков очередей в случае возникновения ошибки с ограниченным кол-вом повторов.
// Для объектов, не полученных при помощи [Retryable2.UnmarshalAddon], не делает ничего.
func (r Retryable2) Retry1() {
	if r.data != nil {
		r.data.needRetry1 = true
	}
}

// Retry2 помечает событие подлежащим повторной обработке с учётом второго счётчика.
// Вызывать из обработчиков очередей в случае возникновения ошибки с ограниченным кол-вом повторов.
// Для объектов, не полученных при помощи [Retryable2.UnmarshalAddon], не делает ничего.
func (r Retryable2) Retry2() {
	if r.data != nil {
		r.data.needRetry2 = true
	}
}

// GetRetry1Count возвращает первый счётчик повторных попыток обработки.
// Для объектов, не полученных при помощи [Retryable2.UnmarshalAddon], возвращает 0.
func (r Retryable2) GetRetry1Count() uint32 {
	if r.data == nil {
		return 0
	}

	return r.data.retry1Count
}

// GetRetry2Count возвращает второй счётчик повторных попыток обработки.
// Для объектов, не полученных при помощи [Retryable2.UnmarshalAddon], возвращает 0.
func (r Retryable2) GetRetry2Count() uint32 {
	if r.data == nil {
		return 0
	}

	return r.data.retry2Count
}

// NeedRetry1 возвращает первый признак необходимости повторной обработки.
// Для объектов, не полученных при помощи [Retryable2.UnmarshalAddon], возвращает false.
func (r Retryable2) NeedRetry1() bool {
	if r.data == nil {
		return false
	}

	return r.data.needRetry1
}

// NeedRetry2 возвращает второй признак необходимости повторной обработки.
// Для объектов, не полученных при помощи [Retryable2.UnmarshalAddon], возвращает false.
func (r Retryable2) NeedRetry2() bool {
	if r.data == nil {
		return false
	}

	return r.data.needRetry2
}

// MarshalAddon кодирует данные аддона.
func (r Retryable2) MarshalAddon() ([]byte, error) {
	data := make([]byte, 8)

	binary.LittleEndian.PutUint32(data, r.GetRetry1Count())
	binary.LittleEndian.PutUint32(data[4:], r.GetRetry2Count())

	return data, nil
}

// UnmarshalAddon декодирует данные аддона.
func (r *Retryable2) UnmarshalAddon(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("addons.Retryable2: got len=%d, expected 8", len(data))
	}

	r.data = new(retry2Data)
	r.data.retry1Count = binary.LittleEndian.Uint32(data)
	r.data.retry2Count = binary.LittleEndian.Uint32(data[4:])

	return nil
}
