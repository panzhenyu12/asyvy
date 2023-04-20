package model

type ImageScanReq struct {
	// 镜像地址
	Image string `json:"image" binding:"required"`
	// 镜像tag
	Tag string `json:"tag" binding:"required"`
	// 镜像仓库用户名
	Username string `json:"username"`
	// 镜像仓库密码
	Password string `json:"password"`
	// 镜像仓库地址
	//Registry string `json:"registry"`
}

type ImageScanResp struct {
	BaseResp
	TaskID string `json:"task_id"`
}
