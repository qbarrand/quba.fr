package handlers

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/qbarrand/quba.fr/internal/image/mock_image"
)

func Test_newImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProcessor := mock_image.NewMockProcessor(ctrl)

	t.Run("processor returns an error", func(t *testing.T) {
		randomError := errors.New("random error")

		mockProcessor.EXPECT().Init().Return(randomError)

		_, err := newImage(mockProcessor, "", nil)
		assert.True(t, errors.Is(err, randomError))
	})

	t.Run("works as expected", func(t *testing.T) {
		const path = "/some/path"
		logger, _ := test.NewNullLogger()

		mockProcessor.EXPECT().Init()

		expected := &image{
			logger:    logger,
			path:      path,
			processor: mockProcessor,
		}

		i, err := newImage(mockProcessor, path, logger)

		assert.NoError(t, err)
		assert.Equal(t, expected, i)
	})
}
