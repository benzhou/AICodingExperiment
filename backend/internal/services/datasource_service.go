package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
)

// Error definitions
var (
	ErrDataSourceNotFound = errors.New("data source not found")
	ErrDataSourceExists   = errors.New("data source with this name already exists")
)

// DataSourceService provides methods for managing data sources
type DataSourceService struct {
	dataSourceRepo repository.DataSourceRepository
}

// NewDataSourceService creates a new data source service
func NewDataSourceService(dataSourceRepo repository.DataSourceRepository) *DataSourceService {
	return &DataSourceService{
		dataSourceRepo: dataSourceRepo,
	}
}

// CreateDataSource creates a new data source
func (s *DataSourceService) CreateDataSource(name, description string) (*models.DataSource, error) {
	dataSource := &models.DataSource{
		Name:        name,
		Description: description,
	}

	err := s.dataSourceRepo.CreateDataSource(dataSource)
	if err != nil {
		if err == repository.ErrDataSourceExists {
			return nil, ErrDataSourceExists
		}
		return nil, err
	}

	return dataSource, nil
}

// GetDataSourceByID retrieves a data source by ID
func (s *DataSourceService) GetDataSourceByID(id string) (*models.DataSource, error) {
	dataSource, err := s.dataSourceRepo.GetDataSourceByID(id)
	if err != nil {
		if err == repository.ErrDataSourceNotFound {
			return nil, ErrDataSourceNotFound
		}
		return nil, err
	}
	return dataSource, nil
}

// GetDataSourceByName retrieves a data source by name
func (s *DataSourceService) GetDataSourceByName(name string) (*models.DataSource, error) {
	dataSource, err := s.dataSourceRepo.GetDataSourceByName(name)
	if err != nil {
		if err == repository.ErrDataSourceNotFound {
			return nil, ErrDataSourceNotFound
		}
		return nil, err
	}
	return dataSource, nil
}

// UpdateDataSource updates a data source
func (s *DataSourceService) UpdateDataSource(id, name, description string) (*models.DataSource, error) {
	dataSource, err := s.dataSourceRepo.GetDataSourceByID(id)
	if err != nil {
		if err == repository.ErrDataSourceNotFound {
			return nil, ErrDataSourceNotFound
		}
		return nil, err
	}

	dataSource.Name = name
	dataSource.Description = description

	err = s.dataSourceRepo.UpdateDataSource(dataSource)
	if err != nil {
		if err == repository.ErrDataSourceExists {
			return nil, ErrDataSourceExists
		}
		return nil, err
	}

	return dataSource, nil
}

// DeleteDataSource deletes a data source
func (s *DataSourceService) DeleteDataSource(id string) error {
	err := s.dataSourceRepo.DeleteDataSource(id)
	if err != nil {
		if err == repository.ErrDataSourceNotFound {
			return ErrDataSourceNotFound
		}
		return err
	}
	return nil
}

// GetAllDataSources retrieves all data sources
func (s *DataSourceService) GetAllDataSources() ([]models.DataSource, error) {
	return s.dataSourceRepo.GetAllDataSources()
}

// SearchDataSources searches for data sources matching the query
func (s *DataSourceService) SearchDataSources(query string, limit, offset int) ([]models.DataSource, int, error) {
	return s.dataSourceRepo.SearchDataSources(query, limit, offset)
}
