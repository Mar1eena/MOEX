package server

import (
	"fmt"
	"net"
	"strings"

	moex_contract_v1 "github.com/Mar1eena/trb_proto/gen/go/moex"
)

func Request(req *moex_contract_v1.Req) (string, error) {
	addr := "91.208.232.200:9229"

	msg, err := message(req, "A", body(req.Body))

	if err != nil {
		return err.Error(), err
	}

	// соединение
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Не удалось подключиться: %v\n", err)
		return err.Error(), err
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Не удалось прочитать ответ: %v\n", err)
		return err.Error(), err
	}

	// логинимся
	conn.Write([]byte(msg))
	fmt.Println("Сообщение отправлено")

	fmt.Printf("Ответ: %s\n", string(buffer[:n]))

	// получаем инструменты
	msg, err = message(req, "AE", instrument(req.Body))
	conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("Не удалось прочитать ответ: %v\n", err)
		return err.Error(), err
	}
	return string(buffer[:n]), nil
}

func message(req *moex_contract_v1.Req, msgtype string, body string) (string, error) {
	mes := "35=" + msgtype + "\x01" +
		req.Header.MsgSeqNum + "\x01" +
		req.Header.SendercompID + "\x01" +
		req.Header.TargetcompID + "\x01" +
		req.Header.SendingTime + "\x01" +
		body
	bodyLength := len(mes)

	header := req.Header.BeginString + "\x01" +
		fmt.Sprintf("9=%03d\x01", bodyLength)

	messageWithoutTrailer := header + mes

	checkSum := 0
	for _, b := range messageWithoutTrailer {
		checkSum += int(b)
	}
	checkSum %= 256

	message := messageWithoutTrailer + fmt.Sprintf("10=%03d\x01", checkSum)

	return message, nil
}

func body(req *moex_contract_v1.Body) string {
	if req == nil {
		return ""
	}

	var parts []string

	// Обработка Logon
	if req.Logon != nil {
		logon := req.Logon
		if logon.EncryptMethod != "" {
			parts = append(parts, logon.EncryptMethod)
		}
		if logon.Heartbtint != "" {
			parts = append(parts, logon.Heartbtint)
		}
		if logon.Password != "" {
			parts = append(parts, logon.Password)
		}
		if logon.ResetSeqNumFlag != "" {
			parts = append(parts, logon.ResetSeqNumFlag)
		}
		if logon.TestReqID != "" {
			parts = append(parts, logon.TestReqID)
		}
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, "\x01") + "\x01"
}

func instrument(req *moex_contract_v1.Body) string {
	for _, v := range req.Instrument {
		return v.Symbol + "\x01" +
			v.Product + "\x01" +
			v.SecurityType + "\x01"
	}

	return ""
}
