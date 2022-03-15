package emailer

import (
	"fmt"
	"github.com/OpenCal-FYDP/AsyncCalendarOptimizer/internal/storer"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"strings"
)

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

type Emailer struct {
	client *ses.SES
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func New() (*Emailer, error) {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return nil, err
	}

	// Create an SES session.
	svc := ses.New(sess)

	return &Emailer{
		client: svc,
	}, nil
}

func (e *Emailer) SendConfirmationEmail(owner string, attendees []string, data *storer.EventData) error {
	const (
		// This address must be verified with Amazon SES.
		Sender = "notification@opencal.link"

		// Replace recipient@example.com with a "To" address. If your account
		// is still in the sandbox, this address must be verified.
		//Recipient = "jspsun@gmail.com"

		// Specify a configuration set. To use a configuration
		// set, comment the next line and line 92.
		//ConfigurationSet = "ConfigSet"

		// The subject line for the email.
		Subject = "You have an OpenCal meeting soon!"

		//The email body for recipients with non-HTML email clients.
		TextBody = "Email unable to render without html support"

		// The character encoding for the email.
		CharSet = "UTF-8"
	)
	if owner != "" && !contains(attendees, owner) {
		attendees = append(attendees, owner)
	}

	attendees = removeDuplicateStr(attendees)

	for _, attendee := range attendees {
		// Assemble the email.
		input := &ses.SendEmailInput{
			Destination: &ses.Destination{
				CcAddresses: []*string{},
				ToAddresses: []*string{
					aws.String(attendee),
				},
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Html: &ses.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(EmailConfirmationHTML(data)),
					},
					Text: &ses.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(TextBody),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(Subject),
				},
			},
			Source: aws.String(Sender),
			// Uncomment to use a configuration set
			//ConfigurationSetName: aws.String(ConfigurationSet),
		}

		// Attempt to send the email.
		result, err := e.client.SendEmail(input)
		if err != nil && !strings.Contains(err.Error(), "Email address is not verified") {
			return nil
		}
		// Display error messages if they occur.
		//if err != nil {
		//	//if aerr, ok := err.(awserr.Error); ok {
		//	//	switch aerr.Code() {
		//	//	case ses.ErrCodeMessageRejected:
		//	//		fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
		//	//	case ses.ErrCodeMailFromDomainNotVerifiedException:
		//	//		fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
		//	//	case ses.ErrCodeConfigurationSetDoesNotExistException:
		//	//		fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
		//	//	default:
		//	//		fmt.Println(aerr.Error())
		//	//	}
		//	//} else {
		//	//	// Print the error, cast err to awserr.Error to get the Code and
		//	//	// Message from an error.
		//	//	fmt.Println(err.Error())
		//	//}
		//	return err
		//}
		if result != nil && result.MessageId != nil {
			fmt.Println("Email Sent to address: " + attendee)
			fmt.Println(result)
		} else {
			fmt.Println("*******Email couldnt be sent to address: " + attendee + " probs need to verify email on aws SES*******")
		}

	}

	return nil
}
