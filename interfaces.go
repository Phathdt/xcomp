package xcomp

type Injectable interface {
	GetServiceName() string
}

type Service interface {
	Injectable
	Initialize(*Container) error
}

type Module interface {
	GetProviders() []Provider
	GetImports() []Module
}

type Provider struct {
	Name    string
	Factory func(*Container) any
	Service any
}

func NewProvider(name string, factory func(*Container) any) Provider {
	return Provider{
		Name:    name,
		Factory: factory,
	}
}

func NewServiceProvider(name string, service any) Provider {
	return Provider{
		Name:    name,
		Service: service,
	}
}

type ModuleBuilder struct {
	providers []Provider
	imports   []Module
}

func NewModule() *ModuleBuilder {
	return &ModuleBuilder{
		providers: make([]Provider, 0),
		imports:   make([]Module, 0),
	}
}

func (mb *ModuleBuilder) AddProvider(provider Provider) *ModuleBuilder {
	mb.providers = append(mb.providers, provider)
	return mb
}

func (mb *ModuleBuilder) AddService(name string, service any) *ModuleBuilder {
	mb.providers = append(mb.providers, NewServiceProvider(name, service))
	return mb
}

func (mb *ModuleBuilder) AddFactory(name string, factory func(*Container) any) *ModuleBuilder {
	mb.providers = append(mb.providers, NewProvider(name, factory))
	return mb
}

func (mb *ModuleBuilder) Import(module Module) *ModuleBuilder {
	mb.imports = append(mb.imports, module)
	return mb
}

func (mb *ModuleBuilder) Build() Module {
	return &BasicModule{
		providers: mb.providers,
		imports:   mb.imports,
	}
}

type BasicModule struct {
	providers []Provider
	imports   []Module
}

func (bm *BasicModule) GetProviders() []Provider {
	return bm.providers
}

func (bm *BasicModule) GetImports() []Module {
	return bm.imports
}

func (c *Container) RegisterModule(module Module) error {
	for _, importedModule := range module.GetImports() {
		if err := c.RegisterModule(importedModule); err != nil {
			return err
		}
	}

	for _, provider := range module.GetProviders() {
		if provider.Factory != nil {
			c.RegisterSingleton(provider.Name, provider.Factory)
		} else if provider.Service != nil {
			c.Register(provider.Name, provider.Service)
		}
	}

	return nil
}
