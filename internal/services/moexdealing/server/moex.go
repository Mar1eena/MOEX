package server

import (
	"fmt"
	"net"
	"time"

	moex_contract_v1 "github.com/Mar1eena/trb_proto/gen/go/moex"
)

func Dealing(req *moex_contract_v1.DealingRequest) (*moex_contract_v1.DealingResponse, error) {
	addr := req.Address
	logon := msgBuild(req.Header, req.Logon, "A")
	instrument := msgBuild(req.Header, req.Instrument, "AE")

	// Устанавливаем соединение с таймаутом
	dialer := &net.Dialer{Timeout: 30 * time.Second}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return &moex_contract_v1.DealingResponse{Logon: "58=Не удалось установить соединение"}, nil
	}
	defer conn.Close()

	// Устанавливаем дедлайны для операций чтения/записи
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// 1. Отправляем запрос логина
	if _, err := conn.Write([]byte(logon)); err != nil {
		return &moex_contract_v1.DealingResponse{Logon: "58=Не удалось отправить логин"}, nil
	}

	// 2. Читаем ответ на логин
	logonBuffer := make([]byte, 4096)
	n, err := conn.Read(logonBuffer)
	respLogon := string(logonBuffer[:n])
	if err != nil {
		return &moex_contract_v1.DealingResponse{Logon: fmt.Sprintf("58=Не удалось  получить логин: %v, error: %v", logon, err)}, nil
	}

	// 3. Отправляем запрос инструмента
	if _, err := conn.Write([]byte(instrument)); err != nil {
		return &moex_contract_v1.DealingResponse{Logon: respLogon, Instrument: "58=Не удалось отправить инструмент"}, nil
	}

	// 4. Читаем ответ на инструмент
	instrumentBuffer := make([]byte, 4096)
	n, err = conn.Read(instrumentBuffer)
	respInstrument := string(instrumentBuffer[:n])
	if err != nil {
		return &moex_contract_v1.DealingResponse{Logon: respLogon, Instrument: fmt.Sprintf("58=Не удалось  получить инструмент: %v", err)}, nil
	}

	return &moex_contract_v1.DealingResponse{Logon: respLogon, Instrument: respInstrument}, nil
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
