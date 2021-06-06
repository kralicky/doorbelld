package unifi

import "net/url"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type HeaderPacketType byte
type HeaderPayloadFormat byte
type HeaderCompression byte

const (
	PacketTypeActionFrame HeaderPacketType = 1
	PacketTypeDataFrame   HeaderPacketType = 2

	FormatJSONObject HeaderPayloadFormat = 1
	FormatUTF8String HeaderPayloadFormat = 2
	FormatNodeBuffer HeaderPayloadFormat = 3

	Uncompressed HeaderCompression = 0
	Compressed   HeaderCompression = 1
)

type PacketHeader struct {
	Type        HeaderPacketType
	Format      HeaderPayloadFormat
	Compression HeaderCompression
	Reserved    byte
	PayloadSize uint32
}

type ActionKind string

const (
	ActionKindAdd    ActionKind = "add"
	ActionKindUpdate ActionKind = "update"
)

type ModelKeyKind string

const (
	ModelKeyEvent  ModelKeyKind = "event"
	ModelKeyNVR    ModelKeyKind = "nvr"
	ModelKeyCamera ModelKeyKind = "camera"
	ModelKeyUser   ModelKeyKind = "user"
)

type ActionFrame struct {
	Action      ActionKind   `json:"action"`
	ModelKey    ModelKeyKind `json:"modelKey"`
	NewUpdateID string       `json:"newUpdateId"`
	Score       int          `json:"score"`
	Camera      string       `json:"camera"`
	ID          string       `json:"id"`
}

type EventDataFrame struct {
	Type     string `json:"type"`
	Start    int    `json:"start"`
	Score    int    `json:"score"`
	Camera   string `json:"camera"`
	ID       string `json:"id"`
	ModelKey string `json:"modelKey"`
}

type Config struct {
	Endpoint *url.URL
}

type UpdatePayload struct {
	Header          PacketHeader
	ActionFrame     ActionFrame
	SecondaryHeader PacketHeader
	DataFrame       []byte
}
