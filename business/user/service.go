package user

import (
	"errors"
	"strings"

	authPort "github.com/sepulsa/teleco/business/auth/port"
	userPort "github.com/sepulsa/teleco/business/user/port"
	"github.com/sepulsa/teleco/utils/crypto"
)

type (
	service struct {
		userRepository      userPort.Repository
		userTokenRepository authPort.Repository
	}
)

var (
	ErrDuplicateEmail       = "Email already in use"
	ErrUserGeneratePassword = "generate password failed"
)

func New(userRepository userPort.Repository, userTokenRepository authPort.Repository) userPort.Service {
	return &service{
		userRepository,
		userTokenRepository,
	}
}

func (s *service) CreateData(user userPort.UserService) error {
	existingIssuer := s.userRepository.FindByEmail(user.Email)
	if existingIssuer.ID != "" {
		return errors.New(ErrDuplicateEmail)
	}

	var data userPort.UserRepo

	data.Email = user.Email
	data.Fullname = user.Fullname

	if strings.TrimSpace(user.Password) != "" {
		generatePassword(user.Password, &data)
	}

	return s.userRepository.CreateData(data)
}

func (s *service) ReadData(ID string) (userPort.UserService, error) {
	var user userPort.UserService

	data, err := s.userRepository.ReadData(ID)
	if err == nil {
		user.ID = data.ID
		user.Email = data.Email
		user.Fullname = data.Fullname
	}
	return user, err
}

func (s *service) UpdateData(user userPort.UserService) error {
	existingData, err := s.userRepository.ReadData(user.ID)
	if err != nil {
		return err
	}
	if existingData.Email != user.Email {
		existingUser := s.userRepository.FindByEmail(user.Email)
		if existingUser.ID != "" {
			return errors.New(ErrDuplicateEmail)
		}
	}

	var data userPort.UserRepo

	data.ID = user.ID
	data.Email = user.Email
	data.Fullname = user.Fullname

	if strings.TrimSpace(user.Password) != "" {
		generatePassword(user.Password, &data)
	}

	return s.userRepository.UpdateData(data)
}

func (s *service) DeleteData(ID string) error {
	err := s.userRepository.DeleteData(ID)
	if err == nil {
		_ = s.userTokenRepository.DeleteDataByUserID(ID)
	}
	return err
}

func (s *service) ListData() ([]userPort.UserService, error) {
	users := make([]userPort.UserService, 0)

	datas, err := s.userRepository.ListData()
	if err != nil {
		return users, err
	}

	var user userPort.UserService
	for i := range datas {
		user.ID = datas[i].ID
		user.Email = datas[i].Email
		user.Fullname = datas[i].Fullname

		users = append(users, user)
	}

	return users, nil
}

func generatePassword(stringPwd string, data *userPort.UserRepo) {
	hashedPassword, _ := crypto.UserGeneratePassword(stringPwd)
	data.Password = hashedPassword
}
