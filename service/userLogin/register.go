package userlogin

import (
	"errors"

	"github.com/xiongjsh/learn_tiktok_project/middleware"
	"github.com/xiongjsh/learn_tiktok_project/repository"
)

type RegisterFlow struct {
	username string
	password string

	userId int64
	Token string
	Response *repository.LoginResponse
}

func NewRegisterFlow(username, password string) *RegisterFlow {
	return &RegisterFlow{
		username: username,
		password: password,
	}
}

func Register(username, password string) (*repository.LoginResponse, error) {
	return NewRegisterFlow(username, password).Do()
}

func (f *RegisterFlow) Do() (*repository.LoginResponse, error) {
	if err := f.checkPara(); err != nil {
		return nil, err
	}
	if err := f.updateDB(); err != nil {
		return nil, err
	}
	if err := f.packResponse(); err != nil {
		return nil, err
	}
	return f.Response, nil
}

func (f *RegisterFlow) checkPara() error {
	if len(f.username) < 5 || len(f.username) > 50 {
		return errors.New("length of username should be between 5 and 50")
	}
	if len(f.password) < 6 || len(f.password) > 20 {
		return errors.New("length of password should be between 6 and 20")
	}
	return nil
}

func (f *RegisterFlow) updateDB() error {
	userLogin := repository.UserLogin{
		Username: f.username,
		Password: f.password,
	}

	userInfo := &repository.UserInfo{
		Name: f.username,
		UserLogin: &userLogin,
	}
	userLoginDao := repository.NewUserLoginDao()
	if userLoginDao.CheckExistByName(f.username) {
		return errors.New("username already exists, please change another one")
	}
	userInfoDao := repository.NewUserInfoDao()
	if err := userInfoDao.AddUserInfo(*&userInfo); err != nil {
		return err
	}
	token, err := middleware.GenerateToken(userLogin)
	if err == nil {
		return err
	}
	f.userId = userInfo.Id
	f.Token = token
	return nil
}

func (f *RegisterFlow) packResponse() error {
	f.Response = &repository.LoginResponse{
		UserId: f.userId,
		Token: f.Token,
	}
	return nil
}