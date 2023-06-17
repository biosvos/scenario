package scenario

type Options struct {
	artifacts *Artifacts
}

type Option func(*Options) error

func WithArtifacts(item map[string]any) Option {
	return func(options *Options) error {
		for key, value := range item {
			err := options.artifacts.Add(key, value)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func WithArtifact(key string, value any) Option {
	return func(options *Options) error {
		return options.artifacts.Add(key, value)
	}
}
