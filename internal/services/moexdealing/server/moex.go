package server

import (
	"fmt"
	"net"
	"time"

	moex_contract_v1 "github.com/Mar1eena/trb_proto/gen/go/moex"
)

func Dealing(req *moex_contract_v1.DealingRequest) (string, error) {

	addr := req.Address
	logon := msgBuild(req.Header, req.Logon, "A")
	// instrument := msgBuild(req.Header, "AE")

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	buffer := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// Отправляем запрос логина и читаем ответ
	if _, err := conn.Write([]byte(logon)); err != nil {
		return "", err
	}

	// Читаем ответ
	n, err := conn.Read(buffer)
	resp := string(buffer[:n])

	return resp, err
}

func msgBuild(header string, body string, msgtype string) string {
	// Добавляем тип сообщения и соединяем заголовок с телом, чтобы вычислить длину
	msges := "35=" + msgtype + "\x01" + header + "\x01" + body + "\x01"

	// Длина сообщения без первых двух полей BeginString и BodyLength
	bodyLength := len(msges)

	// Формируем заголовок без трейлера
	head := "8=FIX.4.4" + "\x01" +
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
