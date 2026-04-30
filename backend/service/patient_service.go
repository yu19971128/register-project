package service

import (
	"fmt"

	"clinic/models"
	"clinic/repo"
	"clinic/utils"
)

type PatientService struct {
	repo *repo.PatientRepository
}

func NewPatientService(repo *repo.PatientRepository) *PatientService {
	return &PatientService{repo: repo}
}

func (s *PatientService) CreatePatient(p *models.Patient) (*models.Patient, error) {
	encID, err := utils.Encrypt(p.IDCard)
	if err != nil {
		return nil, fmt.Errorf("encrypt id_card: %w", err)
	}
	encPhone, err := utils.Encrypt(p.Phone)
	if err != nil {
		return nil, fmt.Errorf("encrypt phone: %w", err)
	}
	p.IDCardEncrypted = encID
	p.PhoneEncrypted = encPhone
	p.IDCard = utils.MaskIDCard(p.IDCard)
	p.Phone = utils.MaskPhone(p.Phone)

	id, err := s.repo.Create(p)
	if err != nil {
		return nil, fmt.Errorf("create patient: %w", err)
	}
	p.ID = id
	return p, nil
}

func (s *PatientService) ListPatientsByVisitorPhone(vp string, page, pageSize int) ([]*models.Patient, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.ListByVisitorPhone(vp, (page-1)*pageSize, pageSize)
}

func (s *PatientService) ListPatients(keyword string, page, pageSize int) ([]*models.Patient, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(keyword, (page-1)*pageSize, pageSize)
}

func (s *PatientService) GetPatient(id int64) (*models.Patient, error) {
	return s.repo.GetByID(id)
}

func (s *PatientService) UpdatePatient(p *models.Patient) error {
	if p.IDCard != "" {
		encID, err := utils.Encrypt(p.IDCard)
		if err != nil {
			return fmt.Errorf("encrypt id_card: %w", err)
		}
		p.IDCardEncrypted = encID
		p.IDCard = utils.MaskIDCard(p.IDCard)
	}
	if p.Phone != "" {
		encPhone, err := utils.Encrypt(p.Phone)
		if err != nil {
			return fmt.Errorf("encrypt phone: %w", err)
		}
		p.PhoneEncrypted = encPhone
		p.Phone = utils.MaskPhone(p.Phone)
	}
	return s.repo.Update(p)
}

func (s *PatientService) DeletePatient(id int64) error {
	has, err := s.repo.HasActiveOrders(id)
	if err != nil {
		return fmt.Errorf("check active orders: %w", err)
	}
	if has {
		return fmt.Errorf("该就诊人存在未完成的挂号订单，禁止删除。未完成状态包括：confirmed（待就诊）")
	}
	return s.repo.Delete(id)
}
