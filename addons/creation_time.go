package addons

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/my-mail-ru/go-iproto"
)

// CreationTime - время создания события (заполняется автоматически текущим временем)
type CreationTime uint32

// GetCreationTime возвращает значение [CreationTime], для которого он вызван. Использовать для получения времени создания из произвольной структуры события со встроенным [CreationTime].
func (ct CreationTime) GetCreationTime() CreationTime {
	return ct
}

var (
	_ PerlAddon                   = (*CreationTime)(nil)
	_ Builder                     = (*CreationTime)(nil)
	_ iproto.MarshalerUnmarshaler = (*CreationTime)(nil)
)

// AddonID - возвращает CreationTimeID (1)
func (CreationTime) AddonID() uint16 {
	return CreationTimeID
}

// MarshalIProto - no op (чтобы данные аддона не попали в стандартный payload).
//
// TODO выкинуть эти методы, когда библиотека iproto сможет парзить директивы встроенных в структуру типов.
func (CreationTime) MarshalIProto(buf []byte) ([]byte, error) {
	return buf, nil
}

// UnmarshalIProto - no op (не пытаемся читать данные из стандартного payload).
func (CreationTime) UnmarshalIProto(buf []byte) ([]byte, error) {
	return buf, nil
}

// MarshalAddon кодирует данные аддона
func (ct CreationTime) MarshalAddon() ([]byte, error) {
	data := make([]byte, 4)

	binary.LittleEndian.PutUint32(data, uint32(ct))

	return data, nil
}

// UnmarshalAddon декодирует данные аддона
func (ct *CreationTime) UnmarshalAddon(data []byte) error {
	if len(data) != 4 {
		*ct = 0

		return fmt.Errorf("CreationTime.UnmarshalAddon: got len=%d, expected 4", len(data))
	}

	*ct = CreationTime(binary.LittleEndian.Uint32(data))

	return nil
}

// Build записывает в *ct текущее время, если там дефолтное значение (0).
func (ct *CreationTime) Build() error {
	if *ct == 0 {
		*ct = CreationTime(time.Now().Unix()) //nolint:gosec
	}

	return nil
}
