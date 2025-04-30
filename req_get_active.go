package queuedproto

// ReqGetActive - формат tuple запроса на получение Count активных событий из очереди
//
//adv:iproto:
type ReqGetActive struct {
	StorageType uint16 // номер очереди
	Count       uint32 // количество событий в ответе
}

var _ Request = ReqGetActive{}

// Cmd - команда queued: CmdGetActive (21)
func (ReqGetActive) Cmd() uint32 {
	return CmdGetActive
}
