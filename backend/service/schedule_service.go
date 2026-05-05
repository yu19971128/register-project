package service

import (
	"fmt"

	"clinic/models"
	"clinic/repo"
)

type ScheduleService struct {
	repo *repo.ScheduleRepository
}

func NewScheduleService(repo *repo.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) CreateSchedule(schedule *models.Schedule) (*models.Schedule, error) {
	schedule.Remaining = schedule.TotalQuota
	schedule.Status = "available"
	id, err := s.repo.Create(schedule)
	if err != nil {
		return nil, fmt.Errorf("create schedule: %w", err)
	}
	schedule.ID = id
	return schedule, nil
}

func (s *ScheduleService) ListSchedules(date, department, doctorName string, page, pageSize int) ([]*models.Schedule, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize == 0 {
		return s.repo.List(date, department, doctorName, 0, 0)
	}
	if pageSize < 0 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.List(date, department, doctorName, (page-1)*pageSize, pageSize)
}

func (s *ScheduleService) GetSchedule(id int64) (*models.Schedule, error) {
	schedule, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get schedule: %w", err)
	}
	if schedule == nil {
		return nil, fmt.Errorf("号源不存在")
	}
	return schedule, nil
}

func (s *ScheduleService) UpdateSchedule(schedule *models.Schedule) error {
	existing, err := s.repo.GetByID(schedule.ID)
	if err != nil {
		return fmt.Errorf("get schedule: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("号源不存在")
	}
	booked := existing.TotalQuota - existing.Remaining
	if schedule.TotalQuota > 0 && schedule.TotalQuota < booked {
		return fmt.Errorf("总号数不得小于已预约数 %d", booked)
	}
	// Whenever total_quota is provided, recompute remaining from booking count
	// so it stays consistent regardless of whether quota increased or decreased.
	// Otherwise preserve the existing remaining (zero-value would otherwise overwrite it).
	if schedule.TotalQuota > 0 {
		schedule.Remaining = schedule.TotalQuota - booked
	} else {
		schedule.Remaining = existing.Remaining
	}
	return s.repo.Update(schedule)
}

func (s *ScheduleService) DeleteSchedule(id int64) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("get schedule: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("号源不存在")
	}
	if existing.TotalQuota-existing.Remaining > 0 {
		return fmt.Errorf("该号源已有预约，禁止删除")
	}
	return s.repo.Delete(id)
}

func (s *ScheduleService) Deduct(id int64) (bool, error) {
	ok, err := s.repo.Deduct(id)
	if err != nil {
		return false, fmt.Errorf("deduct schedule: %w", err)
	}
	if ok {
		schedule, _ := s.repo.GetByID(id)
		if schedule != nil && schedule.Remaining == 0 {
			_ = s.repo.Update(&models.Schedule{ID: id, Status: "full"})
		}
	}
	return ok, nil
}

func (s *ScheduleService) Rollback(id int64) (bool, error) {
	schedule, err := s.repo.GetByID(id)
	if err != nil {
		return false, fmt.Errorf("get schedule: %w", err)
	}
	if schedule == nil {
		return false, fmt.Errorf("号源不存在")
	}
	if schedule.Remaining >= schedule.TotalQuota {
		return false, nil
	}
	ok, err := s.repo.Rollback(id)
	if err != nil {
		return false, fmt.Errorf("rollback schedule: %w", err)
	}
	if ok && schedule.Status == "full" {
		_ = s.repo.Update(&models.Schedule{ID: id, Status: "available"})
	}
	return ok, nil
}
