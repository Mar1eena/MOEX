package server

import (
	"fmt"
	"net"

	moex_contract_v1 "github.com/Mar1eena/trb_proto/gen/go/moex"
)

func Logon(req *moex_contract_v1.Logonrequest) (string, error) {
	addr := "91.208.232.200:9229"
	body := req.Header.Msgtype +
		"\x01" + req.Header.Sendercompid +
		"\x01" + req.Header.Targetcompid +
		"\x01" + req.Header.Msgseqnum +
		"\x01" + req.Header.Sendingtime +
		"\x01" + req.Logon.Encryptmethod +
		"\x01" + req.Logon.Heartbtint +
		"\x01" + req.Logon.Password +
		"\x01"
	bodyLength := len(body)

	header := req.Header.Beginstring + "\x01" + fmt.Sprintf("9=%03d\x01", bodyLength)

	checkSum := 0
	for _, b := range header + body {
		checkSum += int(b)
	}
	checkSum %= 256

	message := header + body + fmt.Sprintf("10=%03d\x01", checkSum)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Не удалось подключиться: %v\n", err)
		return "", err
	}
	defer conn.Close()

	conn.Write([]byte(message))
	fmt.Println("Сообщение отправлено")

	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	fmt.Printf("Ответ: %s\n", string(buffer[:n]))
	return string(buffer[:n]), nil
}
