package partner

import (
	"encoding/json"

	partnerPort "github.com/sepulsa/teleco/business/partner/port"
)

type (
	service struct {
		partnerRepository partnerPort.Repository
	}
)

var (
	ErrDuplicateName = "Name already in use"
)

func New(partnerRepository partnerPort.Repository) partnerPort.Service {
	return &service{
		partnerRepository,
	}
}

func (s *service) CreateData(partner partnerPort.PartnerService) error {
	data := partnerPort.PartnerRepo{
		Code:        partner.Code,
		Name:        partner.Name,
		Pic:         partner.Pic,
		Address:     partner.Address,
		CallbackUrl: partner.CallbackUrl,
		IpWhitelist: partner.IpWhitelist,
		Status:      partner.Status,
		SecretKey:   partner.SecretKey,
	}
	return s.partnerRepository.CreateData(data)
}

func (s *service) ReadData(ID string) (partner partnerPort.PartnerService, err error) {
	data, err := s.partnerRepository.ReadData(ID)
	if err != nil {
		return
	}

	partner = partnerPort.PartnerService{
		ID:          data.ID,
		Code:        data.Code,
		Name:        data.Name,
		Pic:         data.Pic,
		Address:     data.Address,
		CallbackUrl: data.CallbackUrl,
		IpWhitelist: data.IpWhitelist,
		Status:      data.Status,
		SecretKey:   data.SecretKey,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.CreatedAt,
		DeletedAt:   data.DeletedAt,
	}
	return
}

func (s *service) UpdateData(partner partnerPort.PartnerService) error {
	_, err := s.partnerRepository.ReadData(partner.ID)
	if err != nil {
		return err
	}
	data := partnerPort.PartnerRepo{
		ID:          partner.ID,
		Code:        partner.Code,
		Name:        partner.Name,
		Pic:         partner.Pic,
		Address:     partner.Address,
		CallbackUrl: partner.CallbackUrl,
		IpWhitelist: partner.IpWhitelist,
		Status:      partner.Status,
		SecretKey:   partner.SecretKey,
		CreatedAt:   partner.CreatedAt,
		UpdatedAt:   partner.CreatedAt,
		DeletedAt:   partner.DeletedAt,
	}
	return s.partnerRepository.UpdateData(data)
}

func (s *service) DeleteData(ID string) error {
	return s.partnerRepository.DeleteData(ID)
}

func (s *service) ListData() (partners []partnerPort.PartnerService, err error) {
	datas, err := s.partnerRepository.ListData()
	if err != nil {
		return
	}
	if len(datas) > 0 {
		d, _ := json.Marshal(datas)
		json.Unmarshal(d, &partners)
	}

	return
}
