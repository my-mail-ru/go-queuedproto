package addons

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatedAt_AddonID(t *testing.T) {
	var ca CreatedAt
	assert.Equal(t, CreatedAtID, ca.AddonID())
}

func TestCreatedAt_GetCreatedAt(t *testing.T) {
	ca := CreatedAt(1234567890)
	assert.Equal(t, ca, ca.GetCreatedAt())
}

func TestCreatedAt_MarshalAddon(t *testing.T) {
	ca := CreatedAt(0x0102030405060708)

	data, err := ca.MarshalAddon()
	require.NoError(t, err)
	assert.Equal(t, []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}, data)
}

func TestCreatedAt_UnmarshalAddon(t *testing.T) {
	var ca CreatedAt

	err := ca.UnmarshalAddon([]byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01})
	require.NoError(t, err)
	assert.Equal(t, CreatedAt(0x0102030405060708), ca)
}

func TestCreatedAt_Unmarshal_InvalidLength(t *testing.T) {
	var ca CreatedAt

	err := ca.UnmarshalAddon([]byte{0x01, 0x02, 0x03, 0x04})
	require.Error(t, err)
	assert.Equal(t, CreatedAt(0), ca)
}

func TestCreatedAt_Roundtrip(t *testing.T) {
	original := CreatedAt(1712140800000000000)

	data, err := original.MarshalAddon()
	require.NoError(t, err)

	var decoded CreatedAt
	err = decoded.UnmarshalAddon(data)
	require.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestCreatedAt_Roundtrip_Negative(t *testing.T) {
	original := CreatedAt(-1)

	data, err := original.MarshalAddon()
	require.NoError(t, err)

	var decoded CreatedAt
	err = decoded.UnmarshalAddon(data)
	require.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestCreatedAt_Build_FillsZero(t *testing.T) {
	var ca CreatedAt

	err := ca.Build()
	require.NoError(t, err)
	assert.NotEqual(t, CreatedAt(0), ca)
}

func TestCreatedAt_Build_KeepsExisting(t *testing.T) {
	ca := CreatedAt(1234567890)

	err := ca.Build()
	require.NoError(t, err)
	assert.Equal(t, CreatedAt(1234567890), ca)
}
