package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type EventMessage struct {
	EventID   string `json:"event_id"`
	EventName string `json:"event_name"`
	Action    string `json:"action"`
}

type SQSClient struct {
	Client   *sqs.Client
	QueueURL string
}

func (s *SQSClient) SendMessage(message string) error {
	ctx := context.Background()
	_, err := s.Client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.QueueURL),
		MessageBody: aws.String(message),
	})
	if err != nil {
		return fmt.Errorf("error sending SQS message: %w", err)
	}
	return nil
}

func (s *SQSClient) SendEventMessage(ctx context.Context, msg EventMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling SQS message: %w", err)
	}

	_, err = s.Client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.QueueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("error sending SQS message: %w", err)
	}
	return nil
}

func (s *SQSClient) ReceiveEventMessages(ctx context.Context, maxMessages int32) ([]EventMessage, error) {
	resp, err := s.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.QueueURL),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     10,
	})
	if err != nil {
		return nil, fmt.Errorf("error receiving SQS messages: %w", err)
	}

	var messages []EventMessage
	for _, m := range resp.Messages {
		var msg EventMessage
		if err := json.Unmarshal([]byte(*m.Body), &msg); err == nil {
			messages = append(messages, msg)
		}
	}
	return messages, nil
} 