package provider

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

type Provider interface {
	Init() error
	Close() error
	Run() error
	IsRunning() bool
}

type DefaultProvider struct {
	Provider
	running bool
}

func (p *DefaultProvider) Init() error {
	return nil
}

func (p *DefaultProvider) Close() error {
	p.SetRunning(false)
	return nil
}

func (p *DefaultProvider) SetRunning(running bool) {
	p.running = running
}

func (p *DefaultProvider) IsRunning() bool {
	return p.running
}

func WaitForProvider(p Provider, timeoutSeconds time.Duration) error {
	if p.IsRunning() {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()

	name := Name(p)
	logrus.WithField("timeout", timeoutSeconds).Debugf("Waiting for %s to run...", name)
	for {
		if p.IsRunning() {
			return nil
		}

		if ctx.Err() != nil {
			return fmt.Errorf("time exceeded for %s to run", name)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func Name(provider Provider) string {
	return reflect.ValueOf(provider).Elem().String()
}
