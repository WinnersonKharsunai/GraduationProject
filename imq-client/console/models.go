package console

type clientRole int
type choice string

const (
	unknown clientRole = iota
	publisher
	subscriber

	showAllTopics        choice = "1"
	register             choice = "2"
	deregister           choice = "3"
	publishMessage       choice = "4"
	exitPublisher        choice = "5"
	subscribe            choice = "2"
	showSubscribedTopics choice = "3"
	unsubscribe          choice = "4"
	readMessage          choice = "5"
	exitSubscriber       choice = "6"

	welcome           = "Welcome to ITT Messaging Queue"
	welcomepublisher  = "You are logged in as publisher"
	welcomeSubscriber = "You are logged in as subscriber"
)
