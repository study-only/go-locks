package di

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/study-only/go-locks"
	"github.com/study-only/go-locks/mock"
)

type usecase struct {
	lockFactory golocks.ExpiryLockFactory
}

func (u *usecase) CheckFrequntSubmit(key string, expiry time.Duration) (ok bool) {
	lock := u.lockFactory.NewLock(key, expiry)
	if err := lock.TryLock(); err != nil {
		return false
	}

	// let the lock expire, and unlock automatically
	return true
}

func TestUsecase_CheckFrequentSubmit(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	key := "book:1:submit"
	expiry := time.Second

	lock := mock.NewMockTryLocker(ctl)
	lock.EXPECT().TryLock().Return(nil)
	lockFactory := mock.NewMockExpiryLockFactory(ctl)
	lockFactory.EXPECT().NewLock(key, expiry).Return(lock)

	bookUsecase := usecase{lockFactory}
	ok := bookUsecase.CheckFrequntSubmit(key, expiry)
	assert.True(t, ok)
}
