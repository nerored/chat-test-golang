/*
	打包/解包数据到protobuf
*/
package utils

import (
	"encoding/binary"

	"github.com/nerored/chat-test-golang/cli/log"
	"github.com/nerored/chat-test-golang/message"
	"google.golang.org/protobuf/proto"
)

func TranslateMsg(data []byte) (api message.ChatServiceAPI, msg []byte) {
	if len(data) < 4 {
		return
	}

	return message.ChatServiceAPI(binary.BigEndian.Uint32(data[:4])), data[4:]
}

func PackMsg(api message.ChatServiceAPI, msg proto.Message) (data []byte) {
	packet, err := proto.Marshal(msg)

	if err != nil {
		log.Erro("packMsg failed")
		return
	}

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(api))
	return append(header, packet...)
}

func UnpackMsg(data []byte, msg proto.Message) (ok bool) {
	err := proto.Unmarshal(data, msg)

	if err != nil {
		log.Erro("unpackMsg err %v", err)
		return
	}

	return true
}
