package addons //nolint:dupl

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/my-mail-ru/go-iproto"
)

// CreatedAt - время создания события (заполняется автоматически текущим временем)
type CreatedAt int64

// GetCreatedAt возвращает значение [CreatedAt], для которого он вызван. Использовать для получения времени создания из произвольной структуры события со встроенным [CreatedAt].
func (ca CreatedAt) GetCreatedAt() CreatedAt {
	return ca
}

var (
	_ PerlAddon                   = (*CreatedAt)(nil)
	_ Builder                     = (*CreatedAt)(nil)
	_ iproto.MarshalerUnmarshaler = (*CreatedAt)(nil)
)

// AddonID - возвращает CreatedAtID (5)
func (CreatedAt) AddonID() uint16 {
	return CreatedAtID
}

// MarshalIProto - no op (чтобы данные аддона не попали в стандартный payload).
//
// TODO выкинуть эти методы, когда библиотека iproto сможет парсить директивы встроенных в структуру типов.
func (CreatedAt) MarshalIProto(buf []byte) ([]byte, error) {
	return buf, nil
}

// UnmarshalIProto - no op (не пытаемся читать данные из стандартного payload).
func (CreatedAt) UnmarshalIProto(buf []byte) ([]byte, error) {
	return buf, nil
}

// MarshalAddon кодирует данные аддона
func (ca CreatedAt) MarshalAddon() ([]byte, error) {
	data := make([]byte, 8)

	binary.LittleEndian.PutUint64(data, uint64(ca))

	return data, nil
}

// UnmarshalAddon декодирует данные аддона
func (ca *CreatedAt) UnmarshalAddon(data []byte) error {
	if len(data) != 8 {
		*ca = 0

		return fmt.Errorf("CreatedAt.UnmarshalAddon: got len=%d, expected 8", len(data))
	}

	*ca = CreatedAt(binary.LittleEndian.Uint64(data))

	return nil
}

// Build записывает в *ca текущее время в наносекундах, если там дефолтное значение (0).
func (ca *CreatedAt) Build() error {
	if *ca == 0 {
		*ca = CreatedAt(time.Now().UnixNano())
	}

	return nil
}
