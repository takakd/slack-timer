package db

type NoSqlDb interface {
	ReadString(key string) string
	WriteString(key, value string)
}
