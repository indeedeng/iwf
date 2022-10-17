package interpreter

import "github.com/cadence-oss/iwf-server/gen/iwfidl"

type InterStateChannel struct {
	// key is channel name
	receivedData map[string][]*iwfidl.EncodedObject
}

func NewInterStateChannel() *InterStateChannel {
	return &InterStateChannel{
		receivedData: map[string][]*iwfidl.EncodedObject{},
	}
}

func (i *InterStateChannel) HasData(channelName string) bool {
	l := i.receivedData[channelName]
	return len(l) > 0
}

func (i *InterStateChannel) ProcessPublishing(publishes []iwfidl.InterStateChannelPublishing) {
	for _, pub := range publishes {
		i.receive(pub.ChannelName, pub.Value)
	}
}

func (i *InterStateChannel) receive(channelName string, data *iwfidl.EncodedObject) {
	l := i.receivedData[channelName]
	l = append(l, data)
	i.receivedData[channelName] = l
}

func (i *InterStateChannel) Retrieve(channelName string) *iwfidl.EncodedObject {
	l := i.receivedData[channelName]
	if len(l) <= 0 {
		panic("critical bug, this shouldn't happen")
	}
	data := l[0]
	l = l[1:]
	i.receivedData[channelName] = l
	return data
}
