package adapters

type Adapters struct {
	baseUrl	 string
}

// NewModels returns models with db pool
func NewAdapter(url string) Adapters {
	return Adapters{
		baseUrl: url,
	}
}
