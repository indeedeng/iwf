package interpreter

import (
	"testing"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/stretchr/testify/assert"
)

func TestInternalChannelRetrieveUpToN(t *testing.T) {
	ch := NewInternalChannel()
	ch.ProcessPublishing([]iwfidl.InterStateChannelPublishing{
		{ChannelName: "ch", Value: testEncodedObject("1")},
		{ChannelName: "ch", Value: testEncodedObject("2")},
		{ChannelName: "ch", Value: testEncodedObject("3")},
	})

	assert.True(t, ch.HasAtLeastN("ch", 2))
	assert.False(t, ch.HasAtLeastN("ch", 4))
	assert.Equal(t, 3, ch.Size("ch"))

	values := ch.RetrieveUpToN("ch", 2)
	assert.Len(t, values, 2)
	assert.Equal(t, "1", values[0].GetData())
	assert.Equal(t, "2", values[1].GetData())
	assert.Equal(t, 1, ch.Size("ch"))

	values = ch.RetrieveUpToN("ch", 10)
	assert.Len(t, values, 1)
	assert.Equal(t, "3", values[0].GetData())
	assert.False(t, ch.HasData("ch"))
}

func TestInternalChannelRetrieveUpToNEmpty(t *testing.T) {
	ch := NewInternalChannel()

	values := ch.RetrieveUpToN("missing", 3)

	assert.Empty(t, values)
	assert.False(t, ch.HasData("missing"))
}

func TestInternalChannelRetrieveAtLeastUpToN(t *testing.T) {
	ch := NewInternalChannel()
	ch.ProcessPublishing([]iwfidl.InterStateChannelPublishing{
		{ChannelName: "ch", Value: testEncodedObject("1")},
	})

	values, ok := ch.RetrieveAtLeastUpToN("ch", 2, 2)
	assert.False(t, ok)
	assert.Nil(t, values)
	assert.Equal(t, 1, ch.Size("ch"))

	values, ok = ch.RetrieveAtLeastUpToN("ch", 1, 2)
	assert.True(t, ok)
	assert.Len(t, values, 1)
	assert.Equal(t, "1", values[0].GetData())
	assert.False(t, ch.HasData("ch"))
}

func testEncodedObject(data string) *iwfidl.EncodedObject {
	return &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString(data),
	}
}
