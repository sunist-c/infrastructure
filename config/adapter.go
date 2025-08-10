package config

type Loader interface {
	Load(source, namespace string, receiver any) error
	MustLoad(source, namespace string, receiver any)
}
