package server

import (
	"fmt"
	"net"
	"strings"
	"time"

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

	instrumentBuffer := make([]byte, 1600000)
	for {
		d, err := conn.Read(instrumentBuffer)
		response := string(instrumentBuffer[:d])
		switch {
		case strings.Contains(response, "35=AE"):
			// return &moex_contract_v1.DealingResponse{Response: response}, nil
		case strings.Contains(response, "35=5"):
			return nil, err
		case err != nil:
			return nil, err
		}
	}

}

func msgBuild(header string, body string) string {
	loc := time.FixedZone("UTC-0", 0)

	// Получаем текущее время в этой локации
	now := time.Now().In(loc)
	// Добавляем тип сообщения и соединяем заголовок с телом, чтобы вычислить длину
	msges := header + "\x01" + "52=" + now.Format("00000000-00:00:00.000") + "\x01" + body + "\x01"

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
