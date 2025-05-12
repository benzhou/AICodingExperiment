package services

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

// Message types for the queue
const (
	MessageTypeProcessDataSource = "process_data_source"
	MessageTypeRunMatchSet       = "run_match_set"
)

// QueueMessage represents a message in the queue
type QueueMessage struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
}

// ProcessDataSourcePayload contains data for processing a data source
type ProcessDataSourcePayload struct {
	DataSourceID string `json:"data_source_id"`
	SchemaID     string `json:"schema_id"`
	UploadID     string `json:"upload_id"`
	TenantID     string `json:"tenant_id"`
}

// RunMatchSetPayload contains data for running a match set
type RunMatchSetPayload struct {
	MatchSetID string `json:"match_set_id"`
	TenantID   string `json:"tenant_id"`
}

// QueueService provides methods for handling messages from a queue
type QueueService struct {
	schemaService      *SchemaService
	dataSourceService  *DataSourceService
	matchSetService    *MatchSetService
	transactionService *TransactionService
	uploadService      *UploadService
	// In a real implementation, we would have an SQS client here
	// sqsClient       *sqs.SQS
}

// NewQueueService creates a new queue service
func NewQueueService(
	schemaService *SchemaService,
	dataSourceService *DataSourceService,
	matchSetService *MatchSetService,
	transactionService *TransactionService,
	uploadService *UploadService,
) *QueueService {
	return &QueueService{
		schemaService:      schemaService,
		dataSourceService:  dataSourceService,
		matchSetService:    matchSetService,
		transactionService: transactionService,
		uploadService:      uploadService,
	}
}

// SendProcessDataSourceMessage sends a message to process a data source
func (s *QueueService) SendProcessDataSourceMessage(dataSourceID, schemaID, uploadID, tenantID string) error {
	payload := ProcessDataSourcePayload{
		DataSourceID: dataSourceID,
		SchemaID:     schemaID,
		UploadID:     uploadID,
		TenantID:     tenantID,
	}

	return s.sendMessage(MessageTypeProcessDataSource, payload)
}

// SendRunMatchSetMessage sends a message to run a match set
func (s *QueueService) SendRunMatchSetMessage(matchSetID, tenantID string) error {
	payload := RunMatchSetPayload{
		MatchSetID: matchSetID,
		TenantID:   tenantID,
	}

	return s.sendMessage(MessageTypeRunMatchSet, payload)
}

// sendMessage marshals and sends a message to the queue
func (s *QueueService) sendMessage(messageType string, payload interface{}) error {
	// Marshal the payload
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create the message
	message := QueueMessage{
		Type:      messageType,
		Payload:   payloadJSON,
		Timestamp: time.Now(),
	}

	// Marshal the full message
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// In a real implementation, we would send this to SQS
	// For now, we'll just log it and handle it directly for development
	log.Printf("Sending message to queue: %s", string(messageJSON))

	// Directly process the message for development
	return s.HandleMessage(message)
}

// StartListener starts listening for messages from the queue
func (s *QueueService) StartListener() {
	// In a real implementation, this would start a long-running goroutine
	// that pulls messages from SQS and processes them
	log.Println("Queue listener started")
}

// HandleMessage processes a single message from the queue
func (s *QueueService) HandleMessage(message QueueMessage) error {
	log.Printf("Processing message of type: %s", message.Type)

	switch message.Type {
	case MessageTypeProcessDataSource:
		var payload ProcessDataSourcePayload
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return err
		}
		return s.handleProcessDataSource(payload)

	case MessageTypeRunMatchSet:
		var payload RunMatchSetPayload
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return err
		}
		return s.handleRunMatchSet(payload)

	default:
		return errors.New("unknown message type")
	}
}

// handleProcessDataSource processes a data source upload
func (s *QueueService) handleProcessDataSource(payload ProcessDataSourcePayload) error {
	log.Printf("Processing data source %s for tenant %s", payload.DataSourceID, payload.TenantID)

	// In a real implementation, this would:
	// 1. Get the upload file from storage
	// 2. Parse it according to the schema
	// 3. Insert transactions into the database
	// 4. Update the upload status

	// For development, we'll just log it
	log.Printf("Data source %s processed successfully", payload.DataSourceID)

	return nil
}

// handleRunMatchSet runs a match set
func (s *QueueService) handleRunMatchSet(payload RunMatchSetPayload) error {
	log.Printf("Running match set %s for tenant %s", payload.MatchSetID, payload.TenantID)

	// In a real implementation, this would:
	// 1. Get the match set and its data sources
	// 2. Get the matching rules
	// 3. Apply the rules to find matches
	// 4. Create matched_transaction and unmatched_transaction records
	// 5. Update the match progress

	// For development, we'll just log it
	log.Printf("Match set %s completed successfully", payload.MatchSetID)

	return nil
}
