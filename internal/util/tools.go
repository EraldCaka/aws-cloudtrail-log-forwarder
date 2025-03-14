package tools

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
)

func ConvertEventToMap(event types.Event) map[string]interface{} {
	eventMap := map[string]interface{}{
		"accessKeyId":     event.AccessKeyId,
		"cloudTrailEvent": event.CloudTrailEvent,
		"eventId":         event.EventId,
		"eventName":       event.EventName,
		"eventSource":     event.EventSource,
		"eventTime":       event.EventTime.Format(time.RFC3339),
		"readOnly":        event.ReadOnly,
		"resources":       event.Resources,
		"username":        event.Username,
	}

	return eventMap
}
