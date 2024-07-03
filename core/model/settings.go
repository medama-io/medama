package model

type GetSettingsUsage struct {
	Threads     int    `db:"threads"`
	MemoryLimit string `db:"memory_limit"`
}
