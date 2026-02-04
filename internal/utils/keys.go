package utils

import "fmt"

type TableKey struct {
	Database string
	Schema   string
	Table    string
}

func GenTablePrefix(key *TableKey) string {
	return fmt.Sprintf("%s:%s:%s:", key.Database, key.Schema, key.Table)
}

func ParseTable(prefix string) (*TableKey, error) {
	var key TableKey
	_, err := fmt.Sscanf(prefix, "%s:%s:%s:", &key.Database, &key.Schema, &key.Table)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

type SequenceKey struct {
	Database string
	Schema   string
	Table    string
	Column   string
}

func GenSequenceKey(key *SequenceKey) string {
	return fmt.Sprintf("%s:%s:%s:%s:", key.Database, key.Schema, key.Table, key.Column)
}

func GenMetaPrefix(key string) string {
	return fmt.Sprintf("schema:%s", key)
}

func GenSequencePrefixKey(key *SequenceKey) string {
	return fmt.Sprintf("%s:%s:%s:", key.Database, key.Schema, key.Table)
}

func GenSequencePrefix(key string) string {
	return fmt.Sprintf("sequence:%s", key)
}

func ParseMetaPrefix(prefix string) (string, error) {
	var key string
	_, err := fmt.Sscanf(prefix, "schema:%s", &key)
	if err != nil {
		return "", err
	}
	return key, nil
}

type SchemaKey struct {
	Database string
	Schema   string
}

func GenSchemaPrefix(key *SchemaKey) string {
	return fmt.Sprintf("%s:%s:", key.Database, key.Schema)
}

func ParseSchema(prefix string) (*SchemaKey, error) {
	var key SchemaKey
	_, err := fmt.Sscanf(prefix, "%s:%s:", &key.Database, &key.Schema)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func GenTableRowKey(tablePrefix string, rowID string) string {
	tableRowPrefix := GenTableRowPrefix(tablePrefix)
	return fmt.Sprintf("%s%s:", tableRowPrefix, rowID)
}

func GenTableRowPrefix(tablePrefix string) string {
	return fmt.Sprintf("%srow:", tablePrefix)
}

type DatabaseKey struct {
	Database string
}

func GenDatabasePrefix(key *DatabaseKey) string {
	return fmt.Sprintf("%s:", key.Database)
}

func ParseDatabase(prefix string) (*DatabaseKey, error) {
	var key DatabaseKey
	_, err := fmt.Sscanf(prefix, "%s:", &key.Database)
	if err != nil {
		return nil, err
	}
	return &key, nil
}
