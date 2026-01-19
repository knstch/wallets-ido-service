package enum

import "fmt"

type Provider string

func (r Provider) String() string {
	return string(r)
}

const (
	ProviderPhantom Provider = "phantom"
)

func GetProvider(provider string) (Provider, error) {
	switch provider {
	case "phantom":
		return ProviderPhantom, nil
	default:
		return "", fmt.Errorf("unknown provider: %s", provider)
	}
}
