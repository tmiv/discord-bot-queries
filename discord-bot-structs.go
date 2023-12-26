package main

type VerifyMembershipResponse struct {
	Member    bool   `json:"member"`
	UserId    string `json:"user_id"`
	ChannelId string `json:"channel_id"`
}
