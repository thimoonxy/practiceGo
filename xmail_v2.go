package main

import (
	"bufio"
	"flag"
	"io"
	"os"

	"gopkg.in/gomail.v2"
)

func main() {
	body := flag.String("b", "mail_body", "mail body strings")
	subject := flag.String("s", "mail_subject", "mail subject strings")
	att := flag.String("f", "mail_attachment", "mail attachment filepath")
	has_stdin := flag.Bool("i", false, "mail_body pipe from stdin. ")
	flag.Parse()
	msg := gomail.NewMessage()
	msg.SetHeader("From", "pwrdbugs@sina.com")
	msg.SetHeader("To", "22749752@qq.com", "pwrdbugs@sina.com", "405796346@qq.com", "airfly555@qq.com", "553668853@qq.com")
	//	msg.SetAddressHeader("Cc", "xieyin@pwrd.com", "Simon")
	if *subject != "mail_subject" {
		msg.SetHeader("Subject", *subject)
	} else {
		msg.SetHeader("Subject", "[CrossOut]Alarm")
	}
	var stdinstring string
	if *has_stdin {
		stdinstring = readstdin()
	}

	if len(stdinstring) != 0 {
		msg.SetBody("text/html", stdinstring)
	} else if *body != "mail_body" {
		msg.SetBody("text/html", *body)
	} else {
		msg.SetBody("text/html", "Failed to read <b>mail_body</b> !")

	}

	if *att != "mail_attachment" {
		msg.Attach(*att)
	}

	mailer := gomail.NewDialer("smtp.sina.com", 25, "pwrdbugs", "xxx")
	if err := mailer.DialAndSend(msg); err != nil {
		panic(err)
	}

}

func readstdin() string {
	reader := bufio.NewReader(os.Stdin)
	var (
		result, line string
		err          error
	)
	for {
		line, err = reader.ReadString('\r')
		result = result + line
		if err == io.EOF {
			break
		}

	}
	return result
}
