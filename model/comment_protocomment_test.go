package model

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

const ISO8601 = "2006-01-02T15:04:05.999Z"func TestCommentProtoComment(t *testing.T) {
	tests := []struct {
		name     string
		comment  Comment
		expected struct {
			id        string
			body      string
			createdAt string
			updatedAt string
		}
	}{
		{
			name: "Scenario 1: Valid Comment Conversion",
			comment: Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				},
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			expected: struct {
				id        string
				body      string
				createdAt string
				updatedAt string
			}{
				id:        "1",
				body:      "Test comment",
				createdAt: "2023-01-01T12:00:00.000Z",
				updatedAt: "2023-01-01T12:00:00.000Z",
			},
		},
		{
			name: "Scenario 2: Zero Values",
			comment: Comment{
				Model: gorm.Model{
					ID:        0,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				Body: "",
			},
			expected: struct {
				id        string
				body      string
				createdAt string
				updatedAt string
			}{
				id:        "0",
				body:      "",
				createdAt: "0001-01-01T00:00:00.000Z",
				updatedAt: "0001-01-01T00:00:00.000Z",
			},
		},
		{
			name: "Scenario 3: Special Characters",
			comment: Comment{
				Model: gorm.Model{
					ID:        123,
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				},
				Body: "Test 😊 特殊文字 !@#$%^&*()",
			},
			expected: struct {
				id        string
				body      string
				createdAt string
				updatedAt string
			}{
				id:        "123",
				body:      "Test 😊 特殊文字 !@#$%^&*()",
				createdAt: "2023-01-01T12:00:00.000Z",
				updatedAt: "2023-01-01T12:00:00.000Z",
			},
		},
		{
			name: "Scenario 4: Maximum Values",
			comment: Comment{
				Model: gorm.Model{
					ID:        ^uint(0),
					CreatedAt: time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
					UpdatedAt: time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
				},
				Body: string(make([]byte, 1000)),
			},
			expected: struct {
				id        string
				body      string
				createdAt string
				updatedAt string
			}{
				id:        "18446744073709551615",
				body:      string(make([]byte, 1000)),
				createdAt: "9999-12-31T23:59:59.999Z",
				updatedAt: "9999-12-31T23:59:59.999Z",
			},
		},
		{
			name: "Scenario 5: Different Timestamps",
			comment: Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				},
				Body: "Updated comment",
			},
			expected: struct {
				id        string
				body      string
				createdAt string
				updatedAt string
			}{
				id:        "1",
				body:      "Updated comment",
				createdAt: "2023-01-01T12:00:00.000Z",
				updatedAt: "2023-01-02T12:00:00.000Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Testing:", tt.name)

			result := tt.comment.ProtoComment()

			assert.Equal(t, tt.expected.id, result.Id, "ID mismatch")
			assert.Equal(t, tt.expected.body, result.Body, "Body mismatch")
			assert.Equal(t, tt.expected.createdAt, result.CreatedAt, "CreatedAt mismatch")
			assert.Equal(t, tt.expected.updatedAt, result.UpdatedAt, "UpdatedAt mismatch")

			t.Log("Test completed successfully")
		})
	}
}
