package migrations

import "embed"

// FS embeds all migration SQL files and atlas.sum for database initialization
//go:embed *.sql atlas.sum
var FS embed.FS
