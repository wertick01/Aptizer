package processors

import (
	"errors"

	"aptizer.com/internal/app/models"
)

func (processor *Processor) ListUsers() ([]*models.User, error) {
	return processor.store.UsersStorage.GetUsersList()
}

func (processor *Processor) CreateUser(user *models.User) (*models.User, error) {
	if user.Name == "" {
		return nil, errors.New("name should not be empty")
	}
	return processor.store.UsersStorage.CreateNewUser(user)
}

func (processor *Processor) FindUser(id int64) (*models.User, error) {
	user, err := processor.store.UsersStorage.GetUserByID(id)
	if err != nil {
		return user, err
	}
	return user, nil

}

func (processor *Processor) FindByPhone(phone string) (*models.User, error) {
	user, err := processor.store.UsersStorage.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (processor *Processor) SetRefrToken(rt *models.RefreshToken) error {
	return processor.store.UsersStorage.SetRefreshToken(rt)
}

func (processor *Processor) CheckRefrToken(rt *models.RefreshToken) (*models.RefreshToken, error) {
	return processor.store.UsersStorage.CheckRefreshToken(rt)
}

func (processor *Processor) UpdateUser(user *models.User) (*models.User, error) {
	changeduser, err := processor.store.UsersStorage.ChangeUser(user)
	if err != nil {
		return user, err
	}
	return changeduser, nil
}

func (processor *Processor) DeleteUser(id int64) (*models.User, error) {
	user, err := processor.FindUser(id)
	if err != nil {
		return nil, err
	}
	_, err = processor.store.UsersStorage.DeleteUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
