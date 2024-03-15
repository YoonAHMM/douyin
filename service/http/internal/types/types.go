// Code generated by goctl. DO NOT EDIT.
package types

type Comment struct {
	ID         uint64 `json:"id"`          // 评论id
	User       User   `json:"user"`        // 评论用户信息
	Content    string `json:"content"`     // 评论内容
	CreateDate string `json:"create_date"` // 评论发布日期，格式 mm-dd
}

type CommentActionReq struct {
	Token       string  `form:"token"`                 // 用户鉴权 token
	VideoId     string  `form:"video_id"`              // 视频id
	ActionType  string  `form:"action_type"`           // 1-发布评论，2-删除评论
	CommentText *string `form:"comment_text,optional"` // 用户填写的评论内容，在action_type=1的时候使用
	CommentId   *string `form:"comment_id,optional"`   // 要删除的评论id，在action_type=2的时候使用
}

type CommentActionResp struct {
	Status
	Comment *Comment `json:"comment"` // 评论成功返回评论内容，不需要重新拉取整个列表
}

type CommentListReq struct {
	Token   string `form:"token"`    // 用户鉴权 token
	VideoId string `form:"video_id"` // 视频id
}

type CommentListResp struct {
	Status
	CommentList []Comment `json:"comment_list"` // 评论列表
}

type FavoriteActionReq struct {
	Token      string `form:"token"`       // 用户鉴权 token
	VideoId    string `form:"video_id"`    // 视频id
	ActionType string `form:"action_type"` // 1-点赞，2-取消点赞
}

type FavoriteActionResp struct {
	Status
}

type FavoriteListReq struct {
	UserId string `form:"user_id"` // 用户 id
	Token  string `form:"token"`   // 用户鉴权 token
}

type FavoriteListResp struct {
	Status
	VideoList []Video `json:"video_list"` // 用户发布的视频列表
}

type FeedReq struct {
	LatestTime *string `form:"latest_time,optional"` // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	Token      *string `form:"token,optional"`       // 用户登录状态下设置 可选
}

type FeedResp struct {
	Status
	NextTime  *int64  `json:"next_time"`  // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	VideoList []Video `json:"video_list"` // 视频列表
}

type FollowListReq struct {
	Token  string `form:"token"`   // 用户鉴权 token
	UserId string `form:"user_id"` // 用户id
}

type FollowListResp struct {
	Status
	UserList []User `json:"user_list"`
}

type FollowerListReq struct {
	Token  string `form:"token"`   // 用户鉴权 token
	UserId string `form:"user_id"` // 用户id
}

type FollowerListResp struct {
	Status
	UserList []User `json:"user_list"`
}

type GetUserReq struct {
	UserID string `form:"user_id"` // 用户id
	Token  string `form:"token"`   // 用户鉴权token
}

type GetUserResp struct {
	Status
	User *User `json:"user"` // 用户信息
}

type LoginReq struct {
	Username string `form:"username"` // 登入用户名
	Password string `form:"password"` // 登入密码
}

type LoginResp struct {
	Status
	UserID uint64  `json:"user_id,omitempty"` // 用户id
	Token  *string `json:"token"`             // 用户鉴权token
}

type PublishActionReq struct {
	Token string `form:"token"` // 用户鉴权 token
	Title string `form:"title"` // 视频标题
}

type PublishActionResp struct {
	Status
}

type PublishListReq struct {
	Token  string `form:"token"`   // 用户鉴权 token
	UserId string `form:"user_id"` // 用户 id
}

type PublishListResp struct {
	Status
	VideoList []Video `json:"video_list"` // 用户发布的视频列表
}

type RegisterReq struct {
	Username string `form:"username"` // 注册用户名，最长32个字符
	Password string `form:"password"` // 密码，最长32个字符
}

type RegisterResp struct {
	Status
	UserID uint64  `json:"user_id"` // 用户id
	Token  *string `json:"token"`   // 用户鉴权token
}

type RelationActionReq struct {
	Token      string `form:"token"`       // 用户鉴权 token
	ToUserId   string `form:"to_user_id"`  // 对方用户id
	ActionType string `form:"action_type"` // 1-关注，2-取消关注
}

type RelationActionResp struct {
	Status
}

type Status struct {
	StatusCode string `json:"status_code"`          // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg,omitempty"` // 返回状态描述
}

type User struct {
	ID            uint64 `json:"id"`             // 用户id
	Name          string `json:"name"`           // 用户名称
	FollowCount   int64  `json:"follow_count"`   // 关注总数
	FollowerCount int64  `json:"follower_count"` // 粉丝总数
	IsFollow      bool   `json:"is_follow"`      // true-已关注，false-未关注
}

type Video struct {
	ID            uint64 `json:"id"`             // 视频唯一标识
	Author        User   `json:"author"`         // 视频作者信息
	PlayURL       string `json:"play_url"`       // 视频播放地址
	CoverURL      string `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"` // 视频的点赞总数
	CommentCount  int64  `json:"comment_count"`  // 视频的评论总数
	IsFavorite    bool   `json:"is_favorite"`    // true-已点赞，false-未点赞
	Title         string `json:"title"`          // 视频标题
}
