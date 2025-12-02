package server

import (
	"fmt"
	"io"
	"net"
	"strings"

	moex_contract_v1 "github.com/Mar1eena/trb_proto/gen/go/moex"
)

func Dealing(req *moex_contract_v1.DealingRequest) (*moex_contract_v1.DealingResponse, error) {

	logon := msgBuild(req.Header, req.Logon)

	dialer := &net.Dialer{}
	conn, err := dialer.Dial("tcp", req.Address)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(logon)); err != nil {
		return nil, fmt.Errorf("failed to send login: %w", err)
	}

	var fullResponse strings.Builder
	buffer := make([]byte, 16384)

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

func msgBuild(header string, body string) string {
	msges := "35=A" + "\x01" + "34=1" + "\x01" + header + "\x01" + body + "\x01" + "141=Y" + "\x01" + "1137=9" + "\x01"

	// Длина сообщения без первых двух полей BeginString и BodyLength
	bodyLength := len(msges)

	// Формируем заголовок без трейлера
	head := "8=FIXT.1.1" + "\x01" +
		fmt.Sprintf("9=%03d\x01", bodyLength)
	messageWithoutTrailer := head + msges

	// считаем контрольную сумму
	checkSum := 0
	for _, b := range messageWithoutTrailer {
		checkSum += int(b)
	}
	checkSum %= 256

	// Формируем конечное сообщение с заголовком, телом и трейлером
	message := messageWithoutTrailer + fmt.Sprintf("10=%03d\x01", checkSum)

	return message
}
