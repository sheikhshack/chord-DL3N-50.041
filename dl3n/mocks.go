package dl3n

type MockNodeDiscovery struct {
	addr string
}

func NewMockNodeDiscovery(addr string) *MockNodeDiscovery {
	return &MockNodeDiscovery{
		addr: addr,
	}
}

func (m *MockNodeDiscovery) FindSeederAddr(_ string) string {
	return m.addr
}
