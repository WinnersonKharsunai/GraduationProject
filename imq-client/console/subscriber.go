package console

import (
	"context"
	"errors"
	"fmt"

	"github.com/WinnersonKharsunai/GraduationProject/client/cmd/services/subscriber"
)

func (c *Console) subscriberHandler(ctx context.Context) {
	fmt.Println(welcomeSubscriber)

	shutdown := false
	for !shutdown {
		displayActions(subscriberWelcomeMenu())

		input := getStringInput("Enter your choice")

		switch choice(input) {
		case showAllTopics:
			response, err := processShowSubscriberTopics(ctx, c.subscriberSvc, c.clientID)
			if err != nil {
				displayError(err)
				continue
			}
			displayTopics(response.Topics)

		case subscribe:
			response, err := procesSubscribeToTopic(ctx, c.subscriberSvc, c.clientID)
			if err != nil {
				displayError(err)
				continue
			}
			displayStatus(response.Status)

		case showSubscribedTopics:
			response, err := processShowSubscribedTopics(ctx, c.subscriberSvc, c.clientID)
			if err != nil {
				displayError(err)
				continue
			}
			displayTopics(response.Topics)

		case unsubscribe:
			response, err := processUnsubscribeFromTopic(ctx, c.subscriberSvc, c.clientID)
			if err != nil {
				displayError(err)
				continue
			}
			displayStatus(response.Status)

		case readMessage:
			response, err := processReadMessage(ctx, c.subscriberSvc, c.clientID)
			if err != nil {
				displayError(err)
				continue
			}
			displayMessage(response.Message)

		case exitSubscriber:
			shutdown = true
			c.ShutdwonChan <- struct{}{}

		default:
			displayError(errors.New(invalidChoice))
		}
	}
}

func processShowSubscriberTopics(ctx context.Context, svc subscriber.Service, id int) (*subscriber.ShowTopicResponse, error) {
	showTopicResponse, err := svc.ShowTopics(ctx, &subscriber.ShowTopicRequest{SubscriberID: id})
	if err != nil {
		return nil, err
	}
	return showTopicResponse, nil
}

func procesSubscribeToTopic(ctx context.Context, svc subscriber.Service, id int) (*subscriber.SubscribeToTopicResponse, error) {
	showTopicResponse, err := processShowSubscriberTopics(ctx, svc, id)
	if err != nil {
		return nil, err
	}

	displayTopics(showTopicResponse.Topics)

	input := getIntegerInput("Choose topic")
	if input == 0 || input > len(showTopicResponse.Topics) {
		return nil, errors.New(invalidChoice)
	}

	topicName := showTopicResponse.Topics[input-1]

	subscribeToTopicResponse, err := svc.SubscribeToTopic(ctx, &subscriber.SubscribeToTopicRequest{SubscriberID: id, TopicName: topicName})
	if err != nil {
		return nil, err
	}
	return subscribeToTopicResponse, nil
}

func processShowSubscribedTopics(ctx context.Context, svc subscriber.Service, id int) (*subscriber.GetSubscribedTopicsResponse, error) {
	getSubscribedTopicsResponse, err := svc.GetSubscribedTopics(ctx, &subscriber.GetSubscribedTopicsRequest{SubscriberID: id})
	if err != nil {
		return nil, err
	}
	return getSubscribedTopicsResponse, nil
}

func processUnsubscribeFromTopic(ctx context.Context, svc subscriber.Service, id int) (*subscriber.UnsubscribeFromTopicResponse, error) {
	getSubscribedTopicsResponse, err := processShowSubscribedTopics(ctx, svc, id)
	if err != nil {
		return nil, err
	}
	displayTopics(getSubscribedTopicsResponse.Topics)

	input := getIntegerInput("Choose topic")
	if input == 0 || input > len(getSubscribedTopicsResponse.Topics) {
		return nil, errors.New(invalidChoice)
	}

	topicName := getSubscribedTopicsResponse.Topics[input-1]
	unsubscribeFromTopicResponse, err := svc.UnsubscribeFromTopic(ctx, &subscriber.UnsubscribeFromTopicRequest{SubscriberID: id, TopicName: topicName})
	if err != nil {
		return nil, err
	}
	return unsubscribeFromTopicResponse, nil
}

func processReadMessage(ctx context.Context, svc subscriber.Service, id int) (*subscriber.GetMessageFromTopicResponse, error) {
	getSubscribedTopicsResponse, err := processShowSubscribedTopics(ctx, svc, id)
	if err != nil {
		return nil, err
	}
	displayTopics(getSubscribedTopicsResponse.Topics)

	input := getIntegerInput("Choose topic")
	if input == 0 || input > len(getSubscribedTopicsResponse.Topics) {
		return nil, errors.New(invalidChoice)
	}

	topicName := getSubscribedTopicsResponse.Topics[input-1]
	getMessageFromTopicResponse, err := svc.GetMessageFromTopic(ctx, &subscriber.GetMessageFromTopicRequest{SubscriberID: id, TopicName: topicName})
	if err != nil {
		return nil, err
	}
	return getMessageFromTopicResponse, nil
}

func subscriberWelcomeMenu() []string {
	return []string{
		"1. Show all Topics",
		"2. Subscribe to Topic",
		"3. Show all Subscribed Topics",
		"4. Unsubscribe from Topic",
		"5. Read Message from Topic",
		"6. Exit",
	}
}
