package app

import (
	"context"

	"firebase.google.com/go"
	"google.golang.org/api/option"

	w "github.com/core-go/firestore/writer"
	"github.com/core-go/health"
	fh "github.com/core-go/health/firestore"
	"github.com/core-go/mq"
	v "github.com/core-go/mq/validator"
	"github.com/core-go/mq/zap"
	"github.com/core-go/pubsub"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	Subscribe     func(ctx context.Context, handle func(context.Context, []byte, map[string]string))
	Handle        func(context.Context, []byte, map[string]string)
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	log.Initialize(cfg.Log)
	opts := option.WithCredentialsJSON([]byte(cfg.Firestore.Credentials))
	app, er0 := firebase.NewApp(ctx, nil, opts)
	if er0 != nil {
		return nil, er0
	}

	client, er1 := app.Firestore(ctx)
	if er1 != nil {
		return nil, er1
	}

	logError := log.ErrorMsg
	var logInfo func(context.Context, string)
	if log.IsInfoEnable() {
		logInfo = log.InfoMsg
	}

	subscriber, er2 := pubsub.NewSubscriberByConfig(ctx, cfg.Sub, logError, true)
	if er2 != nil {
		log.Error(ctx, "Cannot create a new subscriber. Error: "+er2.Error())
		return nil, er2
	}

	validator, er3 := v.NewValidator[*User]()
	if er3 != nil {
		return nil, er3
	}
	errorHandler := mq.NewErrorHandler[*User](logError)
	publisher, er4 := pubsub.NewPublisherByConfig(ctx, *cfg.Pub)
	if er4 != nil {
		log.Error(ctx, "Cannot new a new publisher. Error: "+er3.Error())
		return nil, er4
	}
	writer := w.NewWriter[*User](client, "user")
	handler := mq.NewRetryHandlerByConfig[User](cfg.Retry, writer.Write, validator.Validate, errorHandler.RejectWithMap, errorHandler.HandleErrorWithMap, publisher.Publish, logError, logInfo)
	firestoreChecker, er5 := fh.NewHealthChecker(ctx, []byte(cfg.Firestore.Credentials), cfg.Firestore.ProjectId)
	if er5 != nil {
		return nil, er5
	}
	subscriberChecker := pubsub.NewSubHealthChecker("pubsub_subscriber", subscriber.Client, cfg.Sub.SubscriptionId)
	publisherChecker := pubsub.NewPubHealthChecker("pubsub_publisher", publisher.Client, cfg.Pub.TopicId)
	healthHandler := health.NewHandler(firestoreChecker, subscriberChecker, publisherChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		Subscribe:     subscriber.Subscribe,
		Handle:        handler.Handle,
	}, nil
}
