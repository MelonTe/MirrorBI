package user

import (
	"mrbi/internal/common"
	uservo "mrbi/internal/model/vo/user"
)

type ListUserVOResponse struct {
	common.PageResponse
	Records []uservo.UserVO `json:"records" `
}
