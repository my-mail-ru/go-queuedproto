package addons

// ID перловых аддонов (queuedproto.PerlAddonData.ID)
const (
	CreationTimeID = uint16(1) // время создания события
	RetryID        = uint16(2) // сохраняет в событии счётчик ретраев, при объявлении очереди указывается лимит ретраев, и задержка
	Retry2ID       = uint16(3) // два независимых счётчика ретраев
	ProducerID     = uint16(4) // в перле: UNKNOWN, UWSGI, TP, SCRIPT. пока при создании событий в go пишем сюда TP-GO, хотя это и некорректно (продьюсят не только обработчики очередей)
	CreatedAtID    = uint16(5) // время создания события в наносекундах (int64)
	MaxAddonID
)

type (
	// PerlAddon - интерфейс всех перловых аддонов с поддержкой специфических для аддонов де/кодирующих функций
	PerlAddon interface {
		MarshalAddon() ([]byte, error) // не стал делать encoding.BinaryMarshaler, во избежание проброса на уровень юзерской структуры
		UnmarshalAddon([]byte) error
		AddonID() uint16
	}

	// Builder - для автозаполнения структур, используется в tp (копия интерфейса tp.Builder)
	Builder interface {
		Build() error
	}
)
