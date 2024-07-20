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
	Receive       func(ctx context.Context, handle func(context.Context, []byte, map[string]string))
	Handle        func(context.Context, []byte, map[string]string)
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	log.Initialize(cfg.Log)
	opts := option.WithCredentialsJSON([]byte(cfg.Firestore.Credentials))
	app, err := firebase.NewApp(ctx, nil, opts)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	logError := log.ErrorMsg
	var logInfo func(context.Context, string)
	if log.IsInfoEnable() {
		logInfo = log.InfoMsg
	}

	receiver, er2 := pubsub.NewSubscriberByConfig(ctx, cfg.Sub, logError, true)
	if er2 != nil {
		log.Error(ctx, "Cannot create a new receiver. Error: "+er2.Error())
		return nil, er2
	}

	validator, err := v.NewValidator[*User]()
	if err != nil {
		return nil, err
	}
	errorHandler := mq.NewErrorHandler[*User](logError)
	sender, er3 := pubsub.NewPublisherByConfig(ctx, *cfg.Pub)
	if er3 != nil {
		log.Error(ctx, "Cannot new a new sender. Error: "+er3.Error())
		return nil, er3
	}
	if err != nil {
		return nil, err
	}
	writer := w.NewWriter[*User](client, "user")
	handler := mq.NewRetryHandlerByConfig[User](cfg.Retry, writer.Write, validator.Validate, errorHandler.RejectWithMap, errorHandler.HandleErrorWithMap, sender.Publish, logError, logInfo)
	firestoreChecker, er5 := fh.NewHealthChecker(ctx, []byte(cfg.Firestore.Credentials), cfg.Firestore.ProjectId)
	if er5 != nil {
		return nil, er5
	}
	receiverChecker := pubsub.NewSubHealthChecker("pubsub_subscriber", receiver.Client, cfg.Sub.SubscriptionId)
	senderChecker := pubsub.NewPubHealthChecker("pubsub_publisher", sender.Client, cfg.Pub.TopicId)
	healthHandler := health.NewHandler(firestoreChecker, receiverChecker, senderChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		Receive:       receiver.Subscribe,
		Handle:        handler.Handle,
	}, nil
}
