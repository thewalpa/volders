package repository

import (
	"context"
	"database/sql"

	reps "github.com/thewalpa/volders/internal/representations"
)

// Interface for repository layer, provides methods for data access
type VolderRepository interface {
	GetFolder(id reps.ObjectID) (*reps.Folder, error)
	GetFile(id reps.ObjectID) (*reps.File, error)
	CreateFolder(folder *reps.Folder) error
	CreateFile(file *reps.File) error
	UpdateFolder(folder *reps.Folder) error
	UpdateFile(file *reps.File) error
	DeleteFolder(id reps.ObjectID) error
	DeleteFile(id reps.ObjectID) error
	GetFolderHierarchy(ctx context.Context, folderID reps.ObjectID) ([]reps.Folder, error)
}

// Implements VolderRepository for PostgreSQL Data Access
type PGVolderRepository struct {
	DB *sql.DB
}

func NewPGVolderRepository(db *sql.DB) VolderRepository {
	return &PGVolderRepository{DB: db}
}

func (repo *PGVolderRepository) GetFolder(id reps.ObjectID) (*reps.Folder, error) {
	folder := &reps.Folder{}
	query := `SELECT id, parent_id, name, creation_date, modified_date FROM folders WHERE id = $1`
	err := repo.DB.QueryRow(query, id).Scan(&folder.ID, &folder.ParentID, &folder.Name, &folder.CreationDate, &folder.ModifiedDate)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (repo *PGVolderRepository) GetFile(id reps.ObjectID) (*reps.File, error) {
	file := &reps.File{}
	query := `SELECT id, folder_id, name, content_type, size, creation_date, modified_date, data FROM files WHERE id = $1`
	err := repo.DB.QueryRow(query, id).Scan(&file.ID, &file.FolderID, &file.Name, &file.ContentType, &file.Size, &file.CreationDate, &file.ModifiedDate, &file.Data)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (repo *PGVolderRepository) CreateFolder(folder *reps.Folder) error {
	query := `INSERT INTO folders (parent_id, name, creation_date, modified_date) VALUES ($1, $2, $3, $4) RETURNING id`
	err := repo.DB.QueryRow(query, folder.ParentID, folder.Name, folder.CreationDate, folder.ModifiedDate).Scan(&folder.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PGVolderRepository) CreateFile(file *reps.File) error {
	query := `INSERT INTO files (folder_id, name, content_type, size, creation_date, modified_date, data) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := repo.DB.QueryRow(query, file.FolderID, file.Name, file.ContentType, file.Size, file.CreationDate, file.ModifiedDate, file.Data).Scan(&file.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PGVolderRepository) UpdateFolder(folder *reps.Folder) error {
	query := `UPDATE folders SET parent_id = $1, name = $2, modified_date = $3 WHERE id = $4`
	_, err := repo.DB.Exec(query, folder.ParentID, folder.Name, folder.ModifiedDate, folder.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PGVolderRepository) UpdateFile(file *reps.File) error {
	query := `UPDATE files SET folder_id = $1, name = $2, content_type = $3, size = $4, modified_date = $5, data = $6 WHERE id = $7`
	_, err := repo.DB.Exec(query, file.FolderID, file.Name, file.ContentType, file.Size, file.ModifiedDate, file.Data, file.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PGVolderRepository) DeleteFolder(id reps.ObjectID) error {
	query := `DELETE FROM folders WHERE id = $1`
	_, err := repo.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PGVolderRepository) DeleteFile(id reps.ObjectID) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := repo.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PGVolderRepository) GetFolderHierarchy(ctx context.Context, folderID reps.ObjectID) ([]reps.Folder, error) {
	var folders []reps.Folder
	query := `
WITH RECURSIVE folder_hierarchy AS (
    SELECT id, parent_id, name, creation_date, modified_date
    FROM folders
    WHERE id = $1
    UNION ALL
    SELECT f.id, f.parent_id, f.name, f.creation_date, f.modified_date
    FROM folders f
    INNER JOIN folder_hierarchy fh ON f.parent_id = fh.id
)
SELECT id, parent_id, name, creation_date, modified_classifieded_date
FROM folder_hierarchy;`

	rows, err := repo.DB.QueryContext(ctx, query, folderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var folder reps.Folder
		err := rows.Scan(&folder.ID, &folder.ParentID, &folder.Name, &folder.CreationDate, &folder.ModifiedDate)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return folders, nil
}
