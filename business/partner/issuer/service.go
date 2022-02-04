package issuer

import (
	"encoding/json"

	partnerIssuerPort "github.com/sepulsa/teleco/business/partner/issuer/port"
)

type (
	service struct {
		partnerIssuerRepository partnerIssuerPort.Repository
	}
)

var (
	ErrDuplicateCode = "Code already in use"
)

func New(partnerIssuerRepository partnerIssuerPort.Repository) partnerIssuerPort.Service {
	return &service{
		partnerIssuerRepository,
	}
}

func (s *service) CreateData(partnerIssuer partnerIssuerPort.PartnerIssuerService) error {
	data := partnerIssuerPort.PartnerIssuerRepo{
		PartnerId: partnerIssuer.PartnerId,
		IssuerId:  partnerIssuer.IssuerId,
		Config:    partnerIssuer.Config,
	}
	return s.partnerIssuerRepository.CreateData(data)
}

func (s *service) ReadData(ID string) (partnerIssuer partnerIssuerPort.PartnerIssuerService, err error) {
	data, err := s.partnerIssuerRepository.ReadData(ID)
	if err != nil {
		return
	}
	partnerIssuer = partnerIssuerPort.PartnerIssuerService{
		ID:        data.ID,
		PartnerId: data.PartnerId,
		IssuerId:  data.IssuerId,
		Config:    data.Config,
	}
	return
}

func (s *service) UpdateData(partnerIssuer partnerIssuerPort.PartnerIssuerService) error {
	_, err := s.partnerIssuerRepository.ReadData(partnerIssuer.ID)
	if err != nil {
		return err
	}
	data := partnerIssuerPort.PartnerIssuerRepo{
		ID:        partnerIssuer.ID,
		PartnerId: partnerIssuer.PartnerId,
		IssuerId:  partnerIssuer.IssuerId,
		Config:    partnerIssuer.Config,
	}
	return s.partnerIssuerRepository.UpdateData(data)
}

func (s *service) DeleteData(ID string) error {
	return s.partnerIssuerRepository.DeleteData(ID)
}

func (s *service) ListData() (partnerIssuers []partnerIssuerPort.PartnerIssuerService, err error) {
	datas, err := s.partnerIssuerRepository.ListData()
	if err != nil {
		return
	}
	if len(datas) > 0 {
		d, _ := json.Marshal(datas)
		json.Unmarshal(d, &partnerIssuers)
	}

	return
}
