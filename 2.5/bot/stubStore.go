package bot

// Stub store for testing.
type stubStore struct {
	findSubsBySchedule func() ([]*Subscription, error)
}

func NewStubStore() *stubStore {
	return &stubStore{}
}

func (s *stubStore) Save(_ Update) error {
	panic("not implemented") // TODO: Implement
}

func (s *stubStore) SetFindSubsBySchedule(f func() ([]*Subscription, error)) {
	s.findSubsBySchedule = f
}
func (s *stubStore) FindSubsBySchedule(_ string) ([]*Subscription, error) {
	return s.findSubsBySchedule()
}

func (s *stubStore) SaveLocation(_ *Subscription) error {
	panic("not implemented") // TODO: Implement
}

func (s *stubStore) SaveSchedule(_ string, _ string) error {
	panic("not implemented") // TODO: Implement
}

func (s *stubStore) GetSubscriptionByChatId(_ int) *Subscription {
	panic("not implemented") // TODO: Implement
}

func (s *stubStore) GetSubscription(_ string) (*Subscription, error) {
	panic("not implemented") // TODO: Implement
}
