package provider

type Provider interface {
	Request(string) (string, error)
}
