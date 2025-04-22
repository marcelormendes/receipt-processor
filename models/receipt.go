package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Date represents a date in YYYY-MM-DD format
type Date string

// Validate checks if the date is in the correct format
func (d Date) Validate() error {
	if _, err := time.Parse("2006-01-02", string(d)); err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}
	return nil
}

// Time represents a time in HH:MM format
type Time string

// Validate checks if the time is in the correct format
func (t Time) Validate() error {
	if _, err := time.Parse("15:04", string(t)); err != nil {
		return fmt.Errorf("invalid time format: %w", err)
	}
	return nil
}

// Price represents a price value that can be parsed from a string
type Price float64

// UnmarshalJSON custom unmarshaler for Price
func (p *Price) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("invalid price format: %w", err)
	}
	*p = Price(val)
	return nil
}

// MarshalJSON custom marshaler for Price
func (p Price) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%.2f", float64(p)))
}

// Receipt represents a receipt submitted for processing
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate Date   `json:"purchaseDate"` // YYYY-MM-DD
	PurchaseTime Time   `json:"purchaseTime"` // HH:MM
	Items        []Item `json:"items"`
	Total        Price  `json:"total"` // Price type handles parsing
}

// Validate performs validation on the receipt and all its fields
func (r *Receipt) Validate() error {
	// Check required fields
	if strings.TrimSpace(r.Retailer) == "" {
		return fmt.Errorf("retailer is required")
	}

	// Validate purchase date
	if err := r.PurchaseDate.Validate(); err != nil {
		return err
	}

	// Validate purchase time
	if err := r.PurchaseTime.Validate(); err != nil {
		return err
	}

	// Require at least one item
	if len(r.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}

	// Validate each item
	for i, item := range r.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item %d: %w", i+1, err)
		}
	}

	return nil
}

// Item represents an individual item on a receipt
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            Price  `json:"price"` // Price type handles parsing
}

// Validate checks if the item is valid
func (i *Item) Validate() error {
	if strings.TrimSpace(i.ShortDescription) == "" {
		return fmt.Errorf("short description is required")
	}

	return nil
}

// ReceiptResponse is returned when processing a receipt
type ReceiptResponse struct {
	ID string `json:"id"`
}

// PointsResponse is returned when querying points for a receipt
type PointsResponse struct {
	Points int `json:"points"`
}
