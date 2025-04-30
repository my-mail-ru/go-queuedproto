package queuedproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// тест на правильную обработку опционального поля Flags

var testReqAddItemHeader = []byte{
	0x5C, 0x61, 0x00, 0x00,
	0x01, 0x00, 0x00, 0x00, 0x64, 0xE1, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xA2, 0xE0,
	0x14, 0x62, 0x24, 0x00,
}

var testReqAddItem, testReqAddItemWithFlags = func() ([]byte, []byte) {
	reqAddItemLen := len(testReqAddItemHeader) + len(testPerlBytes)
	reqAddItem := make([]byte, len(testReqAddItemHeader), reqAddItemLen)
	reqAddItemWithFlags := make([]byte, reqAddItemLen, reqAddItemLen+4)

	copy(reqAddItem, testReqAddItemHeader)
	reqAddItem = append(reqAddItem, testPerlBytes...)
	copy(reqAddItemWithFlags, reqAddItem)
	reqAddItemWithFlags = append(reqAddItemWithFlags, 1, 0, 0, 0)

	return reqAddItem, reqAddItemWithFlags
}()

func testAddItem(t *testing.T, orig []byte, what string) ReqAddItem {
	var req ReqAddItem

	if _, err := req.UnmarshalIProto(orig); err != nil {
		t.Fatalf("decoding ReqAddItem %s flags failed: %v", what, err)
	}

	encoded, err := req.MarshalIProto(nil)
	if err != nil {
		t.Fatalf("encoding ReqAddItem %s flags failed: %v", what, err)
	}

	assert.Equalf(t, orig, encoded, "original and reencoded data (%s flags) differ", what)

	return req
}

func TestAddItemWithoutFlags(t *testing.T) {
	testAddItem(t, testReqAddItem, "without")
}

func TestAddItemWithFlags(t *testing.T) {
	req := testAddItem(t, testReqAddItemWithFlags, "with")
	assert.Equal(t, req.Flags, uint32(1), "flags are missing")
}
