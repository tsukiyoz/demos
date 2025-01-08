package smtp

import (
	"fmt"
	"net/smtp"
	"strings"
	"testing"
	"time"
)

func TestSMTPUsage(t *testing.T) {
}

func SendToMail(user, password, host, subject, date, body, mailtype, replyToAddress string, to, cc, bcc []string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	cc_address := strings.Join(cc, ";")
	bcc_address := strings.Join(bcc, ";")
	to_address := strings.Join(to, ";")
	msg := []byte("To: " + to_address + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\nDate: " + date + "\r\nReply-To: " + replyToAddress + "\r\nCc: " + cc_address + "\r\nBcc: " + bcc_address + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := MergeSlice(to, cc)
	send_to = MergeSlice(send_to, bcc)
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func main() {
	user := "控制台创建的发信地址"
	password := "控制台设置的SMTP密码"
	host := "smtpdm.aliyun.com:80"
	to := []string{"收件人地址", "收件人地址1"}
	cc := []string{"抄送地址", "抄送地址1"}
	bcc := []string{"密送地址", "密送地址1"}
	subject := "test Golang to sendmail"
	date := fmt.Sprintf("%s", time.Now().Format(time.RFC1123Z))
	mailtype := "html"
	replyToAddress := "a***@example.net"
	body := "<html><body><h3>Test send to email</h3></body></html>"
	fmt.Println("send email")
	err := SendToMail(user, password, host, subject, date, body, mailtype, replyToAddress, to, cc, bcc)
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}
}

func MergeSlice(s1 []string, s2 []string) []string {
	slice := make([]string, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}
