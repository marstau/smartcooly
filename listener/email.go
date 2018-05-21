package listener

import (
	"fmt"
	"net/smtp"
	"strings"
	"github.com/matcornic/hermes"
	"github.com/marstau/samaritan/config"
)

var (
	user = config.String("emailfrom")
	password = config.String("emailfrompass")
	to = config.String("emailto")
	host = "smtp.qq.com:25"
	HOSTS = map[string]string {
		"qq.com" : "smtp.qq.com:25",
		"hotmail.com" : "smtp.live.com:587",
		"gmail.com" : "smtp.gmail.com:587",
		"163.com" : "smtp.163.com:25",
		"126.com" : "smtp.126.com:25",
		"sina.com" : "smtp.sina.com.cn:25",
		"sina.cn" : "smtp.sina.com:25",
	}
)

// Email struct
type Email struct {
	Title     int64
	Content   string
	Type      int64
}


func init() {
	hostStart := strings.Index(user, "@")
	host = HOSTS[strings.ToLower(user[hostStart+1:])]
}

func send(user, password, host, to, title, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + title + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func (e Email) SendEmail(title string, bodyTitle string) {
	// Configure hermes by setting a theme and your product info
	h := hermes.Hermes{
	    // Optional Theme
	    // Theme: new(Default) 
	    Product: hermes.Product{
	        // Appears in header & footer of e-mails
	        Name: "Hermes",
	        Link: "https://example-hermes.com/",
	        // Optional product logo
	        Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
	    },
	}

	email := hermes.Email{
		Body: hermes.Body{
			Title: bodyTitle,
			Name: "Jon Snow",
			Intros: []string{
				"Your order has been processed successfully.",
			},
			Table: hermes.Table{
				Data: [][]hermes.Entry{
					{
						{Key: "Item", Value: "Golang"},
						{Key: "Description", Value: "Open source programming language that makes it easy to build simple, reliable, and efficient software"},
						{Key: "Price", Value: "$10.99"},
					},
					{
						{Key: "Item", Value: "Hermes"},
						{Key: "Description", Value: "Programmatically create beautiful e-mails using Golang."},
						{Key: "Price", Value: "$1.99"},
					},
				},
				Columns: hermes.Columns{
					CustomWidth: map[string]string{
						"Item":  "20%",
						"Price": "15%",
					},
					CustomAlignment: map[string]string{
						"Price": "right",
					},
				},
			},
			Actions: []hermes.Action{
				{
					Instructions: "You can check the status of your order and more in your dashboard:",
					Button: hermes.Button{
						Text: "Go to Dashboard",
						Link: "https://hermes-example.com/dashboard",
					},
				},
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
	    panic(err) // Tip: Handle error with something else than a panic ;)
	}

	// Generate the plaintext version of the e-mail (for clients that do not support xHTML)
	emailText, err := h.GeneratePlainText(email)
	if err != nil {
	    panic(err) // Tip: Handle error with something else than a panic ;)
	}

	err = send(user, password, host, to, title, emailText, emailBody, "html")
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}

}