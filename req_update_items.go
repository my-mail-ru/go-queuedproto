package queuedproto

// ReqUpdateItems - формат tuple запроса на обновление времени срабатывания события
//
//adv:iproto:
type ReqUpdateItems struct {
	StorageType   uint16   // номер очереди
	UnixTime      uint32   // новое время активации для всех указанных событий
	IDs           []uint64 `iproto:"u32"`
	TimestampType uint8    // TimestampAbsolute, TimestampRelativePlus или TimestampRelativeMinus
}

var (
	_ Request     = ReqUpdateItems{}
	_ withQueueID = ReqUpdateItems{}
)

// Cmd - команда queued: CmdUpdateItems(24)
func (ReqUpdateItems) Cmd() Cmd {
	return CmdUpdateItems
}

func (req ReqUpdateItems) QueueID() uint16 {
	return req.StorageType
}
