package socket

import (
	"encoding/binary"
	"testing"
)

func TestPackHeader(t *testing.T) {
	client := newClient(&Socket{})

	if client == nil {
		t.Fatalf("client create failed")
	}

	type sizeAndEOF struct {
		size   int
		eofs   bool
		hasErr bool
	}

	sizeAndEOFSlice := []sizeAndEOF{
		{
			size:   0,
			eofs:   true,
			hasErr: false,
		},
		{
			size:   0,
			eofs:   false,
			hasErr: false,
		},
		{
			size:   -100,
			eofs:   true,
			hasErr: true,
		},
		{
			size:   -100,
			eofs:   false,
			hasErr: true,
		},
		{
			size:   1276,
			eofs:   true,
			hasErr: false,
		},
		{
			size:   9527,
			eofs:   false,
			hasErr: false,
		},
		{
			size:   PACKET_LIMIT_SIZE,
			eofs:   true,
			hasErr: false,
		},
		{
			size:   PACKET_LIMIT_SIZE,
			eofs:   false,
			hasErr: false,
		},
		{
			size:   PACKET_LIMIT_SIZE + 1,
			eofs:   true,
			hasErr: true,
		},
		{
			size:   PACKET_LIMIT_SIZE + 1,
			eofs:   false,
			hasErr: true,
		},
		{
			size:   PACKET_LIMIT_SIZE - 1,
			eofs:   true,
			hasErr: false,
		},
		{
			size:   PACKET_LIMIT_SIZE - 1,
			eofs:   false,
			hasErr: false,
		},
	}

	// bigendian test
	for i, encoder := range []binary.ByteOrder{binary.BigEndian, binary.LittleEndian} {
		t.Logf("testing byteorder index %v", i)
		client.SetHeaderEncoder(encoder)

		for _, data := range sizeAndEOFSlice {
			header := []byte{0, 0}
			err := client.encodeHeader(data.size, data.eofs, header)

			if !data.hasErr && err != nil {
				t.Fatalf("encode data %+v failed err %v", data, err)
			}

			if err != nil {
				continue
			}

			n, eofs, err := client.decodeHeader(header)

			if err != nil {
				t.Fatalf("decode data %+v failed err %v,header %+v", data, err, header)
			}

			if n != data.size || eofs != data.eofs {
				t.Fatalf("decode data %+v failed err %v,header %+v result n %v eofs %v", data, err, header, n, eofs)
			}
		}
	}
}
