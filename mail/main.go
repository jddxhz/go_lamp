package mail

import (
	// "encoding/base64"
	"fmt"
	// "io/ioutil"
	"net/smtp"
)

func main() {
	SendMailByNetSMTP(`你好啊`)
}

//SendMailByNetSMTP 发送邮件
func SendMailByNetSMTP(str string) {
	auth := smtp.PlainAuth("", "writecycle@163.com", "zhang524399", "smtp.163.com")

	to := []string{"writecycle@qq.com"}
	// image, _ := ioutil.ReadFile("c://1.jpg")
	// imageBase64 := base64.StdEncoding.EncodeToString(image)
	msg := []byte("from:writecycle@163.com\r\n" +
		"to: writecycle@qq.com\r\n" +
		"Subject: hello,subject!\r\n" +
		"Content-Type:multipart/mixed;boundary=a\r\n" +
		"Mime-Version:1.0\r\n" +
		"\r\n" +
		"--a\r\n" +
		"Content-type:text/plain;charset=utf-8\r\n" +
		"Content-Transfer-Encoding:quoted-printable\r\n" +
		"\r\n" +
		str + "\r\n" +
		// "--a\r\n" +
		// "Content-type:image/jpg;name=1.jpg\r\n" +
		// "Content-Transfer-Encoding:base64\r\n" +
		// "\r\n" +
		// imageBase64 + "\r\n" +
		"--a--\r\n")
	err := smtp.SendMail("smtp.163.com:25", auth, "writecycle@163.com", to, msg)
	if err != nil {
		fmt.Println(err)
	}
}
