package app

import (
	"github.com/core-go/health/server"
	"github.com/core-go/mq"
	"github.com/core-go/mq/zap"
	"github.com/core-go/pubsub"
)

type Config struct {
	Server    server.ServerConf       `mapstructure:"server"`
	Log       log.Config              `mapstructure:"log"`
	Firestore FirestoreConfig         `mapstructure:"firestore"`
	Handler   mq.HandlerConfig        `mapstructure:"handler"`
	Sub       pubsub.SubscriberConfig `mapstructure:"sub"`
	Pub       *pubsub.PublisherConfig `mapstructure:"pub"`
}

type FirestoreConfig struct {
	ProjectId   string `yaml:"project_id" mapstructure:"project_id" json:"projectId,omitempty" gorm:"column:projectid" bson:"projectId,omitempty" dynamodbav:"projectId,omitempty" firestore:"projectId,omitempty"`
	Credentials string `yaml:"credentials" mapstructure:"credentials" json:"credentials,omitempty" gorm:"column:credentials" bson:"credentials,omitempty" dynamodbav:"credentials,omitempty" firestore:"credentials,omitempty"`
}
