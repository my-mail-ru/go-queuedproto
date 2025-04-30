package addons

import (
	"encoding/binary"
	"fmt"
)

// Retry2 - два независимых счётчика повторных попыток.
type Retry2 struct {
	data *retry2Data `iproto:"-"`
}

type retry2Data struct {
	retry1Count uint32
	retry2Count uint32
	hasFailed1  bool
	hasFailed2  bool
}

var _ PerlAddon = &Retry2{}

// AddonID - возвращает Retry2ID (3)
func (Retry2) AddonID() uint16 {
	return Retry2ID
}

// IncRetry1Count увеличивает первый счётчик ошибок.
// Для использования из tp. Не следует вызывать этот метод из кода обработчиков.
// Для объектов, не полученных при помощи [Retry2.UnmarshalAddon], не делает ничего.
func (r Retry2) IncRetry1Count() {
	if r.data != nil {
		r.data.retry1Count++
	}
}

// IncRetry2Count увеличивает первый счётчик ошибок.
// Для использования из tp. Не следует вызывать этот метод из кода обработчиков.
// Для объектов, не полученных при помощи [Retry2.UnmarshalAddon], не делает ничего.
func (r Retry2) IncRetry2Count() {
	if r.data != nil {
		r.data.retry2Count++
	}
}

// Fail1 помечает событие ошибочным с учётом первого счётчика.
// Вызывать из обработчиков очередей в случае возникновения ошибки с ограниченным кол-вом повторов.
// Для объектов, не полученных при помощи [Retry2.UnmarshalAddon], не делает ничего.
func (r Retry2) Fail1() {
	if r.data != nil {
		r.data.hasFailed1 = true
	}
}

// Fail2 помечает событие ошибочным с учётом первого счётчика.
// Вызывать из обработчиков очередей в случае возникновения ошибки с ограниченным кол-вом повторов.
// Для объектов, не полученных при помощи [Retry2.UnmarshalAddon], не делает ничего.
func (r Retry2) Fail2() {
	if r.data != nil {
		r.data.hasFailed2 = true
	}
}

// GetRetry1Count возвращает первый счётчик ошибок.
// Для объектов, не полученных при помощи [Retry2.UnmarshalAddon], возвращает 0.
func (r Retry2) GetRetry1Count() uint32 {
	if r.data != nil {
		return r.data.retry1Count
	}

	return 0
}

// GetRetry2Count возвращает второй счётчик ошибок.
// Для объектов, не полученных при помощи [Retry2.UnmarshalAddon], возвращает 0.
func (r Retry2) GetRetry2Count() uint32 {
	if r.data != nil {
		return r.data.retry2Count
	}

	return 0
}

// HasFailed1 возвращает первый признак ошибки.
// Для объектов, не полученных при помощи [Retry2.UnmarshalAddon], возвращает false.
func (r Retry2) HasFailed1() bool {
	if r.data != nil {
		return r.data.hasFailed1
	}

	return false
}

// HasFailed2 возвращает первый признак ошибки.
// Для объектов, не полученных при помощи [Retry2.UnmarshalAddon], возвращает false.
func (r Retry2) HasFailed2() bool {
	if r.data != nil {
		return r.data.hasFailed2
	}

	return false
}

// MarshalAddon кодирует данные аддона.
func (r Retry2) MarshalAddon() ([]byte, error) {
	data := make([]byte, 8)

	binary.LittleEndian.PutUint32(data, r.GetRetry1Count())
	binary.LittleEndian.PutUint32(data[4:], r.GetRetry2Count())

	return data, nil
}

// UnmarshalAddon декодирует данные аддона.
func (r *Retry2) UnmarshalAddon(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("addons.Retry2: got len=%d, expected 8", len(data))
	}

	r.data = new(retry2Data)
	r.data.retry1Count = binary.LittleEndian.Uint32(data)
	r.data.retry2Count = binary.LittleEndian.Uint32(data[4:])

	return nil
}
