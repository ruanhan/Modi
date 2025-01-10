package service

import (
	"fmt"
	UserModel "github.com/bigartists/Modi/src/model/UserModel"
	"github.com/bigartists/Modi/src/repo"
	"github.com/bigartists/Modi/src/result"
)

type IUser interface {
	GetUserList() []*UserModel.UserImpl
	GetUserDetail(id int64) *result.ErrorResult
	CreateUser(user *UserModel.UserImpl) *result.ErrorResult
	UpdateUser(id int, user *UserModel.UserImpl) *result.ErrorResult
	DeleteUser(id int) *result.ErrorResult
	SignIn(username string, password string) (*UserModel.UserImpl, error)
	SignUp(email string, username string, password string) error
}

type IUserServiceGetterImpl struct {
	userRepo repo.IUserRepo
}

func NewUserServiceImpl(userRepo repo.IUserRepo) *IUserServiceGetterImpl {
	return &IUserServiceGetterImpl{userRepo: userRepo}
}

func (this *IUserServiceGetterImpl) SignIn(username string, password string) (*UserModel.UserImpl, error) {
	user, err := this.userRepo.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}
	//if user.Password != password {
	//	err = fmt.Errorf("用户名%s或密码错误", username)
	//	return nil, err
	//}
	// 校验密码
	if !user.CheckPassword(password) {
		err = fmt.Errorf("用户名%s或密码错误", username)
		return nil, err
	}
	return user, nil
}

func (this *IUserServiceGetterImpl) SignUp(email string, username string, password string) error {
	//return repo.UserGetter.CreateUser(user)

	if _, err := this.userRepo.FindUserByUsername(username); err != nil {
		return fmt.Errorf("用户名%s已存在", username)
	}
	if _, err := this.userRepo.FindUserByEmail(email); err != nil {
		return fmt.Errorf("邮箱%s已存在", email)
	}

	user := UserModel.New(UserModel.WithEmail(email), UserModel.WithUsername(username), UserModel.WithPassword(password))
	err := user.GeneratePassword()
	if err != nil {
		return fmt.Errorf("密码加密失败")
	}
	err = this.userRepo.CreateUser(user)

	if err != nil {
		return fmt.Errorf("用户注册失败")
	}

	return nil
}

func (this *IUserServiceGetterImpl) GetUserList() []*UserModel.UserImpl {
	users := this.userRepo.FindUserAll()
	return users
}

func (this *IUserServiceGetterImpl) GetUserDetail(id int64) *result.ErrorResult {
	//TODO implement me
	user := UserModel.New()
	_, err := this.userRepo.FindUserById(id, user)
	if err != nil {
		return result.Result(nil, err)
	}
	return result.Result(user, nil)
}

func (this *IUserServiceGetterImpl) CreateUser(user *UserModel.UserImpl) *result.ErrorResult {
	//TODO implement me
	panic("implement me")
}

func (this *IUserServiceGetterImpl) UpdateUser(id int, user *UserModel.UserImpl) *result.ErrorResult {
	//TODO implement me
	panic("implement me")
}

func (this *IUserServiceGetterImpl) DeleteUser(id int) *result.ErrorResult {
	//TODO implement me
	panic("implement me")
}

//
//// 创建用户
//func (this *IUserGetterImpl) CreateUser(user *UserModel.UserModelImpl) *result.ErrorResult {
//	db := dbs.Orm.Create(user)
//	if db.Error != nil {
//		return result.Result(nil, db.Error)
//	}
//	return result.Result(user, nil)
//}
//

//
//// 更新用户
//func (this *IUserGetterImpl) UpdateUser(id int, user *UserModel.UserModelImpl) *result.ErrorResult {
//	db := dbs.Orm.Where("id=?", id).Updates(user)
//	if db.Error != nil {
//		return result.Result(nil, db.Error)
//	}
//	return result.Result(user, nil)
//}
//
//// 删除用户
//func (this *IUserGetterImpl) DeleteUser(id int) *result.ErrorResult {
//	user := UserModel.New()
//	db := dbs.Orm.Where("id=?", id).Delete(user)
//	if db.Error != nil {
//		return result.Result(nil, db.Error)
//	}
//	return result.Result(user, nil)
//}
