import api from './api';

// Define the SchemaField interface for field definition
export interface SchemaField {
  name: string;
  displayName: string;
  type: string;  // string, number, date, boolean
  required: boolean;
  format?: string;
  description?: string;
}

// Define the SchemaDefinition interface
export interface SchemaDefinition {
  fields: SchemaField[];
  dateFormat: string;
  defaultMappings: Record<string, string>;
  requiredFields: string[];
}

// Define the DataSource interface matching backend model
export interface DataSource {
  id: string;
  name: string;
  description: string;
  schemaDefinition?: SchemaDefinition;
  created_at: number; // Unix timestamp in milliseconds
  updated_at: number; // Unix timestamp in milliseconds
}

// Define ImportRecord interface
export interface ImportRecord {
  id: string;
  dataSourceId: string;
  fileName: string;
  fileSize: number;
  status: 'Processing' | 'Completed' | 'Failed';
  rowCount: number;
  successCount: number;
  errorCount: number;
  importedBy: string;
  createdAt: string;
  updatedAt: string;
  metadata?: any;
}

// Define RawTransaction interface
export interface RawTransaction {
  id: string;
  importId: string;
  dataSourceId: string;
  rowNumber: number;
  data: any;
  errorMessage?: string;
  createdAt: string;
}

// Define the request interfaces
export interface CreateDataSourceRequest {
  name: string;
  description: string;
  schemaDefinition?: SchemaDefinition;
}

export interface UpdateDataSourceRequest {
  name: string;
  description: string;
  schemaDefinition?: SchemaDefinition;
}

// Define pagination response interface
export interface PaginationResponse {
  total: number;
  limit: number;
  offset: number;
  hasMore: boolean;
}

// Define paginated response interface
export interface PaginatedResponse<T> {
  data: T[];
  pagination: PaginationResponse;
}

// API endpoints for data sources
const ENDPOINTS = {
  BASE: '/api/v1/datasources',
  SEARCH: '/api/v1/datasources/search',
  byId: (id: string) => `/api/v1/datasources/${id}`,
  imports: (dataSourceId: string) => `/api/v1/datasources/${dataSourceId}/imports`,
  importById: (id: string) => `/api/v1/imports/${id}`,
  rawTransactionsByImport: (importId: string) => `/api/v1/imports/${importId}/raw-transactions`,
  rawTransactionById: (id: string) => `/api/v1/raw-transactions/${id}`,
};

// DataSource service using the centralized API client
export const dataSourceService = {
  // Get all data sources
  async getDataSources(): Promise<DataSource[]> {
    try {
      const response = await api.get(ENDPOINTS.BASE);
      console.log('Retrieved data sources:', response.data);
      return response.data || [];
    } catch (error) {
      console.error('Error fetching data sources:', error);
      return [];
    }
  },

  // Search data sources with server-side pagination and filtering
  async searchDataSources(
    query: string,
    limit: number = 10,
    offset: number = 0
  ): Promise<PaginatedResponse<DataSource>> {
    try {
      const response = await api.get(
        `${ENDPOINTS.SEARCH}?q=${encodeURIComponent(query)}&limit=${limit}&offset=${offset}`
      );
      return response.data;
    } catch (error) {
      console.error('Error searching data sources:', error);
      throw error;
    }
  },

  // Get data source by ID
  async getDataSourceById(id: string): Promise<DataSource> {
    try {
      const response = await api.get(ENDPOINTS.byId(id));
      return response.data;
    } catch (error) {
      console.error(`Error fetching data source with ID ${id}:`, error);
      throw error;
    }
  },

  // Create a new data source
  async createDataSource(data: CreateDataSourceRequest): Promise<DataSource> {
    try {
      const response = await api.post(ENDPOINTS.BASE, data);
      return response.data;
    } catch (error) {
      console.error('Error creating data source:', error);
      throw error;
    }
  },

  // Update a data source
  async updateDataSource(id: string, data: UpdateDataSourceRequest): Promise<DataSource> {
    try {
      const response = await api.put(ENDPOINTS.byId(id), data);
      return response.data;
    } catch (error) {
      console.error(`Error updating data source with ID ${id}:`, error);
      throw error;
    }
  },

  // Delete a data source
  async deleteDataSource(id: string): Promise<void> {
    try {
      await api.delete(ENDPOINTS.byId(id));
    } catch (error) {
      console.error(`Error deleting data source with ID ${id}:`, error);
      throw error;
    }
  },

  // Get imports for a data source with pagination
  async getImportsByDataSource(
    dataSourceId: string, 
    limit: number = 10, 
    offset: number = 0
  ): Promise<PaginatedResponse<ImportRecord>> {
    try {
      const response = await api.get(
        `${ENDPOINTS.imports(dataSourceId)}?limit=${limit}&offset=${offset}`
      );
      return response.data;
    } catch (error) {
      console.error(`Error fetching imports for data source ${dataSourceId}:`, error);
      throw error;
    }
  },

  // Get import by ID
  async getImportById(id: string): Promise<ImportRecord> {
    try {
      const response = await api.get(ENDPOINTS.importById(id));
      return response.data;
    } catch (error) {
      console.error(`Error fetching import with ID ${id}:`, error);
      throw error;
    }
  },

  // Delete an import
  async deleteImport(id: string): Promise<void> {
    try {
      await api.delete(ENDPOINTS.importById(id));
    } catch (error) {
      console.error(`Error deleting import with ID ${id}:`, error);
      throw error;
    }
  },

  // Get raw transactions for an import with pagination
  async getRawTransactionsByImport(
    importId: string, 
    limit: number = 10, 
    offset: number = 0
  ): Promise<PaginatedResponse<RawTransaction>> {
    try {
      const response = await api.get(
        `${ENDPOINTS.rawTransactionsByImport(importId)}?limit=${limit}&offset=${offset}`
      );
      return response.data;
    } catch (error) {
      console.error(`Error fetching raw transactions for import ${importId}:`, error);
      throw error;
    }
  },

  // Get raw transaction by ID
  async getRawTransactionById(id: string): Promise<RawTransaction> {
    try {
      const response = await api.get(ENDPOINTS.rawTransactionById(id));
      return response.data;
    } catch (error) {
      console.error(`Error fetching raw transaction with ID ${id}:`, error);
      throw error;
    }
  }
}; 