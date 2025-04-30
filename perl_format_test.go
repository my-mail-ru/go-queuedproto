package queuedproto

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/my-mail-ru/go-queuedproto/addons"
)

var testPerlBytes = []byte{
	0x62, 0x02, 0x01, 0x00, 0x04, 0xA1, 0xE0, 0x14, 0x62, 0x04, 0x00, 0x10,
	0x0F, 0x55, 0x4E, 0x4B, 0x4E, 0x4F, 0x57, 0x4E, 0x2E, 0x55, 0x4E, 0x4B, 0x4E, 0x4F, 0x57, 0x4E,
	0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
}

//adv:iproto:
type testPerlData struct {
	ID    uint32
	Value uint32
}

func TestPerlDataDecode(t *testing.T) {
	var (
		decoded  PerlData
		event    testPerlData
		producer addons.Producer
	)

	if _, err := decoded.UnmarshalIProto(testPerlBytes); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, uint8(2), decoded.Version, "version doen't match")
	assert.Len(t, decoded.Addons.List, 2, "addons count doesn't match")
	assert.Equal(t, addons.CreationTimeID, decoded.Addons.List[0].ID, "addon #1 ID doesn't match")
	assert.Equal(t, addons.ProducerID, decoded.Addons.List[1].ID, "addon #2 ID doesn't match")
	assert.Len(t, decoded.Extensions.Data, 0, "unexpected extensions data")

	if err := producer.UnmarshalAddon(decoded.Addons.List[1].Data); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "UNKNOWN.UNKNOWN", string(producer), "addon #2 producer doesn't match")

	if _, err := event.UnmarshalIProto(decoded.Data); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, uint32(1), event.ID, "ID doesn't match")
	assert.Equal(t, uint32(2), event.Value, "Value doesn't match")
}

func TestPerlDataEncode(t *testing.T) {
	var decoded PerlData

	if _, err := decoded.UnmarshalIProto(testPerlBytes); err != nil {
		t.Fatal(err)
	}

	encoded, err := decoded.MarshalIProto(nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testPerlBytes, encoded, "encoded data doesn't match")
}
