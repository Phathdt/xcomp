package xcomp

import (
	"fmt"
	"reflect"
	"sync"
)

type Container struct {
	services map[string]any
	mutex    sync.RWMutex
}

func NewContainer() *Container {
	return &Container{
		services: make(map[string]any),
	}
}

func (c *Container) Register(name string, service any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.services[name] = service
}

func (c *Container) RegisterSingleton(name string, factory func(*Container) any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.services[name] = &lazyService{factory: factory, container: c}
}

type lazyService struct {
	factory   func(*Container) any
	container *Container
	instance  any
	once      sync.Once
}

func (ls *lazyService) getInstance() any {
	ls.once.Do(func() {
		ls.instance = ls.factory(ls.container)
	})
	return ls.instance
}

func (c *Container) Get(name string) any {
	c.mutex.RLock()
	service := c.services[name]
	c.mutex.RUnlock()

	if lazyService, ok := service.(*lazyService); ok {
		return lazyService.getInstance()
	}
	return service
}

func (c *Container) GetTyped(name string, target any) bool {
	service := c.Get(name)
	if service == nil {
		return false
	}

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return false
	}

	serviceValue := reflect.ValueOf(service)
	if !serviceValue.Type().AssignableTo(targetValue.Elem().Type()) {
		return false
	}

	targetValue.Elem().Set(serviceValue)
	return true
}

func (c *Container) Inject(target any) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	targetValue = targetValue.Elem()
	targetType := targetValue.Type()

	for i := 0; i < targetValue.NumField(); i++ {
		field := targetValue.Field(i)
		fieldType := targetType.Field(i)

		injectTag := fieldType.Tag.Get("inject")
		if injectTag == "" {
			continue
		}

		if !field.CanSet() {
			continue
		}

		service := c.Get(injectTag)
		if service == nil {
			return fmt.Errorf("service '%s' not found for field '%s'", injectTag, fieldType.Name)
		}

		serviceValue := reflect.ValueOf(service)
		if !serviceValue.Type().AssignableTo(field.Type()) {
			return fmt.Errorf("service '%s' is not assignable to field '%s'", injectTag, fieldType.Name)
		}

		field.Set(serviceValue)
	}

	return nil
}

func (c *Container) AutoWire(target any) error {
	return c.Inject(target)
}

func (c *Container) ListServices() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	services := make([]string, 0, len(c.services))
	for name := range c.services {
		services = append(services, name)
	}
	return services
}
