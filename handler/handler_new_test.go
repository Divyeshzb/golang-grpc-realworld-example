package handler

import (
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	type args struct {
		logger       *zerolog.Logger
		userStore    *store.UserStore
		articleStore *store.ArticleStore
	}

	logger := zerolog.New(nil)

	userStore := &store.UserStore{}
	articleStore := &store.ArticleStore{}

	tests := []struct {
		name     string
		args     args
		wantNil  bool
		validate func(*testing.T, *Handler)
	}{
		{
			name: "Scenario 1: Successfully Create Handler with Valid Parameters",
			args: args{
				logger:       &logger,
				userStore:    userStore,
				articleStore: articleStore,
			},
			wantNil: false,
			validate: func(t *testing.T, h *Handler) {
				assert.NotNil(t, h.logger, "logger should not be nil")
				assert.NotNil(t, h.us, "user store should not be nil")
				assert.NotNil(t, h.as, "article store should not be nil")
				assert.Equal(t, &logger, h.logger, "logger should match input")
				assert.Equal(t, userStore, h.us, "user store should match input")
				assert.Equal(t, articleStore, h.as, "article store should match input")
			},
		},
		{
			name: "Scenario 2: Create Handler with Nil Logger",
			args: args{
				logger:       nil,
				userStore:    userStore,
				articleStore: articleStore,
			},
			wantNil: false,
			validate: func(t *testing.T, h *Handler) {
				assert.Nil(t, h.logger, "logger should be nil")
				assert.NotNil(t, h.us, "user store should not be nil")
				assert.NotNil(t, h.as, "article store should not be nil")
			},
		},
		{
			name: "Scenario 3: Create Handler with Nil UserStore",
			args: args{
				logger:       &logger,
				userStore:    nil,
				articleStore: articleStore,
			},
			wantNil: false,
			validate: func(t *testing.T, h *Handler) {
				assert.NotNil(t, h.logger, "logger should not be nil")
				assert.Nil(t, h.us, "user store should be nil")
				assert.NotNil(t, h.as, "article store should not be nil")
			},
		},
		{
			name: "Scenario 4: Create Handler with Nil ArticleStore",
			args: args{
				logger:       &logger,
				userStore:    userStore,
				articleStore: nil,
			},
			wantNil: false,
			validate: func(t *testing.T, h *Handler) {
				assert.NotNil(t, h.logger, "logger should not be nil")
				assert.NotNil(t, h.us, "user store should not be nil")
				assert.Nil(t, h.as, "article store should be nil")
			},
		},
		{
			name: "Scenario 5: Create Handler with All Nil Parameters",
			args: args{
				logger:       nil,
				userStore:    nil,
				articleStore: nil,
			},
			wantNil: false,
			validate: func(t *testing.T, h *Handler) {
				assert.Nil(t, h.logger, "logger should be nil")
				assert.Nil(t, h.us, "user store should be nil")
				assert.Nil(t, h.as, "article store should be nil")
			},
		},
		{
			name: "Scenario 6: Verify Handler Field Isolation",
			args: args{
				logger:       &logger,
				userStore:    userStore,
				articleStore: articleStore,
			},
			wantNil: false,
			validate: func(t *testing.T, h *Handler) {

				logger2 := zerolog.New(nil)
				userStore2 := &store.UserStore{}
				articleStore2 := &store.ArticleStore{}

				h2 := New(&logger2, userStore2, articleStore2)

				assert.NotEqual(t, h.logger, h2.logger, "handlers should have different loggers")
				assert.NotEqual(t, h.us, h2.us, "handlers should have different user stores")
				assert.NotEqual(t, h.as, h2.as, "handlers should have different article stores")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Starting test:", tt.name)

			got := New(tt.args.logger, tt.args.userStore, tt.args.articleStore)

			if tt.wantNil {
				assert.Nil(t, got, "handler should be nil")
			} else {
				assert.NotNil(t, got, "handler should not be nil")
				tt.validate(t, got)
			}

			t.Log("Completed test:", tt.name)
		})
	}
}
