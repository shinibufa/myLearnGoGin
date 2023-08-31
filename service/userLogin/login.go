package userlogin

import (
	"errors"

	"github.com/xiongjsh/learn_tiktok_project/middleware"
	"github.com/xiongjsh/learn_tiktok_project/repository"
)

type LoginFlow struct {
	username string
	password string 

	userId int64
	token string
	responce *repository.LoginResponse
}

func NewLoginFlow(username, password string) *LoginFlow {
	return &LoginFlow{
		username: username,
		password: password,
	}
}

func Login(username, pasword string) (*repository.LoginResponse, error) {
	return NewLoginFlow(username, pasword).Do()
}

func (f *LoginFlow) Do() (*repository.LoginResponse, error) {
	if err := f.checkPara(); err != nil {
		return nil, err
	}
	if err := f.prepareData(); err != nil {
		return nil, err
	}
	if err := f.packData(); err != nil {
		return nil, err
	}
	return f.responce, nil
}

func (f *LoginFlow) checkPara() error {
	if len(f.username) < 5 || len(f.username) > 50 {
		return errors.New("invalid username")
	}
	if len(f.password) < 6 || len(f.password) > 20 {
		return errors.New("invalid password")
	}
	return nil
}

func (f *LoginFlow) prepareData() error {
	var userLogin repository.UserLogin
	err := repository.NewUserLoginDao().QueryUserLogin(f.username, f.password, &userLogin)
	if err != nil {
		return err
	}
	f.userId = userLogin.UserInfoId
	token, err := middleware.GenerateToken(userLogin)
	if err != nil {
		return err
	}
	f.token = token
	return nil
}

func (f *LoginFlow) packData() error {
	f.responce = &repository.LoginResponse{
		UserId: f.userId,
		Token: f.token,
	}
	return nil
}

