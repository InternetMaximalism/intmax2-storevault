package third_party

import (
	"embed"
)

//go:embed OpenAPI/store_vault_service/*
var OpenAPIStoreVault embed.FS
