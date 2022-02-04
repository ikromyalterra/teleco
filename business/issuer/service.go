package issuer

import (
	"encoding/json"
	"errors"

	issuerPort "github.com/sepulsa/teleco/business/issuer/port"
)

type (
	service struct {
		issuerRepository issuerPort.Repository
	}
)

var (
	ErrDuplicateCode = "Code already in use"
)

func New(issuerRepository issuerPort.Repository) issuerPort.Service {
	return &service{
		issuerRepository,
	}
}

func (s *service) CreateData(issuer issuerPort.IssuerService) error {
	existingIssuer := s.issuerRepository.FindByCode(issuer.Code)
	if existingIssuer.ID != "" {
		return errors.New(ErrDuplicateCode)
	}

	data := issuerPort.IssuerRepo{
		Code:          issuer.Code,
		Label:         issuer.Label,
		Config:        issuer.Config,
		ThreadNum:     issuer.ThreadNum,
		ThreadTimeout: issuer.ThreadTimeout,
	}
	return s.issuerRepository.CreateData(data)
}

func (s *service) ReadData(ID string) (issuer issuerPort.IssuerService, err error) {
	data, err := s.issuerRepository.ReadData(ID)
	if err != nil {
		return
	}
	issuer = issuerPort.IssuerService{
		ID:            data.ID,
		Code:          data.Code,
		Label:         data.Label,
		Config:        data.Config,
		ThreadNum:     data.ThreadNum,
		ThreadTimeout: data.ThreadTimeout,
	}
	return
}

func (s *service) UpdateData(issuer issuerPort.IssuerService) error {
	existingData, err := s.issuerRepository.ReadData(issuer.ID)
	if err != nil {
		return err
	}
	if existingData.Code != issuer.Code {
		existingIssuer := s.issuerRepository.FindByCode(issuer.Code)
		if existingIssuer.ID != "" {
			return errors.New(ErrDuplicateCode)
		}
	}
	data := issuerPort.IssuerRepo{
		ID:            issuer.ID,
		Code:          issuer.Code,
		Label:         issuer.Label,
		Config:        issuer.Config,
		ThreadNum:     issuer.ThreadNum,
		ThreadTimeout: issuer.ThreadTimeout,
	}
	return s.issuerRepository.UpdateData(data)
}

func (s *service) DeleteData(ID string) error {
	return s.issuerRepository.DeleteData(ID)
}

func (s *service) ListData() (issuers []issuerPort.IssuerService, err error) {
	datas, err := s.issuerRepository.ListData()
	if err != nil {
		return
	}
	if len(datas) > 0 {
		d, _ := json.Marshal(datas)
		json.Unmarshal(d, &issuers)
	}

	return
}
