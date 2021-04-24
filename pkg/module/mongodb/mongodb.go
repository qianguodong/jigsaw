package mongodb

import (
	"context"
	"github.com/guodongq/jigsaw/pkg/module"
	"github.com/guodongq/jigsaw/pkg/module/app"
	"github.com/guodongq/jigsaw/pkg/module/probes"
	"time"

	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	module.DefaultProvider
	Config         *Config
	appProvider    *app.App
	probesProvider *probes.Probes
	Client         *mongo.Client
	Database       *mongo.Database
}

func New(config *Config, appProvider *app.App, probesProvider *probes.Probes) *MongoDB {
	return &MongoDB{
		Config:         config,
		appProvider:    appProvider,
		probesProvider: probesProvider,
	}
}

func (p *MongoDB) Init() error {
	opts := options.
		Client().
		ApplyURI(p.Config.URI).
		SetConnectTimeout(time.Duration(p.Config.Timeout) * time.Second).
		SetMaxPoolSize(uint64(p.Config.MaxPoolSize)).
		SetMaxConnIdleTime(time.Duration(p.Config.MaxConnIdleTime) * time.Second).
		SetHeartbeatInterval(time.Duration(p.Config.HeartBeatInterval) * time.Second)
	if p.appProvider != nil {
		opts.SetAppName(p.appProvider.Name())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.Config.Timeout)*time.Second)
	defer cancel()

	logEntry := logrus.WithField("address", p.Config.URI).WithField("time_out", p.Config.Timeout)
	logEntry.Debug("Connecting to MongoDB server...")

	client, err := mongo.NewClient(opts)
	if err != nil {
		logEntry.WithError(err).Error("MongoDB client creation failed")
		return err
	}
	err = client.Connect(ctx)
	if err != nil {
		logEntry.WithError(err).Error("MongoDB connection failed")
		return err
	}

	// Check connection by pinging.
	err = client.Ping(ctx, nil)
	if err != nil {
		logEntry.WithError(err).Error("MongoDB ping failed")
		return err
	}
	p.Client = client
	p.Database = client.Database(p.Config.Database)

	// Add live probes if possible.
	if p.probesProvider != nil {
		p.probesProvider.AddLivenessProbes(p.livenessProbe)
	}
	return nil
}

func (p *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.Config.Timeout)*time.Second)
	defer cancel()

	err := p.Client.Disconnect(ctx)
	if err != nil {
		logrus.WithError(err).Info("MongoDB disconnecting failed")
		return err
	}

	return p.DefaultProvider.Close()
}

func (p *MongoDB) livenessProbe() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.Client.Ping(ctx, nil)
	if err != nil {
		logrus.WithError(err).Error("MongoDB liveness probe failed")
		return err
	}

	logrus.Debug("MongoDB liveness probe succeeded")
	return nil
}
