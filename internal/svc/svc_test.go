package svc

import (
	"testing"
)

type mockService struct{}

func (s *mockService) Init() error  { return nil }
func (s *mockService) Start() error { return nil }
func (s *mockService) Stop() error  { return nil }

func TestRun(t *testing.T) {
	s := &mockService{}
	if err := Run(s); err != nil {
		t.Fail()
	}
}
