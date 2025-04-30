package queuedproto

// ReqGetItems - формат tuple запроса на получение событий по ID
//
//adv:iproto:
type ReqGetItems struct {
	StorageType uint16   // номер очереди
	IDs         []uint64 `iproto:"u32"`
}

var _ Request = ReqGetItems{}

// Cmd - команда queued: CmdGetItems (28)
func (ReqGetItems) Cmd() uint32 {
	return CmdGetItems
}
