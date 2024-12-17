package handler

import (
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNew(t *testing.T) {

	mockDB := &gorm.DB{}
	mockUserStore := store.NewUserStore(mockDB)
	mockArticleStore := store.NewArticleStore(mockDB)

	tests := []struct {
		name         string
		logger       *zerolog.Logger
		userStore    *store.UserStore
		articleStore *store.ArticleStore
		wantNil      bool
		description  string
	}{
		{
			name:         "Successful Handler Creation",
			logger:       &zerolog.Logger{},
			userStore:    mockUserStore,
			articleStore: mockArticleStore,
			wantNil:      false,
			description:  "Should successfully create handler with valid parameters",
		},
		{
			name:         "Nil Logger",
			logger:       nil,
			userStore:    mockUserStore,
			articleStore: mockArticleStore,
			wantNil:      false,
			description:  "Should create handler with nil logger",
		},
		{
			name:         "Nil UserStore",
			logger:       &zerolog.Logger{},
			userStore:    nil,
			articleStore: mockArticleStore,
			wantNil:      false,
			description:  "Should create handler with nil UserStore",
		},
		{
			name:         "Nil ArticleStore",
			logger:       &zerolog.Logger{},
			userStore:    mockUserStore,
			articleStore: nil,
			wantNil:      false,
			description:  "Should create handler with nil ArticleStore",
		},
		{
			name:         "All Nil Parameters",
			logger:       nil,
			userStore:    nil,
			articleStore: nil,
			wantNil:      false,
			description:  "Should create handler with all nil parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Testing scenario:", tt.description)

			handler := New(tt.logger, tt.userStore, tt.articleStore)

			if tt.wantNil {
				assert.Nil(t, handler, "Handler should be nil")
			} else {
				assert.NotNil(t, handler, "Handler should not be nil")
				assert.Equal(t, tt.logger, handler.logger, "Logger not properly set")
				assert.Equal(t, tt.userStore, handler.us, "UserStore not properly set")
				assert.Equal(t, tt.articleStore, handler.as, "ArticleStore not properly set")
			}

			if !tt.wantNil {
				handler2 := New(tt.logger, tt.userStore, tt.articleStore)
				assert.NotSame(t, handler, handler2, "Handler instances should be different")
			}

			t.Log("Test completed successfully for:", tt.name)
		})
	}

	t.Run("Verify Handler Field Isolation", func(t *testing.T) {
		t.Log("Testing handler instance isolation")

		logger := &zerolog.Logger{}
		userStore := store.NewUserStore(mockDB)
		articleStore := store.NewArticleStore(mockDB)

		handler1 := New(logger, userStore, articleStore)
		handler2 := New(logger, userStore, articleStore)

		assert.NotSame(t, handler1, handler2, "Handlers should be different instances")
		assert.Equal(t, handler1.logger, handler2.logger, "Loggers should be the same")
		assert.Equal(t, handler1.us, handler2.us, "UserStores should be the same")
		assert.Equal(t, handler1.as, handler2.as, "ArticleStores should be the same")

		t.Log("Handler isolation test completed successfully")
	})
}
