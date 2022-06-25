package blacklist

import "fmt"

type BlockedSeller struct {
	UserID   int64  `json:"userID" bson:"userID"`
	UserName string `json:"userName" bson:"userName"`
	Reason   string `json:"reason" bson:"reason"` // reason text, there would be other flags to differ reasons, but not now.
}

// get http url to seller's homepage
func (s *BlockedSeller) ParseID() string {
	return fmt.Sprintf("%s%d", "https://jp.mercari.com/user/profile/", s.UserID)
}

// generate qq message
func (s *BlockedSeller) FormatSimplifiedChinese() string {
	return fmt.Sprintf("卖家:%s(%d)\n原因:%s", s.UserName, s.UserID, s.Reason)
}
