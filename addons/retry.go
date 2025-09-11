package addons

import (
	"encoding/binary"
	"fmt"
)

// Retryable - счётчик повторных попыток обработки события.
type Retryable struct {
	data *retryData `iproto:"-"`
}

type retryData struct {
	retryCount uint32
	needRetry  bool
}

var _ PerlAddon = &Retryable{}

// AddonID - возвращает RetryID (2)
func (Retryable) AddonID() uint16 {
	return RetryID
}

// IncRetryCount увеличивает счётчик повторных обработок.
// Для использования из tp. Не следует вызывать этот метод из кода обработчиков.
// Для объектов, не полученных при помощи [Retryable.UnmarshalAddon], не делает ничего.
func (r Retryable) IncRetryCount() {
	if r.data != nil {
		r.data.retryCount++
	}
}

// Retry помечает событие подлежащим повторной обработке.
// Вызывать из обработчиков очередей в случае возникновения ошибки с ограниченным кол-вом повторов.
// Для объектов, не полученных при помощи [Retryable.UnmarshalAddon], не делает ничего.
func (r Retryable) Retry() {
	if r.data != nil {
		r.data.needRetry = true
	}
}

// GetRetryCount возвращает счётчик повторных попыток обработки.
// Для объектов, не полученных при помощи [Retryable.UnmarshalAddon], возвращает 0.
func (r Retryable) GetRetryCount() uint32 {
	if r.data == nil {
		return 0
	}

	return r.data.retryCount
}

// NeedRetry возвращает признак необходимости повторной обработки.
// Для объектов, не полученных при помощи [Retryable.UnmarshalAddon], возвращает false.
func (r Retryable) NeedRetry() bool {
	if r.data == nil {
		return false
	}

	return r.data.needRetry
}

// MarshalAddon кодирует данные аддона.
func (r Retryable) MarshalAddon() ([]byte, error) {
	data := make([]byte, 4)

	binary.LittleEndian.PutUint32(data, r.GetRetryCount())

	return data, nil
}

// UnmarshalAddon декодирует данные аддона.
func (r *Retryable) UnmarshalAddon(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("addons.Retryable: got len=%d, expected 4", len(data))
	}

	r.data = new(retryData)
	r.data.retryCount = binary.LittleEndian.Uint32(data)

	return nil
}
