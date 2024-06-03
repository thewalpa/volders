package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	reps "github.com/thewalpa/volders/internal/representations"
)

// In Memory implementation of VoldersRepository for testing and developing
type TempRepository struct {
	folders map[reps.ObjectID]*reps.Folder
	files   map[reps.ObjectID]*reps.File
	mu      sync.RWMutex
}

func NewTempRepository() VolderRepository {
	return &TempRepository{
		folders: make(map[reps.ObjectID]*reps.Folder),
		files:   make(map[reps.ObjectID]*reps.File),
	}
}

func (repo *TempRepository) CreateFolder(folder *reps.Folder) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.folders[folder.ID]; exists {
		return errors.New("folder already exists")
	}

	folder.CreationDate = time.Now()
	folder.ModifiedDate = folder.CreationDate
	repo.folders[folder.ID] = folder
	return nil
}

// GetFolder retrieves a folder by ID
func (repo *TempRepository) GetFolder(id reps.ObjectID) (*reps.Folder, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	folder, exists := repo.folders[id]
	if !exists {
		return nil, errors.New("folder not found")
	}
	return folder, nil
}

// UpdateFolder updates an existing folder
func (repo *TempRepository) UpdateFolder(folder *reps.Folder) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.folders[folder.ID]; !exists {
		return errors.New("folder does not exist")
	}

	folder.ModifiedDate = time.Now()
	repo.folders[folder.ID] = folder
	return nil
}

// DeleteFolder deletes a folder by ID
func (repo *TempRepository) DeleteFolder(id reps.ObjectID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.folders[id]; !exists {
		return errors.New("folder not found")
	}
	delete(repo.folders, id)
	return nil
}

// CreateFile adds a new file to the repository
func (repo *TempRepository) CreateFile(file *reps.File) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.files[file.ID]; exists {
		return errors.New("file already exists")
	}

	file.CreationDate = time.Now()
	file.ModifiedDate = file.CreationDate
	repo.files[file.ID] = file
	return nil
}

// GetFile retrieves a file by ID
func (repo *TempRepository) GetFile(id reps.ObjectID) (*reps.File, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	file, exists := repo.files[id]
	if !exists {
		return nil, errors.New("file not found")
	}
	return file, nil
}

// UpdateFile updates an existing file
func (repo *TempRepository) UpdateFile(file *reps.File) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.files[file.ID]; !exists {
		return errors.New("file does not exist")
	}

	file.ModifiedDate = time.Now()
	repo.files[file.ID] = file
	return nil
}

// DeleteFile deletes a file by ID
func (repo *TempRepository) DeleteFile(id reps.ObjectID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.files[id]; !exists {
		return errors.New("file not found")
	}
	delete(repo.files, id)
	return nil
}

func (repo *TempRepository) GetFolderHierarchy(ctx context.Context, folderID reps.ObjectID) ([]reps.Folder, error) {
	return []reps.Folder{}, errors.New("not implemented")
}
