package storage

type StorageOptions struct {
	Addr     string
	Password string
	DB       int
}

type StorageOptionFunc = func(*StorageOptions)

func DefaultOptions(opts ...StorageOptionFunc) *StorageOptions {
	defaultOpts := &StorageOptions{}

	for _, opt := range opts {
		opt(defaultOpts)
	}

	return defaultOpts
}

func WithAddr(addr string) StorageOptionFunc {
	return func(o *StorageOptions) {
		o.Addr = addr
	}
}

func WithPassword(password string) StorageOptionFunc {
	return func(o *StorageOptions) {
		o.Password = password
	}
}

func WithDB(db int) StorageOptionFunc {
	return func(o *StorageOptions) {
		o.DB = db
	}
}
