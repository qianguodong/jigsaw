package mongodb

import (
	"context"
	"time"

	"go.uber.org/fx"

	"github.com/sirupsen/logrus"

	"github.com/guodongq/jigsaw/pkg/provider/app"

	"github.com/guodongq/jigsaw/pkg/provider"
	"github.com/guodongq/jigsaw/pkg/provider/probes"
	"github.com/guodongq/jigsaw/pkg/provider/settings"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Module = func() fx.Option {
	return fx.Options(
		fx.Provide(New),
		fx.Invoke(func(lc fx.Lifecycle, p *MongoDB) error {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					return p.Close()
				},
			})
			return p.Init()
		}),
	)
}

type MongoDB struct {
	provider.DefaultProvider
	Settings       *settings.Settings
	appProvider    *app.App
	probesProvider *probes.Probes
	Client         *mongo.Client
	Database       *mongo.Database
}

func New(settings *settings.Settings, appProvider *app.App, probesProvider *probes.Probes) *MongoDB {
	return &MongoDB{
		Settings:       settings,
		appProvider:    appProvider,
		probesProvider: probesProvider,
	}
}

var mongodbCfg struct {
	URI               string `yaml:"uri"`
	Timeout           int    `yaml:"timeout"`
	MaxPoolSize       int    `yaml:"maxPoolSize"`
	MaxConnIdleTime   int    `yaml:"maxConnIdleTime"`
	HeartBeatInterval int    `yaml:"heartBeatInterval"`
	Database          string `yaml:"database"`
}

const (
	defaultMongoURI          = "mongodb://root:admin@127.0.0.1:27017"
	defaultTimeout           = 20
	defaultMaxPoolSize       = 16
	defaultMaxConnIdleTime   = 30
	defaultHeartBeatInterval = 15
	defaultMongoDatabase     = "test"
)

func (p *MongoDB) Init() error {
	mongodbCfg.URI = defaultMongoURI
	mongodbCfg.Timeout = defaultTimeout
	mongodbCfg.MaxPoolSize = defaultMaxPoolSize
	mongodbCfg.MaxConnIdleTime = defaultMaxConnIdleTime
	mongodbCfg.HeartBeatInterval = defaultHeartBeatInterval
	mongodbCfg.Database = defaultMongoDatabase

	if p.Settings.Enable() {
		if err := p.Settings.Get("mongodb").Populate(&mongodbCfg); err != nil {
			return err
		}
	}

	opts := options.
		Client().
		ApplyURI(mongodbCfg.URI).
		SetConnectTimeout(time.Duration(mongodbCfg.Timeout) * time.Second).
		SetMaxPoolSize(uint64(mongodbCfg.MaxPoolSize)).
		SetMaxConnIdleTime(time.Duration(mongodbCfg.MaxConnIdleTime) * time.Second).
		SetHeartbeatInterval(time.Duration(mongodbCfg.HeartBeatInterval) * time.Second)
	if p.appProvider != nil {
		opts.SetAppName(p.appProvider.Name())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongodbCfg.Timeout)*time.Second)
	defer cancel()

	logEntry := logrus.WithField("address", mongodbCfg.URI).WithField("time_out", mongodbCfg.Timeout)
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
	p.Database = client.Database(mongodbCfg.Database)

	// Add live probes if possible.
	if p.probesProvider != nil {
		p.probesProvider.AddLivenessProbes(p.livenessProbe)
	}
	return nil
}

func (p *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongodbCfg.Timeout)*time.Second)
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
