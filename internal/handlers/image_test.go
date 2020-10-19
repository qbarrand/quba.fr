package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	img "github.com/qbarrand/quba.fr/internal/image"
	"github.com/qbarrand/quba.fr/internal/image/mock_image"
)

func Test_newImage(t *testing.T) {
	t.Parallel()

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

func TestImage_ServeHTTP(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockProcessor := mock_image.NewMockProcessor(ctrl)

	const (
		acceptHeader      = "Accept"
		contentTypeHeader = "Content-Type"
		mimeWebp          = "image/webp"
	)

	randomError := errors.New("random-error")

	t.Run("processor cannot create a new handler", func(t *testing.T) {
		gomock.InOrder(
			mockProcessor.EXPECT().Init(),
			mockProcessor.EXPECT().NewImageHandler("basepath/image_1.jpg").Return(nil, randomError),
		)

		logger, _ := test.NewNullLogger()

		i, err := newImage(mockProcessor, "basepath", logger)

		require.NoError(t, err)
		assert.HTTPStatusCode(
			t,
			i.ServeHTTP,
			http.MethodGet,
			"/images/image_1.jpg",
			nil,
			http.StatusInternalServerError)
	})

	assertStatusCode := func(t *testing.T, h http.Handler, r *http.Request, statusCode int) *httptest.ResponseRecorder {
		t.Helper()

		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		require.Equal(t, statusCode, w.Result().StatusCode)

		return w
	}

	t.Run("handler cannot set the format", func(t *testing.T) {
		mockHandler := mock_image.NewMockHandler(ctrl)

		gomock.InOrder(
			mockProcessor.EXPECT().Init(),
			mockProcessor.EXPECT().NewImageHandler("basepath/image_1.jpg").Return(mockHandler, nil),
			mockHandler.EXPECT().SetFormat(img.Webp).Return(randomError),
			mockHandler.EXPECT().Destroy(),
		)

		logger, _ := test.NewNullLogger()

		i, err := newImage(mockProcessor, "basepath", logger)

		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, "/images/image_1.jpg", nil)
		r.Header.Set(acceptHeader, mimeWebp)

		assertStatusCode(t, i, r, http.StatusInternalServerError)
	})

	getContext := func(t *testing.T) context.Context {
		t.Helper()

		return context.WithValue(context.Background(), "test", t.Name())
	}

	t.Run("cannot resize", func(t *testing.T) {
		const width = 1234

		ctx := getContext(t)

		mockHandler := mock_image.NewMockHandler(ctrl)

		gomock.InOrder(
			mockProcessor.EXPECT().Init(),
			mockProcessor.EXPECT().NewImageHandler("basepath/image_1.jpg").Return(mockHandler, nil),
			mockHandler.EXPECT().SetFormat(img.Webp),
			mockHandler.EXPECT().Resize(ctx, width, 0).Return(randomError),
			mockHandler.EXPECT().Destroy(),
		)

		logger, _ := test.NewNullLogger()

		i, err := newImage(mockProcessor, "basepath", logger)

		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, "/images/image_1.jpg", nil)
		r.Header.Set(acceptHeader, mimeWebp)
		r.Form = make(url.Values)
		r.Form.Set("width", strconv.Itoa(width))
		r = r.WithContext(ctx)

		assertStatusCode(t, i, r, http.StatusInternalServerError)
	})

	t.Run("no resize + cannot get bytes", func(t *testing.T) {
		ctx := getContext(t)

		mockHandler := mock_image.NewMockHandler(ctrl)

		gomock.InOrder(
			mockProcessor.EXPECT().Init(),
			mockProcessor.EXPECT().NewImageHandler("basepath/image_1.jpg").Return(mockHandler, nil),
			mockHandler.EXPECT().SetFormat(img.Webp),
			mockHandler.EXPECT().Bytes().Return(nil, randomError),
			mockHandler.EXPECT().Destroy(),
		)

		logger, _ := test.NewNullLogger()

		i, err := newImage(mockProcessor, "basepath", logger)

		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, "/images/image_1.jpg", nil)
		r.Header.Set(acceptHeader, mimeWebp)
		r = r.WithContext(ctx)

		assertStatusCode(t, i, r, http.StatusInternalServerError)
	})

	t.Run("resize + return bytes", func(t *testing.T) {
		const width = 1234

		buf := []byte("abcd")
		ctx := getContext(t)

		mockHandler := mock_image.NewMockHandler(ctrl)

		gomock.InOrder(
			mockProcessor.EXPECT().Init(),
			mockProcessor.EXPECT().NewImageHandler("basepath/image_1.jpg").Return(mockHandler, nil),
			mockHandler.EXPECT().SetFormat(img.Webp),
			mockHandler.EXPECT().Resize(ctx, width, 0),
			mockHandler.EXPECT().Bytes().Return(buf, nil),
			mockHandler.EXPECT().Destroy(),
		)

		logger, _ := test.NewNullLogger()

		i, err := newImage(mockProcessor, "basepath", logger)

		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, "/images/image_1.jpg", nil)
		r.Header.Set(acceptHeader, mimeWebp)
		r.Form = make(url.Values)
		r.Form.Set("width", strconv.Itoa(width))
		r = r.WithContext(ctx)

		w := assertStatusCode(t, i, r, http.StatusOK)
		assert.Equal(t, mimeWebp, w.Result().Header.Get(contentTypeHeader))
		// `abcd`'s fnv hash is b9de7375
		assert.Equal(t, "b9de7375", w.Result().Header.Get("ETag"))
		// `abcd` is 4 bytes
		assert.Equal(t, "4", w.Result().Header.Get("Content-Length"))
		assert.Equal(t, buf, w.Body.Bytes())
	})
}
