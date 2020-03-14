package services

import (
	"github.com/statping/statping/types/failures"
	"github.com/statping/statping/types/notifications"
)

var (
	allNotifiers []ServiceNotifier
)

func AllNotifiers() []ServiceNotifier {
	return allNotifiers
}

type ServiceNotifier interface {
	OnSuccess(*Service) error                    // OnSuccess is triggered when a service is successful
	OnFailure(*Service, *failures.Failure) error // OnFailure is triggered when a service is failing
	OnTest() error                               // OnTest is triggered for testing
	Select() *notifications.Notification         // OnTest is triggered for testing
}
