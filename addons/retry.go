package addons

import (
	"encoding/binary"
	"fmt"
)

// Retry - счётчик повторных попыток обработки события.
type Retry struct {
	data *retryData `iproto:"-"`
}

type retryData struct {
	retryCount uint32
	hasFailed  bool
}

var _ PerlAddon = &Retry{}

// AddonID - возвращает RetryID (2)
func (Retry) AddonID() uint16 {
	return RetryID
}

// IncRetryCount увеличивает счётчик ошибок.
// Для использования из tp. Не следует вызывать этот метод из кода обработчиков.
// Для объектов, не полученных при помощи [Retry.UnmarshalAddon], не делает ничего.
func (r Retry) IncRetryCount() {
	if r.data != nil {
		r.data.retryCount++
	}
}

// Fail помечает событие ошибочным.
// Вызывать из обработчиков очередей в случае возникновения ошибки с ограниченным кол-вом повторов.
// Для объектов, не полученных при помощи [Retry.UnmarshalAddon], не делает ничего.
func (r Retry) Fail() {
	if r.data != nil {
		r.data.hasFailed = true
	}
}

// GetRetryCount возвращает счётчик ошибок.
// Для объектов, не полученных при помощи [Retry.UnmarshalAddon], возвращает 0.
func (r Retry) GetRetryCount() uint32 {
	if r.data != nil {
		return r.data.retryCount
	}

	return 0
}

// HasFailed возвращает признак ошибки.
// Для объектов, не полученных при помощи [Retry.UnmarshalAddon], возвращает false.
func (r Retry) HasFailed() bool {
	if r.data != nil {
		return r.data.hasFailed
	}

	return false
}

// MarshalAddon кодирует данные аддона.
func (r Retry) MarshalAddon() ([]byte, error) {
	data := make([]byte, 4)

	binary.LittleEndian.PutUint32(data, r.GetRetryCount())

	return data, nil
}

// UnmarshalAddon декодирует данные аддона.
func (r *Retry) UnmarshalAddon(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("addons.Retry: got len=%d, expected 4", len(data))
	}

	r.data = new(retryData)
	r.data.retryCount = binary.LittleEndian.Uint32(data)

	return nil
}
