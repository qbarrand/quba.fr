package metadata

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

type ImageResult struct {
	Filename  string
	Name      string
	Date      string
	Location  string
	MainColor string
}

type DB struct {
	db *sql.DB
}

func OpenDB(path string, readonly bool) (*DB, error) {
	if readonly {
		path = fmt.Sprintf("file:%s?mode=ro", path)
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) AddImage(ctx context.Context, name, date, location, mainColor string) error {
	stmt := squirrel.
		Insert("image").
		Columns("name", "date", "location", "mainColor").
		Values(name, date, location, mainColor)

	_, err := stmt.RunWith(d.db).ExecContext(ctx)

	return err
}

func (d *DB) AddWebImage(ctx context.Context, filename, imageName string, width, height int, format string) error {
	var w, h *int

	if width != 0 {
		w = &width
	}

	if height != 0 {
		h = &height
	}

	stmt := squirrel.
		Insert("webImage").
		Columns("filename", "imageName", "width", "height", "format").
		Values(filename, imageName, w, h, format)

	_, err := stmt.RunWith(d.db).ExecContext(ctx)

	return err
}

func (d *DB) GetImage(ctx context.Context, name *string, width, height *int, format string) (*ImageResult, error) {
	and := squirrel.And{
		squirrel.Eq{"wi.format": format},
		squirrel.Eq{"wi.width": width},
		squirrel.Eq{"wi.height": height},
	}

	if name != nil {
		and = append(and, squirrel.Eq{"i.name": *name})
	} else {
		and = append(
			and,
			squirrel.Expr("i.name = (select name from image order by random() limit 1)"),
		)
	}

	query := squirrel.
		Select("wi.filename", "i.name", "i.date", "i.location", "i.mainColor").
		From("webImage wi").
		Join("image i on wi.imageName = i.name").
		Where(and).
		Limit(1)

	ir := ImageResult{}

	if err := query.RunWith(d.db).ScanContext(ctx, &ir.Filename, &ir.Name, &ir.Date, &ir.Location, &ir.MainColor); err != nil {
		return nil, fmt.Errorf("error while running the query: %v", err)
	}

	return &ir, nil
}

func (d *DB) Init(ctx context.Context) error {
	const sqlStmt = `
	PRAGMA foreign_keys = ON;

	create table image (
		name text primary key,
		date text,
		location text,
		mainColor text
	) without rowid;
	
	create table webImage (
		filename text primary key,
		format text,
		height integer,
		imageName text,
		width integer,
		foreign key(imageName) references image(name)
	) without rowid;
	`

	_, err := d.db.ExecContext(ctx, sqlStmt)

	return err
}
