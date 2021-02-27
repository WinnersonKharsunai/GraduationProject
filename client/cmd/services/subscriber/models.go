package subscriber

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
	SubscriberID int `json:"publisherId" xml:"publisherId"`
}

// ShowTopicResponse holds the response details for ShowTopics
type ShowTopicResponse struct {
	Topics []string `json:"topics" xml:"topics"`
}

// SubscribeToTopicRequest holds the request details for SubscribeToTopic
type SubscribeToTopicRequest struct {
	SubscriberID int    `json:"subscriberId" xml:"subscriberId"`
	TopicName    string `json:"topicName" xml:"topicName"`
}

// SubscribeToTopicResponse holds the response details for SubscribeToTopic
type SubscribeToTopicResponse struct {
	Status string `json:"status" xml:"status"`
}

// UnsubscribeFromTopicRequest holds the request details for UnsubscribeFromTopic
type UnsubscribeFromTopicRequest struct {
	SubscriberID int    `json:"subscriberId" xml:"subscriberId"`
	TopicName    string `json:"topicName" xml:"topicName"`
}

// UnsubscribeFromTopicResponse holds the response details for UnsubscribeFromTopic
type UnsubscribeFromTopicResponse struct {
	Status string `json:"status" xml:"status"`
}

// GetSubscribedTopicsRequest holds the request details for GetSubscribedTopics
type GetSubscribedTopicsRequest struct {
	SubscriberID int `json:"subscriberId" xml:"subscriberId"`
}

// GetSubscribedTopicsResponse holds the response details for GetSubscribedTopics
type GetSubscribedTopicsResponse struct {
	Topics []string `json:"topics" xml:"topics"`
}

// GetMessageFromTopicRequest holds the request details for GetMessageFromTopic
type GetMessageFromTopicRequest struct {
	SubscriberID int    `json:"subscriberId" xml:"subscriberId"`
	TopicName    string `json:"topicName" xml:"topicName"`
}

// GetMessageFromTopicResponse holds the response details for GetMessageFromTopic
type GetMessageFromTopicResponse struct {
	Message Message `json:"message" xml:"message"`
}

// Message holds the message details
type Message struct {
	Data      string `json:"data" xml:"data"`
	CretedAt  string `json:"cretedAt" xml:"cretedAt"`
	ExpiresAt string `json:"expiredAt" xml:"expiredAt"`
}
