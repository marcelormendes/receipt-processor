package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	rperrors "github.com/marcelorm/receipt-processor/errors"
)

// ReceiptStorage defines the interface for storing and retrieving receipt points
type ReceiptStorage interface {
	// SaveReceipt saves the points for a receipt and returns the generated ID
	SaveReceipt(ctx context.Context, points int) (string, error)

	// GetPoints retrieves the points for a receipt by ID
	GetPoints(ctx context.Context, id string) (int, error)

	// Count returns the number of receipts in the store (for testing)
	Count(ctx context.Context) (int, error)
}

// MemoryStore provides thread-safe in-memory storage for receipt points
type MemoryStore struct {
	points map[string]int
	mutex  sync.RWMutex
}

// Verify MemoryStore implements ReceiptStorage interface
var _ ReceiptStorage = (*MemoryStore)(nil)

// NewMemoryStorage creates a new in-memory receipt store
func NewMemoryStorage() *MemoryStore {
	return &MemoryStore{
		points: make(map[string]int),
	}
}

// SaveReceipt saves the points for a receipt and returns the generated ID
func (s *MemoryStore) SaveReceipt(ctx context.Context, points int) (string, error) {
	if ctx.Err() != nil {
		return "", rperrors.Wrap(rperrors.ErrContextCancelled, ctx.Err(), "context canceled before saving receipt")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := uuid.New().String()
	s.points[id] = points
	return id, nil
}

// GetPoints retrieves the points for a receipt by ID
func (s *MemoryStore) GetPoints(ctx context.Context, id string) (int, error) {
	if ctx.Err() != nil {
		return 0, rperrors.Wrap(rperrors.ErrContextCancelled, ctx.Err(), "context canceled before retrieving points")
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	points, exists := s.points[id]
	if !exists {
		return 0, rperrors.New(rperrors.ErrReceiptNotFound, fmt.Sprintf("receipt with ID %s not found", id))
	}
	return points, nil
}

// Count returns the number of receipts in the store (for testing)
func (s *MemoryStore) Count(ctx context.Context) (int, error) {
	// Create a context with a deadline for this operation (500ms)
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	// Use a channel to handle the operation with timeout
	type result struct {
		count int
		err   error
	}

	ch := make(chan result, 1)

	go func() {
		s.mutex.RLock()
		defer s.mutex.RUnlock()

		ch <- result{count: len(s.points)}
	}()

	// Wait for either the operation to complete or the context to timeout
	select {
	case <-ctx.Done():
		return 0, rperrors.Wrap(rperrors.ErrContextCancelled, ctx.Err(), "timeout or cancellation while counting receipts")
	case res := <-ch:
		return res.count, res.err
	}
}
