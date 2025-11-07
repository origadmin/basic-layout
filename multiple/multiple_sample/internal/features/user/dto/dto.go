/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package dto

// UserProfile represents a DTO for a user's profile, potentially aggregating data
// from multiple domain objects or adapting for a specific view.
type UserProfile struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"` // Example of an optional field
	// Add other fields as needed for specific internal or view-layer data transfer
}