package temporal

import (
	"testing"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
)

const testEncryptedEncoding = "binary/encrypted"

// testEncryptionCodec mirrors the relevant behavior of a real encrypting PayloadCodec
// (e.g. iwf-service's): Encode tags payloads with the "binary/encrypted" encoding, and
// Decode only unwraps payloads carrying that encoding, passing everything else through
// unchanged. The transform is a reversible marshal/unmarshal rather than real AES, which
// is sufficient to exercise the memo decode path.
type testEncryptionCodec struct{}

func (testEncryptionCodec) Encode(payloads []*commonpb.Payload) ([]*commonpb.Payload, error) {
	out := make([]*commonpb.Payload, len(payloads))
	for i, p := range payloads {
		b, err := p.Marshal()
		if err != nil {
			return nil, err
		}
		out[i] = &commonpb.Payload{
			Metadata: map[string][]byte{converter.MetadataEncoding: []byte(testEncryptedEncoding)},
			Data:     b,
		}
	}
	return out, nil
}

func (testEncryptionCodec) Decode(payloads []*commonpb.Payload) ([]*commonpb.Payload, error) {
	out := make([]*commonpb.Payload, len(payloads))
	for i, p := range payloads {
		if string(p.Metadata[converter.MetadataEncoding]) != testEncryptedEncoding {
			out[i] = p // passthrough, as a real codec does for encodings it didn't produce
			continue
		}
		var inner commonpb.Payload
		if err := inner.Unmarshal(p.Data); err != nil {
			return nil, err
		}
		out[i] = &inner
	}
	return out, nil
}

// TestGetMemoAndDecryptIfNeeded verifies that encrypted memos decode correctly in both
// formats iwf may encounter: the new format (newer Temporal SDKs apply the configured
// DataConverter, incl. its codec, to memos — sdk-go #1045) and the legacy format (older
// SDKs encoded memos with the default converter). Both must round-trip to the original
// value so existing workflows remain readable across the SDK upgrade.
func TestGetMemoAndDecryptIfNeeded(t *testing.T) {
	cryptoConverter := converter.NewCodecDataConverter(converter.GetDefaultDataConverter(), testEncryptionCodec{})
	client := &temporalClient{dataConverter: cryptoConverter, memoEncryption: true}

	value := iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("TestValue"),
	}

	// iwf pre-encrypts the memo value before handing it to the SDK.
	innerPayload, err := cryptoConverter.ToPayload(value)
	require.NoError(t, err)

	t.Run("new SDK format - data converter applied to memo", func(t *testing.T) {
		// Newer SDKs encode the memo value with the configured (encrypting) converter.
		memoField, err := cryptoConverter.ToPayload(innerPayload)
		require.NoError(t, err)

		out, err := client.getMemoAndDecryptIfNeeded(&commonpb.Memo{
			Fields: map[string]*commonpb.Payload{"TestKey": memoField},
		})
		require.NoError(t, err)
		assert.Equal(t, value, out["TestKey"])
	})

	t.Run("legacy format - default converter applied to memo", func(t *testing.T) {
		// Older SDKs encoded the memo value with the default converter.
		memoField, err := converter.GetDefaultDataConverter().ToPayload(innerPayload)
		require.NoError(t, err)

		out, err := client.getMemoAndDecryptIfNeeded(&commonpb.Memo{
			Fields: map[string]*commonpb.Payload{"TestKey": memoField},
		})
		require.NoError(t, err)
		assert.Equal(t, value, out["TestKey"])
	})
}
