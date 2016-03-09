package registry

import (
	"fmt"
)

var (
	registryClients map[string]RegistryClient
)

type RegistryClient interface {
	ListImages() (interface{}, error)
	GetImageTags(string) (interface{}, error)
	GetImageDigest(string, string) (interface{}, error)
	DeleteImageTag(string, string) error
}

func RegisterRegistryClient(name string, client RegistryClient) error {
	if registryClients == nil {
		registryClients = make(map[string]RegistryClient)

	}

	if _, exists := registryClients[name]; exists {
		return fmt.Errorf("RegistryClient already registered")
	}

	registryClients[name] = client
	return nil
}

func GetRegistryClient() (RegistryClient, error) {
	/*
		name := os.Getenv("RegistryClient")
		if name == nil {
			return nil, fmt.Errorf("not support")
		}
	*/
	name := "direct"

	if client, ok := registryClients[name]; ok {
		return client, nil
	}

	return nil, fmt.Errorf("registryClients[%s] not support", name)

}
