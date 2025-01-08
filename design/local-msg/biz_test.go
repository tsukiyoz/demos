package localmsg_test

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type BizTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *BizTestSuite) Test1() {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// do something
		return nil
	})
	_ = err
}
