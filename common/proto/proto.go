package proto

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
)

/**
协议 [0 0 0 0, 0 0, 0 0, 0 0 0 0, 0 0 0 0, ...]
	包总长-	协议长度(16),版本,通讯代号,身份,消息体
*/
/*type Msg struct {
	Ver       uint16          `json:"ver"`
	Operation uint32          `json:"operation"`
	SeqId     uint32          `json:"seqId"`
	Body      json.RawMessage `json:"body"`
	Code      int32           `json:"code"`
}
*/

const (
	MaxBodySize = uint32(1 << 10)
)

const (
	PackSize      = 4
	HeaderSize    = 2
	VerSize       = 2
	OperationSize = 4
	SeqIdSize     = 4
	RawHeaderSize = PackSize + HeaderSize + VerSize + OperationSize + SeqIdSize
	MaxPackSize   = MaxBodySize + uint32(RawHeaderSize)
	// offset
	PackOffset      = 0
	HeaderOffset    = PackOffset + PackSize
	VerOffset       = HeaderOffset + HeaderSize
	OperationOffset = VerOffset + VerSize
	SeqIdOffset     = OperationOffset + OperationSize
	BodyOffset      = SeqIdOffset + SeqIdSize
)

const (
	// OpHandshake handshake
	OpHandshake = int32(0)
	// OpHandshakeReply handshake reply
	OpHandshakeReply = int32(1)

	// OpHeartbeat heartbeat
	OpHeartbeat = int32(2)
	// OpHeartbeatReply heartbeat reply
	OpHeartbeatReply = int32(3)

	// OpSendMsg send message.
	OpSendMsg = int32(4)
	// OpSendMsgReply  send message reply
	OpSendMsgReply = int32(5)

	// OpDisconnectReply disconnect reply
	OpDisconnectReply = int32(6)

	// OpAuth auth connnect
	OpAuth = int32(7)
	// OpAuthReply auth connect reply
	OpAuthReply = int32(8)

	// OpRaw raw message
	OpRaw = int32(9)

	// OpProtoReady proto ready
	OpProtoReady = int32(10)
	// OpProtoFinish proto finish
	OpProtoFinish = int32(11)

	OpContact      = int32(12)
	OpContactReply = int32(13)

	OpConversion = int32(14)
)

var (
	ProtoFinish = &Proto{Opr: OpProtoFinish}
	ProtoReady  = &Proto{Opr: OpProtoReady}

	ErrMsgPackLen   = errors.New("default Server codec pack length error")
	ErrMsgHeaderLen = errors.New("default Server codec header length error")
	ErrMsgNotCheck  = errors.New("connect not check")
)

func (p *Proto) ReadWebSocket(ws *websocket.Conn) (err error) {
	var (
		bodyLen   uint32
		headerLen uint16
		packLen   uint32
		allBuf    []byte
	)

	_, allBuf, err = ws.ReadMessage()
	if err != nil {
		return
	}
	fmt.Println(allBuf)
	if len(allBuf) < (RawHeaderSize) {
		return ErrMsgHeaderLen
	}

	packLen = binary.BigEndian.Uint32(allBuf[PackOffset:HeaderOffset])
	headerLen = binary.BigEndian.Uint16(allBuf[HeaderOffset:VerOffset])
	p.Ver = int32(binary.BigEndian.Uint16(allBuf[VerOffset:OperationOffset]))
	p.Opr = int32(binary.BigEndian.Uint32(allBuf[OperationOffset:SeqIdOffset]))
	p.Seq = int32(binary.BigEndian.Uint32(allBuf[SeqIdOffset:]))

	if packLen > MaxPackSize {
		return ErrMsgPackLen
	}

	if headerLen != RawHeaderSize {
		return ErrMsgHeaderLen
	}
	if bodyLen = packLen - uint32(headerLen); bodyLen > 0 {

		p.Body = allBuf[headerLen:packLen]
	} else {
		p.Body = nil
	}
	return nil
}

func (p *Proto) WriteWebSocket(ws *websocket.Conn) (err error) {
	var (
		buf     = make([]byte, RawHeaderSize)
		packLen = uint32(RawHeaderSize) + uint32(len(p.Body))
	)
	binary.BigEndian.PutUint32(buf, packLen)
	binary.BigEndian.PutUint16(buf[HeaderOffset:], RawHeaderSize)
	binary.BigEndian.PutUint16(buf[VerOffset:], uint16(p.Ver))
	binary.BigEndian.PutUint32(buf[OperationOffset:], uint32(p.Opr))
	binary.BigEndian.PutUint32(buf[SeqIdOffset:], uint32(p.Seq))

	buf = append(buf, p.Body...)

	err = ws.WriteMessage(websocket.BinaryMessage, buf)
	return
}
