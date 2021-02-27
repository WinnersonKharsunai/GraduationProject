package publisher

const (
	version              = "1.0"
	contentType          = "json"
	showTopic            = "showTopicRequest"
	connectToTopic       = "connectToTopicRequest"
	disconnectFromTopic  = "disconnectFromTopicRequest"
	publishMessage       = "publishMessageRequest"
	checkMessageStatus   = "checkMessageStatusRequest"
	subscribeToTopic     = "subscribeToTopicRequest"
	unsubscribeFromTopic = "unsubscribeFromTopicRequest"
	getSubscribedTopics  = "getSubscribedTopicsRequest"
	getMessageFromTopic  = "getMessageFromTopicRequest"
)

// ShowTopicRequest holds the request details for ShowTopics
type ShowTopicRequest struct {
	PublisherID int `json:"publisherId" xml:"publisherId"`
}

// ShowTopicResponse holds the response details for ShowTopics
type ShowTopicResponse struct {
	Topics []string `json:"topics" xml:"topics"`
}

// ConnectToTopicRequest holds the request details for ConnectToTopic
type ConnectToTopicRequest struct {
	PublisherID int    `json:"publisherId" xml:"publisherId"`
	TopicName   string `json:"topicName" xml:"topicName"`
}

// ConnectToTopicResponse holds the response details for ConnectToTopic
type ConnectToTopicResponse struct {
	Status string `json:"status" xml:"status"`
}

// DisconnectFromTopicRequest holds the request details for DisconnectFromTopic
type DisconnectFromTopicRequest struct {
	PublisherID int `json:"publisherId" xml:"publisherId"`
}

// DisconnectFromTopicResponse holds the response details for DisconnectFromTopic
type DisconnectFromTopicResponse struct {
	Status string `json:"status" xml:"status"`
}

// PublishMessageRequest holds the request details for PublishMessage
type PublishMessageRequest struct {
	PublisherID int     `json:"publisherId" xml:"publisherId"`
	Message     Message `json:"message" xml:"message"`
}

// Message holds the message details
type Message struct {
	Data      string `json:"data" xml:"data"`
	CretedAt  string `json:"cretedAt" xml:"cretedAt"`
	ExpiresAt string `json:"expiredAt" xml:"expiredAt"`
}

// PublishMessageResponse holds the response details for PublishMessage
type PublishMessageResponse struct {
	Status string `json:"status" xml:"status"`
}

// CheckMessageStatusRequest holds the request details for  CheckMessageStatus
type CheckMessageStatusRequest struct {
	PublisherID int     `json:"publisherId" xml:"publisherId"`
	Message     Message `json:"message" xml:"message"`
}

// CheckMessageStatusResponse holds the response details for CheckMessageStatus
type CheckMessageStatusResponse struct {
	Status string `json:"status" xml:"status"`
}
