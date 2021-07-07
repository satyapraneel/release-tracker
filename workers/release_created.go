package workers

type ReleaseCreateData struct {
	ProjectName string
	ReleaseName string
}

// func SendReleaseCreatedMail(release *models.Release, project *models.Project) {
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		mails.SendReleaseCreatedMail(release, project)
// 	}()
// }
