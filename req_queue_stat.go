package queuedproto

import (
	"fmt"
)

// ReqQueueStat - получение статистики
//
//adv:iproto:
type ReqQueueStat struct {
	Version     uint8
	StorageType uint16
}

// RespQueueStat - ответ на команду запроса статистики
//
//adv:iproto:
type RespQueueStat struct {
	ItemsCount  uint32 // кол-во событий в очереди
	ActiveCount uint32 // кол-во активных событий в очереди
	LockedCount uint32 // кол-во заблокированных событий в очереди
}

var (
	_ Request      = ReqQueueStat{}
	_ fmt.Stringer = RespQueueStat{}
	_ withQueueID  = ReqQueueStat{}
)

// Cmd - команда queued: CmdGetQueueStat (36)
func (ReqQueueStat) Cmd() Cmd {
	return CmdGetQueueStat
}

func (req ReqQueueStat) QueueID() uint16 {
	return req.StorageType
}

// String возвращает статистику в формате строки
func (i RespQueueStat) String() string {
	return fmt.Sprintf("Queued Stat: itemsCount% d, activeCount %d, lockedCount %d\n", i.ItemsCount, i.ActiveCount, i.LockedCount)
}
