package storer

import (
	"github.com/OpenCal-FYDP/AsyncCalendarOptimizer/internal/set"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"time"
)

type EventData struct {
	CalendarEventID string
	//TimeToNotify string
	Start     int64
	End       int64
	Attendees []string
	Location  string
	Summary   string
}

type Storage struct {
	client dynamodbiface.DynamoDBAPI
	set    *set.Set
}

func New() *Storage {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := dynamodb.New(sess)

	return &Storage{
		client: client,
		set:    set.NewSet(),
	}
}

func (s *Storage) GetEvents() ([]*EventData, error) {
	filt := expression.Name("Attendees").NotEqual(expression.Value(nil))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, err
	}

	// Now create put item
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String("calEvents"),
	}

	// we can get away with a non-paginated scan because we have so few events
	res, err := s.client.Scan(input)
	if err != nil {
		return nil, err
	}

	retEvents := []*EventData{}

	count := 0
	for _, item := range res.Items {
		retEvent := &EventData{}
		err = dynamodbattribute.UnmarshalMap(item, retEvent)
		if err != nil {
			return nil, err
		}

		// filter items
		startTime := time.Unix(retEvent.Start, 0)

		// check if time is in less than 1hour
		if time.Until(startTime) > 0 && time.Until(startTime) < time.Hour {
			if !s.set.Contains(retEvent.CalendarEventID) {
				s.set.Add(retEvent.CalendarEventID)
				retEvents = append(retEvents, retEvent)
			}
			count += 1
		}
	}

	return retEvents, nil
}
