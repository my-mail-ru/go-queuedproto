package queuedproto

//go:generate go tool stringer -type Cmd

// Коды команд queued
const (
	CmdAddItem      = Cmd(20)
	CmdGetActive    = Cmd(21)
	CmdDeleteItems  = Cmd(22)
	CmdUpdateItems  = Cmd(24)
	CmdGetItems     = Cmd(28)
	CmdFullUpdate   = Cmd(30)
	CmdAddData      = Cmd(34)
	CmdGetQueueStat = Cmd(36)
)

// Коды ответов queued
const (
	// errors
	RcOK                  = uint8(0)
	RcLogicError          = uint8(1)
	RcWrongItem           = uint8(2)
	RcMemAllocationFailed = uint8(3)
	RcUnknownQueueType    = uint8(4)
	RcReqAlreadyProcessed = uint8(8)
	RcItemLocked          = uint8(9)
	RcWrongRequest        = uint8(10)
	RcBadRequestLength    = uint8(11)
	RcWrongVersion        = uint8(12)

	// warnings
	RcNoActiveItems           = uint8(5)
	RcInsufficientActiveCount = uint8(6)
	RcWrongID                 = uint8(7)
)

// Флаги расширений, указываются в поле Version
const (
	FlagEnableExtensions = uint8(0x80)
	FlagEncodeInUTF8     = uint8(0x40)
	FlagEnableAddons     = uint8(0x20)

	VersionMask = ^(FlagEnableExtensions | FlagEncodeInUTF8 | FlagEnableAddons)
)

// FlagIgnoreWrongItem - флаг команды CmdAddItem
const FlagIgnoreWrongItem = uint32(1)

// Типы времени для команды CmdUpdateItems
const (
	TimestampAbsolute      = uint8(0) // указанный UnixTime - абсолютный (секунды с 1970-01-01T00:00:00Z)
	TimestampRelativePlus  = uint8(1) // в Unixtime относительное время, которое необходимо прибавить к текущему серверному времени
	TimestampRelativeMinus = uint8(2) // в Unixtime относительное время, которое необходимо отнять от текущего серверного времени
)
