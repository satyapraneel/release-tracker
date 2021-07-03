package mails

import (
	"fmt"
	"log"

	"github.com/release-trackers/gin/models"
)

type ReMinderData struct {
	ProjectName  string
	ReminderType string
}

func SendReminderMailTest(dls *models.DLS, mailData *MailData, mailTemplate string) {
	// tagetDate := &release.TargetDate
	dlsList := []string{dls.Email}

	mail := NewMail(dlsList, mailData.Subject, "", "")
	errs := mail.ParseTemplate(mailTemplate, mailData)
	if errs != nil {
		log.Printf("template parse : %v", errs)
	}
	ok, _ := mail.SendEmail()
	fmt.Println(ok)
}

//func SendReminderMail(project *models.Project, dls *models.DLS, reminderType string) {
//	// tagetDate := &release.TargetDate
//	dlsList := []string{dls.Email}
//
//	subject := "Reminder mail for " + reminderType + "!"
//	mail := NewMail(dlsList, subject, "", "")
//	reminderData := &ReMinderData{ProjectName: project.Name, ReminderType: reminderType}
//	errs := mail.ParseTemplate("/ui/html/mails/reminder.html", reminderData)
//	if errs != nil {
//		log.Printf("template parse : %v", errs)
//	}
//	ok, _ := mail.SendEmail()
//	fmt.Println(ok)
//}
