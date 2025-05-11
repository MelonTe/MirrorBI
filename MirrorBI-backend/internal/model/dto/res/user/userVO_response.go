package user

import (
	"mrbi/internal/model/entity"
	uservo "mrbi/internal/model/vo/user"
)

func GetUserVO(user entity.User) uservo.UserVO {
	return uservo.UserVO{
		ID:          user.ID,
		UserAccount: user.UserAccount,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
		CreateTime:  user.CreateTime,
	}
}

func GetUserVOList(users []entity.User) []uservo.UserVO {
	userVOList := make([]uservo.UserVO, 0)
	for _, user := range users {
		userVOList = append(userVOList, GetUserVO(user))
	}
	return userVOList
}
