package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"
)

type UserService struct {
}

func (userService *UserService) Register(u database.User) (database.User, error) {
	//只有当结果是「没查到用户」→ 条件为 false（不返回错误，继续注册）
	if !errors.Is(global.DB.Where("email = ?", u.Email).First(&database.User{}).Error, gorm.ErrRecordNotFound) {
		return database.User{}, errors.New("this email address is already registered, please check the information you filled in, or retrieve your password")
	}
	//对用户注册信息进行预处理
	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.Must(uuid.NewV4())
	u.Avatar = "/image/avatar.jpg"
	u.RoleID = appTypes.User
	u.Register = appTypes.Email
	if err := global.DB.Create(&u).Error; err != nil {
		return database.User{}, err
	}
	return u, nil
}

func (userService *UserService) EmailLogin(u database.User) (database.User, error) {
	var user database.User
	//First(&user),这里没用空结构体，就是为了让查询结果覆盖
	err := global.DB.Where("email = ?", u.Email).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return database.User{}, errors.New("incorrect email or password")
		}
		//数据库查询成功（err == nil），即找到了该邮箱对应的用户
		return user, nil
	}
	//登录失败，返回一个空的用户对象（database.User{}），并携带具体的错误信息（err）。
	return database.User{}, err
}

func (userService *UserService) QQLogin(accessTokenResponse other.AccessTokenResponse) (database.User, error) {
	var user database.User

	// 尝试查找用户
	err := global.DB.Where("openid = ?", accessTokenResponse.Openid).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return database.User{}, err
	}

	// 如果用户不存在，则创建新用户
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userInfoResponse, err := ServiceGroupApp.QQService.GetUserInfoByAccessTokenAndOpenid(accessTokenResponse.AccessToken, accessTokenResponse.Openid)
		if err != nil {
			return database.User{}, err
		}
		user.UUID = uuid.Must(uuid.NewV4())
		user.Username = userInfoResponse.Nickname
		user.Openid = accessTokenResponse.Openid
		user.Avatar = userInfoResponse.FigureurlQQ2
		user.RoleID = appTypes.User
		user.Register = appTypes.QQ
		if err := global.DB.Create(&user).Error; err != nil {
			return database.User{}, err
		}
	}
	return user, nil
}

func (userService *UserService) ForgetPassword(req request.ForgetPassword) error {
	var user database.User
	err := global.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		return err
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userService *UserService) UserCard(req request.UserCard) (response.UserCard, error) {
	var user database.User
	if err := global.DB.Where("uuid = ?", req.UUID).Select("uuid", "username", "avatar", "address", "signature").First(&user).Error; err != nil {
		return response.UserCard{}, err
	}
	return response.UserCard{
		UUID:      user.UUID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Address:   user.Address,
		Signature: user.Signature,
	}, nil
}

func (userService *UserService) Logout(c *gin.Context) {
	//先拿到uuid
	logoutUuid := utils.GetUUID(c)
	//清理 refresh token
	jwtStr := utils.GetRefreshToken(c)
	utils.ClearRefreshToken(c)
	//从redis里面删除这个用户
	global.Redis.Del(logoutUuid.String())
	//加入黑名单
	_ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{Jwt: jwtStr})
}

func (userService *UserService) UserResetPassword(req request.UserResetPassword) error {
	//拿到用户id
	var user database.User
	err := global.DB.Take(&user, req.UserID).Error
	if err != nil {
		return err
	}
	//实现重新设置密码逻辑
	//1.对比旧密码
	if !utils.BcryptCheck(req.Password, user.Password) {
		return errors.New("original password does not match the current account")
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userService *UserService) UserInfo(userID uint) (database.User, error) {
	var user database.User
	//第一个参数 &user：将查询结果映射到 user 变量（需要传递指针才能修改原变量）
	//userID是参数名，和数据库字段名无直接关联
	err := global.DB.Take(&user, userID).Error
	if err != nil {
		return database.User{}, err
	}
	return user, nil
}

func (userService *UserService) UserChangeInfo(req request.UserChangeInfo) error {
	var user database.User
	err := global.DB.Take(&user, req.UserID).Error
	if err != nil {
		return err
	}
	return global.DB.Model(&user).Updates(req).Error
}

func (userService *UserService) UserWeather(ip string) (string, error) {
	//先在redis查询
	result, err := global.Redis.Get("weather-" + ip).Result()
	if err != nil {
		//如果没有数据，则调用高德api进行查询
		ipResponse, err := ServiceGroupApp.GaodeService.GetLocationByIP(ip)
		if err != nil {
			return "", err
		}
		live, err := ServiceGroupApp.GaodeService.GetWeatherByAdcode(ipResponse.Adcode)
		if err != nil {
			return "", err
		}
		weather := "地区：" + live.Province + "-" + live.City + " 天气：" + live.Weather + " 温度：" + live.Temperature + "°C" + " 风向：" + live.WindDirection + " 风级：" + live.WindPower + " 湿度：" + live.Humidity + "%"
		// 将天气数据存入redis
		if err := global.Redis.Set("weather-"+ip, weather, time.Hour*1).Err(); err != nil {
			return "", err
		}
		return weather, nil
	}
	return result, nil
}

func (userService *UserService) UserChart(req request.UserChart) (response.UserChart, error) {
	//先在数据库查讯数据
	//where变量是携带 "创建时间在最近 N 天内" 的数据，后续可复用
	where := global.DB.Where(fmt.Sprintf("date_sub(curdate(), interval %d day) <= created_at", req.Date))

	var resp response.UserChart
	// 生成日期列表
	startDate := time.Now().AddDate(0, 0, -req.Date)
	for i := 1; i <= req.Date; i++ {
		resp.DateList = append(resp.DateList, startDate.AddDate(0, 0, i).Format("2006-01-02"))
	}
	//获取注册数据
	registerCounts := utils.FetchDateCounts(global.DB.Model(&database.User{}), where)
	//获取登录数据
	loginCounts := utils.FetchDateCounts(global.DB.Model(&database.Login{}), where)
	for _, data := range resp.DateList {
		//依照日期信息一一对应
		registerCount := registerCounts[data]
		loginCount := loginCounts[data]
		resp.LoginData = append(resp.LoginData, loginCount)
		resp.RegisterData = append(resp.RegisterData, registerCount)
	}
	return resp, nil
}

func (userService *UserService) UserList(req request.UserList) (interface{}, int64, error) {
	//两条查询数据（uuid和注册来源 ）
	db := global.DB
	if req.UUID != nil {
		db = db.Where("uuid = ?", req.UUID)
	}
	if req.Register != nil {
		db = db.Where("register = ?", req.Register)
	}
	option := other.MySQLOption{
		PageInfo: req.PageInfo,
		Where:    db,
	}
	return utils.MySQLPagination(&database.User{}, option)
}

func (userService *UserService) UserFreeze(req request.UserOperation) error {
	//首先拿到id
	var user database.User
	//用 Update方法修改结构体属性
	err := global.DB.Take(&user, req.ID).Update("freeze", true).Error
	if err != nil {
		return err
	}
	jwtStr, _ := ServiceGroupApp.JwtService.GetRedisJWT(user.UUID)
	if jwtStr != "" {
		_ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{Jwt: jwtStr})
	}
	return nil
}

func (userService *UserService) UserUnfreeze(req request.UserOperation) error {
	var user database.User
	return global.DB.Take(&user, req.ID).Update("freeze", false).Error
}

func (userService *UserService) UserLoginList(req request.UserLoginList) (interface{}, int64, error) {
	db := global.DB
	var userID uint
	//uuid转换成userid
	if req.UUID != nil {
		if err := global.DB.Model(&database.User{}).Where("uuid = ?", *req.UUID).Pluck("id", &userID).Error; err != nil {
			return nil, 0, err
		}
		db = db.Where("user_id = ?", userID)
	}
	option := other.MySQLOption{
		PageInfo: req.PageInfo,
		Where:    db,
		Preload:  []string{"User"},
	}
	return utils.MySQLPagination(&database.Login{}, option)
}
