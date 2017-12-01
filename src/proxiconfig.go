package main

import "persistence"

// ProxiConfig acts as a service level config
// Allows customization of service listener port and service backend
type ProxiConfig struct {
	Port                int                              `json:"port"`
	PersistenceConfig   *persistence.PersistenceConfig   `json:"persistence"`
}

