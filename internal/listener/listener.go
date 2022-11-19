package listener

func Init() {
	go MailLinter()
	go BlogLinter()
}
