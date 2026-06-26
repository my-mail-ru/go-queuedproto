package queuedproto

// ReqDeleteItems - формат tuple запроса на удаление событий по списку ID
//
//adv:iproto:
type ReqDeleteItems struct {
	StorageType uint16   // номер очереди
	IDs         []uint64 `iproto:"u32"`
}

var (
	_ Request     = ReqDeleteItems{}
	_ withQueueID = ReqDeleteItems{}
)

// Cmd - команда queued: CmdDeleteItems (22)
func (ReqDeleteItems) Cmd() Cmd {
	return CmdDeleteItems
}

func (req ReqDeleteItems) QueueID() uint16 {
	return req.StorageType
}
