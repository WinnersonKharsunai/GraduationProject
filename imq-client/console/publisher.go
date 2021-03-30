package console

import (
	"context"
	"errors"
	"fmt"

	"github.com/WinnersonKharsunai/GraduationProject/client/cmd/services/publisher"
)

func (c *Console) publisherHandler(ctx context.Context) {
	fmt.Println(welcomepublisher)

	shutdown := false
	for !shutdown {
		displayActions(publisherWelcomeMenu())

		input := getStringInput("Enter your choice")

		switch choice(input) {
		case showAllTopics:
			response, err := processShowAllTopics(ctx, c.publisherSvc, c.clientID)
			if err != nil {
				displayError(err)
				continue
			}
			displayTopics(response.Topics)

		case register:
			response, err := processRegisterToTopic(ctx, c.publisherSvc, c.clientID)
			if err != nil {
				displayError(err)
				continue
			}
			displayStatus(response.Status)

		case deregister:
			response, err := processDeregisterFromTopic(ctx, c.publisherSvc, c.clientID)
			if err != nil {
				displayError(err)
				continue
			}
			displayStatus(response.Status)

		case publishMessage:
			response, err := processPublishMessage(ctx, c.publisherSvc, c.clientID)
			if err != nil {
				displayError(err)
				continue
			}
			displayStatus(response.Status)

		case exitPublisher:
			shutdown = true
			c.ShutdwonChan <- struct{}{}

		default:
			displayError(errors.New(invalidChoice))
		}
	}
}

func processShowAllTopics(ctx context.Context, svc publisher.Service, id int) (*publisher.ShowTopicResponse, error) {
	showTopicResponse, err := svc.ShowTopics(ctx, &publisher.ShowTopicRequest{PublisherID: id})
	if err != nil {
		return nil, err
	}
	return showTopicResponse, nil
}

func processRegisterToTopic(ctx context.Context, svc publisher.Service, id int) (*publisher.ConnectToTopicResponse, error) {
	showTopicResponse, err := processShowAllTopics(ctx, svc, id)
	if err != nil {
		return nil, err
	}

	displayTopics(showTopicResponse.Topics)

	input := getIntegerInput("Choose topic")
	if input == 0 || input > len(showTopicResponse.Topics) {
		return nil, errors.New(invalidChoice)
	}

	topicName := showTopicResponse.Topics[input-1]

	connectToTopicResponse, err := svc.ConnectToTopic(ctx, &publisher.ConnectToTopicRequest{PublisherID: id, TopicName: topicName})
	if err != nil {
		return nil, err
	}
	return connectToTopicResponse, nil
}

func processDeregisterFromTopic(ctx context.Context, svc publisher.Service, id int) (*publisher.DisconnectFromTopicResponse, error) {
	disconnectFromTopicResponse, err := svc.DisconnectFromTopic(ctx, &publisher.DisconnectFromTopicRequest{PublisherID: id})
	if err != nil {
		return nil, err
	}
	return disconnectFromTopicResponse, nil
}

func processPublishMessage(ctx context.Context, svc publisher.Service, id int) (*publisher.PublishMessageResponse, error) {
<<<<<<< HEAD
	msg, err := getMessage()
	if err != nil {
		return nil, err
	}

	publishMessageResponse, err := svc.PublishMessage(ctx, &publisher.PublishMessageRequest{PublisherID: id, Message: msg})
=======
	publishMessageResponse, err := svc.PublishMessage(ctx, &publisher.PublishMessageRequest{PublisherID: id, Message: getMessage()})
>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469
	if err != nil {
		return nil, err
	}
	return publishMessageResponse, nil
}

func publisherWelcomeMenu() []string {
	return []string{
		"1. Show all Topics",
		"2. Register to Topic",
		"3. Deregister from Topic",
		"4. Publish message to Topic",
		"5. Exit",
	}
}
