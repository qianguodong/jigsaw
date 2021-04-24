package module

type Provider interface {
	Init() error
	Close() error
}

type RunProvider interface {
	Provider

	Run() error
	IsRunning() bool
}

type DefaultProvider struct {
	Provider
}

func (*DefaultProvider) Init() error {
	return nil
}

func (*DefaultProvider) Close() error {
	return nil
}

type DefaultRunProvider struct {
	RunProvider
	running bool
}

func (p *DefaultRunProvider) Init() error {
	return nil
}

func (p *DefaultRunProvider) Close() error {
	p.SetRunning(false)
	return nil
}

func (p *DefaultRunProvider) IsRunning() bool {
	return p.running
}

func (p *DefaultRunProvider) SetRunning(running bool) {
	p.running = running
}
