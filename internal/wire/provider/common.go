package provider

import (
	"gofile/internal/repository"

	"github.com/google/wire"
)

/*
- Trong trường hợp có 1 nhóm repository được sử dụng cho các service
ta có thể gom nhóm lại và inject vào từng service (thay vì phải inject từng repo)
*/
var CommonRepositoryProviderSet = wire.NewSet(
	repository.NewAccountRepository,
	// repository.NewUserRepository,
	// ... other common repositories or services
)
