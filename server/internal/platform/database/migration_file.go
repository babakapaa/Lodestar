package database

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"strings"
)

type migrationFile struct {
	version int
	name    string
	sql     string
}

func readMigrationFiles(source fs.FS) ([]migrationFile, error) {
	entries, err := fs.ReadDir(source, ".")
	if err != nil {
		return nil, fmt.Errorf("read migrations: %w", err)
	}
	files := make([]migrationFile, 0, len(entries))
	seen := make(map[int]string, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || path.Ext(entry.Name()) != ".sql" {
			continue
		}
		prefix, _, ok := strings.Cut(entry.Name(), "_")
		if !ok {
			return nil, fmt.Errorf("migration %s: name must be NNNN_description.sql", entry.Name())
		}
		version, err := strconv.Atoi(prefix)
		if err != nil || version <= 0 {
			return nil, fmt.Errorf("migration %s: bad version %q", entry.Name(), prefix)
		}
		if other, dup := seen[version]; dup {
			return nil, fmt.Errorf("migrations %s and %s share version %d", other, entry.Name(), version)
		}
		seen[version] = entry.Name()
		body, err := fs.ReadFile(source, entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}
		files = append(files, migrationFile{version: version, name: entry.Name(), sql: string(body)})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].version < files[j].version })
	return files, nil
}
