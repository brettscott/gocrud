package crud

import (
	"fmt"
	"net/http"
)

type Crud struct {
	entities           Entities
	config             *Config
	log                Logger
	statsd             StatsDer
	stores             []Storer
	apiService         apiService
	elementsValidators []elementsValidatorer
	mutators           []mutatorer
}

// NewCrud creates a new CRUD instance
func NewCrud(config *Config, log Logger, statsd StatsDer) *Crud {
	return &Crud{
		config:   config,
		log:      log,
		statsd:   statsd,
		entities: make(map[string]*Entity),
	}
}

// AddStore allows user to define databases to be used
func (c *Crud) AddStore(store Storer) {
	c.stores = append(c.stores, store)
}

// GetStores returns all registered stores
func (c *Crud) GetStores() []Storer {
	return c.stores
}

// AddEntity for each entity type (eg User)
func (c *Crud) AddEntity(entity *Entity) {
	c.entities[entity.ID] = entity
}

func (c *Crud) GetEntities() Entities {
	return c.entities
}

// AddElementsValidator for all entities
func (c *Crud) AddElementsValidator(elementsValidator elementsValidatorer) {
	c.elementsValidators = append(c.elementsValidators, elementsValidator)
}

func (c *Crud) GetElementsValidators() []elementsValidatorer {
	return c.elementsValidators
}

// AddEntityElementsValidator adds an elements validator to a specific entity
func (c *Crud) AddEntityElementsValidator(entityID string, elementsValidator elementsValidatorer) {
	if _, ok := c.entities[entityID]; !ok {
		panic(fmt.Sprintf("Entity %s is not yet registered.  Please register first.", entityID))
	}
	c.entities[entityID].AddElementsValidator(elementsValidator)
}

func (c *Crud) GetEntityElementsValidators(entityID string) ([]elementsValidatorer, error) {
	if _, ok := c.entities[entityID]; !ok {
		return nil, fmt.Errorf("Entity %s is not yet registered.  Please register first.", entityID)
	}
	return c.entities[entityID].ElementsValidators, nil
}

// AddMutator for all entities
func (c *Crud) AddMutator(mutator mutatorer) {
	c.mutators = append(c.mutators, mutator)
}

func (c *Crud) GetMutators() []mutatorer {
	return c.mutators
}

// AddEntityMutator for entity
func (c *Crud) AddEntityMutator(entityID string, mutator mutatorer) {
	if _, ok := c.entities[entityID]; !ok {
		panic(fmt.Sprintf("Entity %s is not yet registered.  Please register first.", entityID))
	}
	c.entities[entityID].AddMutator(mutator)
}

func (c *Crud) GetEntityMutators(entityID string) ([]mutatorer, error) {
	if _, ok := c.entities[entityID]; !ok {
		return nil, fmt.Errorf("Entity %s is not yet registered.  Please register first.", entityID)
	}
	return c.entities[entityID].Mutators, nil
}

// Handler for mounting routes for CRUD
func (c *Crud) Handler() http.Handler {

	if len(c.elementsValidators) == 0 {
		defaultElementsValidator := NewElementsValidator()
		c.AddElementsValidator(defaultElementsValidator)
	}

	c.apiService = newApiService(c.stores, c.elementsValidators, c.mutators)

	healthcheckHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Healthy")
	})

	apiRouteHandler := NewApiRoute(c.entities, &c.apiService, c.log, c.statsd)
	uiRouteHandler := NewUiRoute(c.entities, &c.apiService, c.log, c.statsd)

	return newRouter(c.log, c.statsd, healthcheckHandler, apiRouteHandler, uiRouteHandler)
}
