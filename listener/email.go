package listener

import (
	"fmt"
	"os"
	"net"
	"crypto/tls"
    "net/mail"
    "net/smtp"
    "strings"
	"github.com/matcornic/hermes"
	"github.com/marstau/smartcooly/config"
)

var (
	user = config.String("emailfrom")
	password = config.String("emailfrompass")
	host = config.String("smtphost")
)

// Email struct
type Email struct {
	Title     int64
	Content   string
	Type      int64
}


func init() {

	envSMTP := os.Getenv("MAIL_SMTP")
	if envSMTP != "" {
		host = envSMTP + os.Getenv("MAIL_SMTP_PORT")
		user = os.Getenv("MAIL_SENDER_ADDRESS")
		password = os.Getenv("MAIL_PASSWORD")
	}
}

func sendPlain(user, password, host, to, title, subject, body, mailtype string) error {
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

// dial using TLS/SSL
func dial(addr string) (*tls.Conn, error) {
    /*
        // TLS config
        tlsconfig := &tls.Config{
            // InsecureSkipVerify controls whether a client verifies the
            // server's certificate chain and host name.
            // If InsecureSkipVerify is true, TLS accepts any certificate
            // presented by the server and any host name in that certificate.
            // In this mode, TLS is susceptible to man-in-the-middle attacks.
            // This should be used only for testing.
            InsecureSkipVerify: false,
            // ServerName indicates the name of the server requested by the client
            // in order to support virtual hosting. ServerName is only set if the
            // client is using SNI (see
            // http://tools.ietf.org/html/rfc4366#section-3.1).
            // ServerName: host,
            // MinVersion contains the minimum SSL/TLS version that is acceptable.
            // If zero, then TLS 1.0 is taken as the minimum.
            MinVersion: tls.VersionSSL30,
            // MaxVersion contains the maximum SSL/TLS version that is acceptable.
            // If zero, then the maximum version supported by this package is used,
            // which is currently TLS 1.2.
            MaxVersion: tls.VersionSSL30,
        }
    */
    // Here is the key, you need to call tls.Dial instead of smtp.Dial
    // for smtp servers running on 465 that require an ssl connection
    // from the very beginning (no starttls)
    return tls.Dial("tcp", addr, nil)
}

// compose message according to "from, to, subject, body"
func composeMsg(from string, to string, subject string, body string) (message string) {
    // Setup headers
    headers := make(map[string]string)
    headers["From"] = from
    headers["To"] = to
    headers["Subject"] = subject
    // Setup message
    for k, v := range headers {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + body
    return
}

// send email over SSL
func sendSSL(to string, title string, subject string, body string) (err error) {
    localhost, _, _ := net.SplitHostPort(host)
    // get SSL connection
    conn, err := dial(host)
    if err != nil {
        return
    }
    // create new SMTP client
    smtpClient, err := smtp.NewClient(conn, localhost)
    if err != nil {
        return
    }
    // Set up authentication information.
    auth := smtp.PlainAuth("", user, password, localhost)
    // auth the smtp client
    err = smtpClient.Auth(auth)
    if err != nil {
        return
    }
    // set To && From address, note that from address must be same as authorization user.
    from := mail.Address{"", user}
    localto := mail.Address{"", to}
    err = smtpClient.Mail(from.Address)
    if err != nil {
        return
    }
    err = smtpClient.Rcpt(localto.Address)
    if err != nil {
        return
    }
    // Get the writer from SMTP client
    writer, err := smtpClient.Data()
    if err != nil {
        return
    }
    // compose message body
    message := composeMsg(from.String(), localto.String(), subject, body)
    // write message to recp
    _, err = writer.Write([]byte(message))
    if err != nil {
        return
    }
    // close the writer
    err = writer.Close()
    if err != nil {
        return
    }
    // Quit sends the QUIT command and closes the connection to the server.
    smtpClient.Quit()
    return nil
}

func (e Email) SendEmail(to string, title string, bodyTitle string) {
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

	err = sendSSL(to, title, emailText, emailBody)
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}

}