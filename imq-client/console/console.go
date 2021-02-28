package console

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/WinnersonKharsunai/GraduationProject/client/cmd/services/publisher"
	"github.com/WinnersonKharsunai/GraduationProject/client/cmd/services/subscriber"
	"github.com/WinnersonKharsunai/GraduationProject/client/pkg/client"
)

// Console is the reciever type for Console
type Console struct {
	clientID      int
	client        client.Service
	publisherSvc  publisher.Service
	subscriberSvc subscriber.Service
	ShutdwonChan  chan struct{}
}

// NewConsole is the factory function for Console type
func NewConsole(clientID int, client client.Service, pSvc publisher.Service, sSvc subscriber.Service) *Console {
	return &Console{
		client:        client,
		clientID:      clientID,
		publisherSvc:  pSvc,
		subscriberSvc: sSvc,
		ShutdwonChan:  make(chan struct{}),
	}
}

// Start starts running the console
func (c *Console) Start(ctx context.Context) error {
	fmt.Println(welcome)

	switch c.getRole() {
	case publisherRole:
		go c.publisherHandler(ctx)
	case subscriberRole:
		go c.subscriberHandler(ctx)
	default:
		return errors.New("client not recognised")
	}
	return nil
}

func getMessage() publisher.Message {
	msg := getStringInput("Enter mesage here")
	now := time.Now().UTC()

	return publisher.Message{
		Data:      msg,
		CretedAt:  now.Format("2006-01-02 15:04:05"),
		ExpiresAt: now.Add(time.Duration(time.Second * 60)).Format("2006-01-02 15:04:05"),
	}
}

func getStringInput(msg string) string {
	var value string

	fmt.Printf("\n%s: ", msg)
	fmt.Scanln(&value)

	return value
}

func getIntegerInput(msg string) int {
	var value int

	fmt.Printf("\n%s: ", msg)
	fmt.Scanln(&value)

	return value
}

func (c *Console) getRole() clientRole {
	if c.clientID >= 5000 && c.clientID < 6000 {
		return publisherRole
	} else if c.clientID >= 6000 && c.clientID < 7000 {
		return subscriberRole
	}
	return unknown
}

func displayActions(actions []string) {
	fmt.Println()
	for _, option := range actions {
		fmt.Println(option)
	}
}

func displayError(err error) {
	fmt.Printf("\nERROR: %v\n", err)
}

func displayTopics(topics []string) {
	fmt.Print("\nTOPICS:")
	for i, topic := range topics {
		fmt.Printf("\t%d.%v\n", i+1, topic)
	}
}

func displayStatus(status string) {
	fmt.Printf("\nSTATUS: %v\n", status)
}

func displayMessage(msg subscriber.Message) {
	if msg.Data != "" {
		fmt.Printf("\nMESSAGE\n\tData: %v\n\tCreatedAt: %v\n\tExpiredAt: %v\n",
			msg.Data, msg.CretedAt, msg.ExpiresAt)
	} else {
		fmt.Println("MESSAGE:\tNo message recieved")
	}
}
