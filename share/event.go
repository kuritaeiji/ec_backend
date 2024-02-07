package share

import "github.com/cockroachdb/errors"

type (
	DomainEvent interface {
		Name() DomainEventName
	}

	DomainEventName string

	DomainEventPublisher interface {
		Publish(events []DomainEvent) error
		Subscribe(event DomainEvent, subscriber DomainEventSubscriber)
	}

	domainEventPublisher struct {
		subscribers map[DomainEventName][]DomainEventSubscriber
	}

	DomainEventSubscriber interface {
		Subscribe(event DomainEvent) error
	}
)

func NewDomainEventPublisher() DomainEventPublisher {
	return &domainEventPublisher{}
}

func (publisher domainEventPublisher) Publish(events []DomainEvent) error {
	for _, event := range events {
		subscribers, ok := publisher.subscribers[event.Name()]
		if !ok {
			return errors.New(string(event.Name()) + "に対するサブスクライバーが存在しません")
		}

		for _, subscriber := range subscribers {
			err := subscriber.Subscribe(event)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (publisher *domainEventPublisher) Subscribe(event DomainEvent, subscriber DomainEventSubscriber) {
	publisher.subscribers[event.Name()] = append(publisher.subscribers[event.Name()], subscriber)
}
