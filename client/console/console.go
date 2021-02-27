package console

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	pub "github.com/WinnersonKharsunai/GraduationProject/client/cmd/services/publisher"
	sub "github.com/WinnersonKharsunai/GraduationProject/client/cmd/services/subscriber"
	"github.com/WinnersonKharsunai/GraduationProject/client/pkg/client"
)

// Console ...
type Console struct {
	clientID      int
	client        *client.Client
	publisherSvc  pub.PublisherService
	subscriberSvc sub.SubscriberService
	shutdwonChan  chan struct{}
	processWg     sync.WaitGroup
}

// NewConsole ...
func NewConsole(clientID int, client *client.Client, pSvc pub.PublisherService, sSvc sub.SubscriberService) *Console {
	return &Console{
		client:        client,
		clientID:      clientID,
		publisherSvc:  pSvc,
		subscriberSvc: sSvc,
		shutdwonChan:  make(chan struct{}),
	}
}

// Start ...
func (c *Console) Start(ctx context.Context) error {
	fmt.Println(welcome)

	switch c.getRole() {
	case publisher:
		c.processWg.Add(1)
		go c.publisherHandler(ctx)
	case subscriber:
		c.processWg.Add(1)
		go c.subscriberHandler(ctx)
	default:
		return errors.New("client not recognised")
	}
	return nil
}

func (c *Console) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		close(c.shutdwonChan)
		c.processWg.Wait()

		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Console) publisherHandler(ctx context.Context) {
	fmt.Println(welcomepublisher)

	stop := false
	for !stop {
		select {
		case <-c.shutdwonChan:
			stop = true
		default:
			shutdown := false
			for !shutdown {
				var input string
				for _, choice := range publisherWelcomeMenu() {
					fmt.Println(choice)
				}

				fmt.Printf("Enter your choice: ")
				fmt.Scanln(&input)

				switch choice(input) {
				case showAllTopics:
					resp, err := c.publisherSvc.ShowTopics(ctx, &pub.ShowTopicRequest{PublisherID: c.clientID})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Print("\nTopics:\t")
					for i, topic := range resp.Topics {
						fmt.Printf("\t%d.%v", i+1, topic)
					}
				case register:
					showTopicsResp, err := c.publisherSvc.ShowTopics(ctx, &pub.ShowTopicRequest{PublisherID: c.clientID})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Print("\nTopics:\t")
					for i, topic := range showTopicsResp.Topics {
						fmt.Printf("\t%d.%v", i+1, topic)
					}

					input := getIntegerInput("\nchoose topic: ")

					if input == 0 || input > len(showTopicsResp.Topics) {
						fmt.Println("Invalid Choice!!!")
						continue
					}

					topicName := showTopicsResp.Topics[input-1]
					resp, err := c.publisherSvc.ConnectToTopic(ctx, &pub.ConnectToTopicRequest{PublisherID: c.clientID, TopicName: topicName})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Printf(resp.Status)

				case deregister:
					resp, err := c.publisherSvc.DisconnectFromTopic(ctx, &pub.DisconnectFromTopicRequest{PublisherID: c.clientID})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Printf("Status: %v", resp.Status)

				case publishMessage:

					resp, err := c.publisherSvc.PublishMessage(ctx, &pub.PublishMessageRequest{PublisherID: c.clientID, Message: getMessage()})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Printf("Status: %v", resp.Status)

				case exitPublisher:
					shutdown = true
					stop = true
				default:
					fmt.Println("Invalid Choice!!!")
				}
				fmt.Print("\n\n")
			}
		}
	}
	c.processWg.Done()
}

func (c *Console) subscriberHandler(ctx context.Context) {
	fmt.Println(welcomeSubscriber)

	stop := false
	for !stop {
		select {
		case <-c.shutdwonChan:
			stop = true
		default:
			shutdown := false
			for !shutdown {
				var input string
				for _, choice := range subscriberWelcomeMenu() {
					fmt.Println(choice)
				}

				fmt.Printf("Enter your choice: ")
				fmt.Scanln(&input)

				switch choice(input) {
				case showAllTopics:
					resp, err := c.subscriberSvc.ShowTopics(ctx, &sub.ShowTopicRequest{SubscriberID: c.clientID})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Print("\nTopics:\t")
					for i, topic := range resp.Topics {
						fmt.Printf("\t%d.%v", i+1, topic)
					}

				case subscribe:
					showTopicsResp, err := c.subscriberSvc.ShowTopics(ctx, &sub.ShowTopicRequest{SubscriberID: c.clientID})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Print("\nTopics:\t")
					for i, topic := range showTopicsResp.Topics {
						fmt.Printf("\t%d.%v", i+1, topic)
					}

					input := getIntegerInput("\nchoose topic: ")

					if input == 0 || input > len(showTopicsResp.Topics) {
						fmt.Println("Invalid Choice!!!")
						continue
					}

					topicName := showTopicsResp.Topics[input-1]

					resp, err := c.subscriberSvc.SubscribeToTopic(ctx, &sub.SubscribeToTopicRequest{SubscriberID: c.clientID, TopicName: topicName})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Printf("Status: %v", resp.Status)

				case showSubscribedTopics:

					resp, err := c.subscriberSvc.GetSubscribedTopics(ctx, &sub.GetSubscribedTopicsRequest{SubscriberID: c.clientID})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Print("\nTopics:\t")
					for i, topic := range resp.Topics {
						fmt.Printf("\t%d.%v", i+1, topic)
					}

				case unsubscribe:
					showTopicsResp, err := c.subscriberSvc.GetSubscribedTopics(ctx, &sub.GetSubscribedTopicsRequest{SubscriberID: c.clientID})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Print("\nTopics:\t")
					for i, topic := range showTopicsResp.Topics {
						fmt.Printf("\t%d.%v", i+1, topic)
					}

					input := getIntegerInput("\nchoose topic")

					if input == 0 || input > len(showTopicsResp.Topics) {
						fmt.Println("Invalid Choice!!!")
						continue
					}

					topicName := showTopicsResp.Topics[input-1]

					resp, err := c.subscriberSvc.UnsubscribeFromTopic(ctx, &sub.UnsubscribeFromTopicRequest{SubscriberID: c.clientID, TopicName: topicName})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Printf("Status: %v", resp.Status)

				case readMessage:
					showTopicsResp, err := c.subscriberSvc.GetSubscribedTopics(ctx, &sub.GetSubscribedTopicsRequest{SubscriberID: c.clientID})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Print("\nTopics:\t")
					for i, topic := range showTopicsResp.Topics {
						fmt.Printf("\t%d.%v", i+1, topic)
					}

					input := getIntegerInput("\nchoose topic")

					if input == 0 || input > len(showTopicsResp.Topics) {
						fmt.Println("Invalid Choice!!!")
						continue
					}

					topicName := showTopicsResp.Topics[input-1]

					resp, err := c.subscriberSvc.GetMessageFromTopic(ctx, &sub.GetMessageFromTopicRequest{SubscriberID: c.clientID, TopicName: topicName})
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					fmt.Printf("Message %+v\n", resp.Message)

				case exitSubscriber:
					shutdown = true
					stop = true

				default:
					fmt.Println("Invalid Choice!!!")
				}
				fmt.Print("\n\n")
			}
		}
	}
	c.processWg.Done()
}

func getMessage() pub.Message {
	now := time.Now().UTC()
	msg := getStringInput("Enter mesage here")
	return pub.Message{
		Data:      msg,
		CretedAt:  now.Format("2006-01-02 15:04:05"),
		ExpiresAt: now.Add(time.Duration(time.Second * 60)).Format("2006-01-02 15:04:05"),
	}
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

func scanInput(msg string) (string, error) {
	var value string

	_, err := fmt.Scanf("%v: %v", msg, &value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getStringInput(msg string) string {
	var value string

	fmt.Printf("%s: ", msg)
	fmt.Scanln(&value)

	return value
}

func getIntegerInput(msg string) int {
	var value int

	fmt.Printf("%s: ", msg)
	fmt.Scanln(&value)

	return value
}

func (c *Console) getRole() clientRole {
	if c.clientID >= 5000 && c.clientID < 6000 {
		return publisher
	} else if c.clientID >= 6000 && c.clientID < 7000 {
		return subscriber
	}
	return unknown
}
