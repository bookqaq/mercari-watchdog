package blacklist

import "fmt"

type BlockedSeller struct {
	UserID   int64
	UserName string
	Reason   string
}

func (s *BlockedSeller) ParseID() string {
	return fmt.Sprintf("%s%d", "https://jp.mercari.com/user/profile/", s.UserID)
}

func (s *BlockedSeller) FormatSimplifiedChinese() string {
	return fmt.Sprintf("卖家:%s(%d)\n原因:%s", s.UserName, s.UserID, s.Reason)
}
