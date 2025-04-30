package queuedproto

// ReqDeleteItems - формат tuple запроса на удаление событий по списку ID
//
//adv:iproto:
type ReqDeleteItems struct {
	StorageType uint16   // номер очереди
	IDs         []uint64 `iproto:"u32"`
}

var _ Request = ReqDeleteItems{}

// Cmd - команда queued: CmdDeleteItems (22)
func (ReqDeleteItems) Cmd() uint32 {
	return CmdDeleteItems
}
