package assistant_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/andrian0vv/chatgpt-cli/internal/dto"
	"github.com/andrian0vv/chatgpt-cli/internal/logger"
	"github.com/andrian0vv/chatgpt-cli/internal/services/assistant"
	"github.com/andrian0vv/chatgpt-cli/internal/services/assistant/mocks"
)

func TestAssistant_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	l := logger.New(nil, logger.WithEnabled(false))

	client := mocks.NewMockclient(ctrl)
	client.EXPECT().Model().Return("test-model").AnyTimes()

	client.EXPECT().
		ModelExists(ctx).
		Return(true, nil)
	_, err := assistant.New(ctx, client, l)
	assert.NoError(t, err)

	client.EXPECT().
		ModelExists(ctx).
		Return(false, nil)
	_, err = assistant.New(ctx, client, l)
	assert.Error(t, err)

	client.EXPECT().
		ModelExists(ctx).
		Return(false, errors.New("some error"))
	_, err = assistant.New(ctx, client, l)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "some error")
}

func TestAssistant_Model(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	l := logger.New(nil, logger.WithEnabled(false))

	client := newMockClient(ctrl)

	a, err := assistant.New(ctx, client, l)
	assert.NoError(t, err)

	expected := "test-model"
	client.EXPECT().
		Model().
		Return(expected)

	model := a.Model()
	assert.Equal(t, expected, model)
}

func TestAssistant_GetModels(t *testing.T) {
	testCases := []struct {
		name      string
		clientFn  func(*gomock.Controller) *mocks.Mockclient
		mockSetup func()
		expected  []string
		wantErr   bool
	}{
		{
			name: "ok",
			clientFn: func(ctrl *gomock.Controller) *mocks.Mockclient {
				c := newMockClient(ctrl)
				c.EXPECT().
					GetModels(gomock.Any()).
					Return([]string{"model1", "model2"}, nil)
				return c
			},
			expected: []string{"model1", "model2"},
		},
		{
			name: "error",
			clientFn: func(ctrl *gomock.Controller) *mocks.Mockclient {
				c := newMockClient(ctrl)
				c.EXPECT().
					GetModels(gomock.Any()).
					Return(nil, errors.New("some error"))
				return c
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			l := logger.New(nil, logger.WithEnabled(false))

			a, err := assistant.New(ctx, tc.clientFn(ctrl), l)
			assert.NoError(t, err)

			models, err := a.GetModels(ctx)
			assert.Equal(t, tc.expected, models)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAssistant_SendChatMessage(t *testing.T) {
	testCases := []struct {
		name     string
		clientFn func(*gomock.Controller) *mocks.Mockclient
		inChat   *dto.Chat
		outChat  *dto.Chat
		in       string
		out      string
		wantErr  bool
	}{
		{
			name: "ok with a new chat",
			clientFn: func(ctrl *gomock.Controller) *mocks.Mockclient {
				c := newMockClient(ctrl)
				c.EXPECT().
					CreateChatCompletion(gomock.Any(), gomock.Any()).
					Return("Hi there!", nil)
				return c
			},
			in:     "Hello",
			out:    "Hi there!",
			inChat: dto.NewChat(),
			outChat: &dto.Chat{
				Messages: []dto.Message{
					{Role: dto.RoleUser, Content: "Hello"},
					{Role: dto.RoleAssistant, Content: "Hi there!"},
				},
			},
		},
		{
			name: "ok with the existing chat",
			clientFn: func(ctrl *gomock.Controller) *mocks.Mockclient {
				c := newMockClient(ctrl)
				c.EXPECT().
					CreateChatCompletion(gomock.Any(), gomock.Any()).
					Return("I can't provide real-time weather updates or current conditions", nil)
				return c
			},
			in:  "What is the weather in Lisbon?",
			out: "I can't provide real-time weather updates or current conditions",
			inChat: &dto.Chat{
				Messages: []dto.Message{
					{Role: dto.RoleUser, Content: "Hello"},
					{Role: dto.RoleAssistant, Content: "Hi there!"},
				},
			},
			outChat: &dto.Chat{
				Messages: []dto.Message{
					{Role: dto.RoleUser, Content: "Hello"},
					{Role: dto.RoleAssistant, Content: "Hi there!"},
					{Role: dto.RoleUser, Content: "What is the weather in Lisbon?"},
					{Role: dto.RoleAssistant, Content: "I can't provide real-time weather updates or current conditions"},
				},
			},
		},
		{
			name: "error",
			clientFn: func(ctrl *gomock.Controller) *mocks.Mockclient {
				c := newMockClient(ctrl)
				c.EXPECT().
					CreateChatCompletion(gomock.Any(), gomock.Any()).
					Return("", errors.New("API error"))
				return c
			},
			inChat:  dto.NewChat(),
			outChat: dto.NewChat(),
			in:      "Hello",
			wantErr: true,
		},
		{
			name: "empty chat",
			clientFn: func(ctrl *gomock.Controller) *mocks.Mockclient {
				c := newMockClient(ctrl)
				c.EXPECT().
					CreateChatCompletion(gomock.Any(), gomock.Any()).
					Return("Hi there!", nil)
				return c
			},
			in:  "How are you?",
			out: "Hi there!",
		},
		{
			name: "empty in",
			clientFn: func(ctrl *gomock.Controller) *mocks.Mockclient {
				return newMockClient(ctrl)
			},
			in:      "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			l := logger.New(nil, logger.WithEnabled(false))

			a, err := assistant.New(ctx, tc.clientFn(ctrl), l)
			assert.NoError(t, err)

			out, err := a.SendChatMessage(ctx, tc.inChat, tc.in)
			assert.Equal(t, tc.out, out)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if tc.inChat != nil {
					assert.Equal(t, tc.outChat, tc.inChat)
				}
			}
		})
	}
}

func newMockClient(ctrl *gomock.Controller) *mocks.Mockclient {
	c := mocks.NewMockclient(ctrl)
	c.EXPECT().
		ModelExists(gomock.Any()).
		Return(true, nil)
	return c
}
