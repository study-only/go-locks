mocks:
	mockgen -package=mock -destination=mock/try_locker.go github.com/study-only/go-locks TryLocker
	mockgen -package=mock -destination=mock/lock_factory.go github.com/study-only/go-locks LockFactory
	mockgen -package=mock -destination=mock/expiry_lock_factory.go github.com/study-only/go-locks ExpiryLockFactory
