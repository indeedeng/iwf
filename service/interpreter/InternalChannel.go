package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service/common/ptr"
)

type InternalChannel struct {
	// key is channel name
	receivedData map[string][]*iwfidl.EncodedObject
}

func NewInternalChannel() *InternalChannel {
	return &InternalChannel{
		receivedData: map[string][]*iwfidl.EncodedObject{},
	}
}

func RebuildInternalChannel(refill map[string][]*iwfidl.EncodedObject) *InternalChannel {
	return &InternalChannel{
		receivedData: refill,
	}
}

func (i *InternalChannel) GetAllReceived() map[string][]*iwfidl.EncodedObject {
	return i.receivedData
}

func (i *InternalChannel) GetInfos() map[string]iwfidl.ChannelInfo {
	infos := make(map[string]iwfidl.ChannelInfo, len(i.receivedData))
	for name, l := range i.receivedData {
		infos[name] = iwfidl.ChannelInfo{
			Size: ptr.Any(int32(len(l))),
		}
	}
	return infos
}

func (i *InternalChannel) HasData(channelName string) bool {
	l := i.receivedData[channelName]
	return len(l) > 0
}

func (i *InternalChannel) HasAtLeastN(channelName string, n int) bool {
	l := i.receivedData[channelName]
	return len(l) >= n
}

func (i *InternalChannel) Size(channelName string) int {
	return len(i.receivedData[channelName])
}

func (i *InternalChannel) ProcessPublishing(publishes []iwfidl.InterStateChannelPublishing) {
	for _, pub := range publishes {
		i.receive(pub.ChannelName, pub.Value)
	}
}

func (i *InternalChannel) receive(channelName string, data *iwfidl.EncodedObject) {
	l := i.receivedData[channelName]
	l = append(l, data)
	i.receivedData[channelName] = l
}

func (i *InternalChannel) Retrieve(channelName string) *iwfidl.EncodedObject {
	l := i.receivedData[channelName]
	if len(l) <= 0 {
		panic("critical bug, this shouldn't happen")
	}
	data := l[0]
	l = l[1:]
	if len(l) == 0 {
		delete(i.receivedData, channelName)
	} else {
		i.receivedData[channelName] = l
	}

	return data
}

// RetrieveUpToN atomically retrieves up to n messages from the channel.
// It consumes min(n, available) messages.
func (i *InternalChannel) RetrieveUpToN(channelName string, n int) []*iwfidl.EncodedObject {
	l := i.receivedData[channelName]
	if len(l) == 0 {
		return []*iwfidl.EncodedObject{}
	}
	count := n
	if count > len(l) {
		count = len(l)
	}
	data := make([]*iwfidl.EncodedObject, count)
	copy(data, l[:count])
	l = l[count:]
	if len(l) == 0 {
		delete(i.receivedData, channelName)
	} else {
		i.receivedData[channelName] = l
	}
	return data
}

// RetrieveAtLeastUpToN atomically retrieves up to n messages only when at least
// atLeast messages are currently available.
func (i *InternalChannel) RetrieveAtLeastUpToN(channelName string, atLeast, n int) ([]*iwfidl.EncodedObject, bool) {
	if !i.HasAtLeastN(channelName, atLeast) {
		return nil, false
	}
	return i.RetrieveUpToN(channelName, n), true
}
