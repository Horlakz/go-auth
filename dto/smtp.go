package dto

type SendEmailDto struct {
	To, Subject, Template string
	Variables             interface{}
}
