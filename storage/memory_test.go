package storage

import (
	"context"
	"sync"
	"testing"
)

func TestSaveReceipt(t *testing.T) {
	store := NewMemoryStorage()
	ctx := context.Background()
	
	id, err := store.SaveReceipt(ctx, 100)
	if err != nil {
		t.Fatalf("Failed to save receipt: %v", err)
	}
	
	if id == "" {
		t.Error("Expected non-empty ID")
	}
	
	points, err := store.GetPoints(ctx, id)
	if err != nil {
		t.Errorf("Failed to retrieve points: %v", err)
	}
	
	if points != 100 {
		t.Errorf("Expected 100 points, got %d", points)
	}
}

func TestGetPoints(t *testing.T) {
	store := NewMemoryStorage()
	ctx := context.Background()
	
	// Test retrieving non-existent ID
	_, err := store.GetPoints(ctx, "non-existent-id")
	if err == nil {
		t.Error("Expected error for non-existent ID, got nil")
	}
	
	// Test retrieving valid ID
	id, err := store.SaveReceipt(ctx, 50)
	if err != nil {
		t.Fatalf("Failed to save receipt: %v", err)
	}
	
	points, err := store.GetPoints(ctx, id)
	if err != nil {
		t.Errorf("Failed to retrieve points: %v", err)
	}
	
	if points != 50 {
		t.Errorf("Expected 50 points, got %d", points)
	}
}

func TestConcurrentAccess(t *testing.T) {
	store := NewMemoryStorage()
	count := 100
	var wg sync.WaitGroup
	
	// Create a background context
	ctx := context.Background()
	
	// Save receipts concurrently
	wg.Add(count)
	ids := make([]string, count)
	errChan := make(chan error, count)
	
	for i := 0; i < count; i++ {
		go func(i int) {
			defer wg.Done()
			id, err := store.SaveReceipt(ctx, i)
			if err != nil {
				errChan <- err
				return
			}
			ids[i] = id
		}(i)
	}
	
	wg.Wait()
	close(errChan)
	
	// Check for any errors during saving
	for err := range errChan {
		t.Fatalf("Error during concurrent save: %v", err)
	}
	
	// Verify all saved receipts
	for i, id := range ids {
		if id == "" {
			continue // Skip if ID wasn't set properly
		}
		
		points, err := store.GetPoints(ctx, id)
		if err != nil {
			t.Errorf("Failed to retrieve points for ID %s: %v", id, err)
		}
		
		if points != i {
			t.Errorf("Expected %d points for ID %s, got %d", i, id, points)
		}
	}
}

func TestContextCancellation(t *testing.T) {
	store := NewMemoryStorage()
	
	// Create canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	// Test saving with canceled context
	_, err := store.SaveReceipt(ctx, 100)
	if err == nil {
		t.Error("Expected error when saving with canceled context, got nil")
	}
	
	// Test getting with canceled context
	_, err = store.GetPoints(ctx, "any-id")
	if err == nil {
		t.Error("Expected error when getting with canceled context, got nil")
	}
	
	// Test count with canceled context
	_, err = store.Count(ctx)
	if err == nil {
		t.Error("Expected error when counting with canceled context, got nil")
	}
}

func TestNewMemoryStorage(t *testing.T) {
	store := NewMemoryStorage()
	
	if store == nil {
		t.Error("Expected store to be non-nil")
		return // Avoid nil pointer dereference
	}
	
	if store.points == nil {
		t.Error("Expected points map to be initialized")
	}
}