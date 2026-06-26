package queuedproto

// ReqAddData - формат tuple запроса на добавление данных в конец события
//
//adv:iproto:
type ReqAddData struct {
	StorageType uint16 // номер очереди
	ID          uint64 // ID изменяемого события
	Data        []byte `iproto:"u16"` // дописываемые данные (в сумме с существующими должно быть до 4кБ)
}

var (
	_ Request     = ReqAddData{}
	_ withQueueID = ReqAddData{}
)

// Cmd - команда queued: CmdAddData (34)
func (ReqAddData) Cmd() Cmd {
	return CmdAddData
}

func (req ReqAddData) QueueID() uint16 {
	return req.StorageType
}
