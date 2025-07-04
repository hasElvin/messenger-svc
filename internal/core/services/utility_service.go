package services

import (
	"github.com/hasElvin/messenger-svc/internal/core/ports"
)

type utilityService struct {
	repo ports.MessageRepository
}

func NewUtilityService(repo ports.MessageRepository) ports.UtilityService {
	return &utilityService{
		repo: repo,
	}
}

func (s utilityService) SeedSampleMessages() error {
	err := s.repo.SeedSampleMessages()
	if err != nil {
		return err
	}
	return nil
}

func (s utilityService) ClearDatabase() error {
	err := s.repo.ClearDatabase()
	if err != nil {
		return err
	}
	return nil
}
