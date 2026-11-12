package database

import "io/fs"

type Migrations struct {
	Owner string
	FS    fs.FS
}
