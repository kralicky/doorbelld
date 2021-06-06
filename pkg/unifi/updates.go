package unifi

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type UpdateListener struct {
	C    chan UpdatePayload
	conn *websocket.Conn
}

func readPayloadSection(reader io.Reader) (PacketHeader, []byte, error) {
	header := PacketHeader{}
	if err := binary.Read(reader, binary.BigEndian, &header); err != nil {
		return PacketHeader{}, nil, err
	}

	jsonData := new(bytes.Buffer)
	if _, err := io.CopyN(jsonData, reader, int64(header.PayloadSize)); err != nil {
		return PacketHeader{}, nil, err
	}
	if header.Compression == Compressed {
		uncompressed := new(bytes.Buffer)
		zr, _ := zlib.NewReader(jsonData)
		if _, err := io.Copy(uncompressed, zr); err != nil {
			return PacketHeader{}, nil, err
		}
		jsonData = uncompressed
	}

	return header, jsonData.Bytes(), nil
}

func NewUpdateListener(ctx context.Context, cfg Config, tokenCookie *http.Cookie) (*UpdateListener, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   cfg.Endpoint.Host,
		Path:   "/proxy/protect/ws/updates",
	}
	log.Info("Connecting to update socket")
	header := http.Header{}
	header.Add("Cookie", tokenCookie.String())
	c, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), header)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("Connected")
	payloadCh := make(chan UpdatePayload, 100)
	listenerCtx, ca := context.WithCancel(ctx)
	go func() {
		defer c.Close()
		defer close(payloadCh)
		<-listenerCtx.Done()
	}()
	go func() {
		defer ca()
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Error(err)
				return
			}
			reader := bytes.NewReader(msg)
			header1, data1, err := readPayloadSection(reader)
			if err != nil {
				log.Error(err)
				return
			}
			af := ActionFrame{}
			if err := json.Unmarshal(data1, &af); err != nil {
				log.Error("Invalid message received")
				return
			}
			header2, data2, err := readPayloadSection(reader)
			if err != nil {
				log.Error(err)
				return
			}
			log.Debug("New update received")
			payloadCh <- UpdatePayload{
				Header:          header1,
				ActionFrame:     af,
				SecondaryHeader: header2,
				DataFrame:       data2,
			}
		}
	}()
	return &UpdateListener{
		C:    payloadCh,
		conn: c,
	}, nil
}
