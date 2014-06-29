package gusher

import (
	"log"
)

type registry struct {
	gsessions map[*gsession]bool
	add       chan *gsession
	remove    chan *gsession
}

func newRegistry() *registry {
	registry := &registry{
		gsessions: make(map[*gsession]bool),
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
			r.register(gsession)
		case gsession := <-r.remove:
			r.unregister(gsession)
		}
	}
}

func (r *registry) teardown() {
	log.Println("registry closing shop")
}

func (r *registry) register(gsession *gsession) {
	r.gsessions[gsession] = true
}

func (r *registry) unregister(gsession *gsession) {
	gsession.close()
	delete(r.gsessions, gsession)
}
