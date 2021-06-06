package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strings"
)

func TriggerMail(reviewers string,releaseName string,projectName string)  {
	// Sender data.
	from := os.Getenv("SENDER_EMAIL")
	password := os.Getenv("SENDER_PASS")

	// Receiver email address.
	reviews := strings.Split(reviewers,",")
	to := reviews
	// smtp server configuration.
	smtpHost := os.Getenv("SMTPHOST")
	smtpPort := os.Getenv("SMTPPORT")
	auth := smtp.PlainAuth("", from, password, smtpHost)

	emailTemplate, errs := template.ParseFiles("cmd/email_template.html")
	if errs != nil {
		log.Printf("template parse : %v",errs)
	}
	var body bytes.Buffer
	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("Subject: Jira Label for "+projectName+"!\r\n", headers)))
	emailTemplate.Execute(&body, struct {
		ProjectName string
		ReleaseName string
	}{
		ProjectName: projectName,
		ReleaseName: releaseName,
	})

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

