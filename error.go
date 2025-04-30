package queuedproto

import (
	"errors"
	"fmt"
)

// Протокольные ошибки queued
var (
	ErrQueued              = errors.New("queued error")
	ErrLogic               = fmt.Errorf("%w: logic error", ErrQueued)
	ErrWrongItem           = fmt.Errorf("%w: wrong item", ErrQueued)
	ErrMemAllocationFailed = fmt.Errorf("%w: queued memory allocation failure", ErrQueued)
	ErrUnknownQueueType    = fmt.Errorf("%w: unknown queue type", ErrQueued)
	ErrReqAlreadyProcessed = fmt.Errorf("%w: request has been already processed (duplicate ReqID)", ErrQueued)
	ErrItemLocked          = fmt.Errorf("%w: item locked", ErrQueued)
	ErrWrongRequest        = fmt.Errorf("%w: wrong request", ErrQueued)
	ErrBadRequestLength    = fmt.Errorf("%w: bad request length", ErrQueued)
	ErrWrongVersion        = fmt.Errorf("%w: wrong version", ErrQueued)
)

var softErrors = map[error]struct{}{
	ErrLogic:               struct{}{},
	ErrWrongItem:           struct{}{},
	ErrUnknownQueueType:    struct{}{},
	ErrReqAlreadyProcessed: struct{}{},
	ErrWrongRequest:        struct{}{},
	ErrBadRequestLength:    struct{}{},
	ErrWrongVersion:        struct{}{},
}

// IsSoftError сообщает, является ли ошибка логической - т.е. не связанной с низкоуровневыми проблемами
// (разрывы соединения, таймауты, блокировки, сбои сервера).
//
// Все неизвестные протокольные ошибки считаются "жёсткими".
//
// Запросы с логическими ошибками повторять не надо, с низкоуровневыми - можно/нужно.
//
// Можно передавать также обёрнутые ошибки.
func IsSoftError(err error) bool {
	if err = unwrapErrQueued(err); err == nil {
		return false
	}

	_, ok := softErrors[err]

	return ok
}

// unwrapErrQueued разворачивает ошибку до ошибки, полученной оборачиванием ErrQueued.
// errors.Is понятнее, но медленнее (2 вложенных цикла, вместо 1 цикла с последующим 1 map lookup).
func unwrapErrQueued(err error) error {
	for {
		switch errPrev := errors.Unwrap(err); errPrev {
		case nil:
			return nil

		case ErrQueued:
			return err

		default:
			err = errPrev
		}
	}
}

// ErrorByRetCode - возвращает ошибку по коду, для неизвестных ошибок в тексте сообщается код.
// Для RcOK и всех ворнингов возвращается nil, если ворнинги нужны - их нужно обрабатывать вручную.
func ErrorByRetCode(rc uint8) error {
	switch rc {
	case RcOK, RcNoActiveItems, RcInsufficientActiveCount, RcWrongID:
		return nil

	case RcLogicError:
		return ErrLogic

	case RcWrongItem:
		return ErrWrongItem

	case RcMemAllocationFailed:
		return ErrMemAllocationFailed

	case RcUnknownQueueType:
		return ErrUnknownQueueType

	case RcReqAlreadyProcessed:
		return ErrReqAlreadyProcessed

	case RcItemLocked:
		return ErrItemLocked

	case RcWrongRequest:
		return ErrWrongRequest

	case RcBadRequestLength:
		return ErrBadRequestLength

	case RcWrongVersion:
		return ErrWrongVersion
	}

	return fmt.Errorf("%w: unknown code: 0x%02X", ErrQueued, rc)
}
