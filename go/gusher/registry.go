package gusher

import (
	"log"
)

type registry struct {
	gsessions map[string]*gsession
	add       chan *gsession
	remove    chan *gsession
}

func newRegistry() *registry {
	registry := &registry{
		gsessions: make(map[string]*gsession),
		add:       make(chan *gsession),
		remove:    make(chan *gsession),
	}
	go registry.run()
	return registry
}

func (r *registry) run() {
	defer r.teardown()
	for {
		select {
		case gsession := <-r.add:
			r.register(gsession.ID(), gsession)
		case gsession := <-r.remove:
			r.unregister(gsession.ID())
		}
	}
}

func (r *registry) teardown() {
	log.Println("registry closing shop")
}

func (r *registry) register(id string, gsession *gsession) {
	r.gsessions[id] = gsession
}

func (r *registry) unregister(id string) {
	delete(r.gsessions, id)
}
