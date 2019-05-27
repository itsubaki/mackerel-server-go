package usecase

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type MonitorInteractor struct {
	MonitorRepository MonitorRepository
}

func (s *MonitorInteractor) List(org string) (*domain.Monitors, error) {
	return s.MonitorRepository.List(org)
}

func (s *MonitorInteractor) Save(org string, monitor *domain.Monitoring) (interface{}, error) {
	sha := sha256.Sum256([]byte(uuid.Must(uuid.NewRandom()).String()))
	hash := hex.EncodeToString(sha[:])
	monitor.ID = hash[:11]

	return s.MonitorRepository.Save(org, monitor)
}

func (s *MonitorInteractor) Update(org string, monitor *domain.Monitoring) (interface{}, error) {
	return s.MonitorRepository.Update(org, monitor)
}

func (s *MonitorInteractor) Monitor(org, monitorID string) (interface{}, error) {
	return s.MonitorRepository.Monitor(org, monitorID)
}

func (s *MonitorInteractor) Delete(org, monitorID string) (interface{}, error) {
	return s.MonitorRepository.Delete(org, monitorID)
}
