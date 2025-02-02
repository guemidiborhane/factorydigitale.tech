package storage

func Setup() {
	SetupPostgres()
	SetupRedis()
	SetupSession()
	SetupMeilisearch()
}
