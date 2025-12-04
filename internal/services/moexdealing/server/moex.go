package server

import (
	"fmt"
	"io"
	"net"
	"strings"

	moex_contract_v1 "github.com/Mar1eena/trb_proto/gen/go/moex"
)

func Dealing(req *moex_contract_v1.DealingRequest) (*moex_contract_v1.DealingResponse, error) {

	var fullResponse strings.Builder
	buffer := make([]byte, 16384)

	logon := msgBuild(req.Header, req.Logon, "A", "1")

	dialer := &net.Dialer{}
	conn, err := dialer.Dial("tcp", req.Address)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(logon)); err != nil {
		return nil, fmt.Errorf("failed to send login: %w", err)
	}

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			return nil, fmt.Errorf("failed: %w, response: %s", err, string(buffer[:n]))
		}

		if n > 0 {
			fullResponse.Write(buffer[:n])
		}
	}
	return &moex_contract_v1.DealingResponse{Response: fullResponse.String()}, nil
}

func msgBuild(header string, body string, msgType string, msgSeq string) string {
	msges := fmt.Sprintf("35=%s\x0134=%s\x01%s\x01%s\x0198=0\x01141=Y\x011137=9\x01",
		msgType, msgSeq, header, body)

	// Длина сообщения без первых двух полей BeginString и BodyLength
	bodyLength := len(msges)

	// Формируем заголовок без трейлера
	head := fmt.Sprintf("8=FIXT.1.1\x019=%03d\x01", bodyLength)
	messageWithoutTrailer := head + msges

	// считаем контрольную сумму
	checkSum := 0
	for _, b := range messageWithoutTrailer {
		checkSum += int(b)
	}
	checkSum %= 256

	// Формируем конечное сообщение с заголовком, телом и трейлером
	message := fmt.Sprintf("%v10=%03d\x01", messageWithoutTrailer, checkSum)

	return message
}
