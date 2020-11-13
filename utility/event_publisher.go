package utility

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
)

func PublishEmailVerificationEvent(payload string) *RestError {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	svc := sqs.New(sess)
	_, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(payload),
		QueueUrl:    aws.String(os.Getenv("EMAIL_VERIFICATION_QUEUE")),
	})
	if err != nil {
		return NewInternalServerError("Could not publish to SQS", &err)
	}
	return nil
}
