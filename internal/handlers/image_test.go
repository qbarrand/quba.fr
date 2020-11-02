package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

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
	mockMeta := mock_image.NewMockMetaDB(ctrl)

	t.Run("processor returns an error", func(t *testing.T) {
		randomError := errors.New("random error")

		mockProcessor.EXPECT().Init().Return(randomError)

		_, err := newImage(mockProcessor, "", nil, nil)
		assert.True(t, errors.Is(err, randomError))
	})

	t.Run("works as expected", func(t *testing.T) {
		const path = "/some/path"
		logger, _ := test.NewNullLogger()

		mockProcessor.EXPECT().Init()

		expected := &image{
			logger:    logger,
			metaDB:    mockMeta,
			path:      path,
			processor: mockProcessor,
		}

		i, err := newImage(mockProcessor, path, mockMeta, logger)

		assert.NoError(t, err)
		assert.Equal(t, expected, i)
	})
}

func TestImage_ServeHTTP(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockProcessor := mock_image.NewMockProcessor(ctrl)
	mockMetaDB := mock_image.NewMockMetaDB(ctrl)

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

		i, err := newImage(mockProcessor, "basepath", mockMetaDB, logger)

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
			mockProcessor.EXPECT().BestFormats().Return([]img.Format{img.Webp}),
			mockHandler.EXPECT().SetFormat(img.Webp).Return(randomError),
			mockHandler.EXPECT().Destroy(),
		)

		logger, _ := test.NewNullLogger()

		i, err := newImage(mockProcessor, "basepath", mockMetaDB, logger)

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
			mockProcessor.EXPECT().BestFormats().Return([]img.Format{img.Webp}),
			mockHandler.EXPECT().SetFormat(img.Webp),
			mockHandler.EXPECT().Resize(ctx, width, 0).Return(randomError),
			mockHandler.EXPECT().Destroy(),
		)

		logger, _ := test.NewNullLogger()

		i, err := newImage(mockProcessor, "basepath", mockMetaDB, logger)

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
			mockProcessor.EXPECT().BestFormats().Return([]img.Format{img.Webp}),
			mockHandler.EXPECT().SetFormat(img.Webp),
			mockHandler.EXPECT().Bytes().Return(nil, randomError),
			mockHandler.EXPECT().Destroy(),
		)

		logger, _ := test.NewNullLogger()

		i, err := newImage(mockProcessor, "basepath", mockMetaDB, logger)

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

		const (
			location = "some-location"
			secs     = 1603842450
		)

		meta := &img.Metadata{
			Date:     time.Unix(secs, 0),
			Location: location,
		}

		gomock.InOrder(
			mockProcessor.EXPECT().Init(),
			mockProcessor.EXPECT().NewImageHandler("basepath/image_1.jpg").Return(mockHandler, nil),
			mockProcessor.EXPECT().BestFormats().Return([]img.Format{img.Webp}),
			mockHandler.EXPECT().SetFormat(img.Webp),
			mockHandler.EXPECT().Resize(ctx, width, 0),
			mockHandler.EXPECT().Bytes().Return(buf, nil),
			mockMetaDB.EXPECT().GetMetadata("image_1.jpg").Return(meta, nil),
			mockHandler.EXPECT().Destroy(),
		)

		logger, _ := test.NewNullLogger()

		i, err := newImage(mockProcessor, "basepath", mockMetaDB, logger)

		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, "/images/image_1.jpg", nil)
		r.Header.Set(acceptHeader, mimeWebp)
		r.Form = make(url.Values)
		r.Form.Set("width", strconv.Itoa(width))
		r = r.WithContext(ctx)

		w := assertStatusCode(t, i, r, http.StatusOK)

		resHeader := w.Result().Header

		assert.Equal(t, mimeWebp, resHeader.Get(contentTypeHeader))
		// `abcd`'s fnv hash is b9de7375
		assert.Equal(t, "b9de7375", resHeader.Get("ETag"))
		// `abcd` is 4 bytes
		assert.Equal(t, "4", resHeader.Get("Content-Length"))
		assert.Equal(t, location, resHeader.Get("X-Quba-Location"))
		assert.Equal(t, strconv.Itoa(secs), resHeader.Get("X-Quba-Date"))
		assert.Equal(t, buf, w.Body.Bytes())
	})
}

func Test_newImageLister(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	metaDB := mock_image.NewMockMetaDB(ctrl)

	t.Run("returns an error", func(t *testing.T) {
		randomError := errors.New("random-error")

		metaDB.EXPECT().AllNames().Return(nil, randomError)

		_, err := newImageLister(metaDB, nil)
		assert.True(t, errors.Is(err, randomError))
	})

	t.Run("should work as expected", func(t *testing.T) {
		metaDB.EXPECT().AllNames().Return([]string{"name-1", "name-2"}, nil)

		logger, _ := test.NewNullLogger()

		il, err := newImageLister(metaDB, logger)

		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		il.ServeHTTP(w, req)

		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(t, `["name-1","name-2"]`+"\n", w.Body.String())
	})
}

func Test_getBestFormat(t *testing.T) {
	t.Run("empty server formats", func(t *testing.T) {
		_, err := getBestFormat(nil, []string{"image/webp"})
		assert.Error(t, err)
	})

	t.Run("empty server formats", func(t *testing.T) {
		_, err := getBestFormat(nil, []string{"image/webp"})
		assert.Error(t, err)
	})
}

func Test_getMIMETypes(t *testing.T) {
	t.Run("one header, no value", func(t *testing.T) {
		ct, err := getMIMETypes(nil)
		require.NoError(t, err)
		assert.Empty(t, ct)
	})

	t.Run("one header, one invalid value", func(t *testing.T) {
		_, err := getMIMETypes([]string{"abcd def"})
		require.Error(t, err)
	})

	t.Run("one header, one valid value", func(t *testing.T) {
		headers := []string{"image/jpeg"}

		ct, err := getMIMETypes(headers)

		require.NoError(t, err)
		assert.Equal(t, headers, ct)
	})

	const (
		jpeg = "image/jpeg"
		webp = "image/webp"
	)

	t.Run("one header, two valid values (comma)", func(t *testing.T) {
		ct, err := getMIMETypes([]string{jpeg + "," + webp})

		require.NoError(t, err)
		assert.Equal(t, []string{jpeg, webp}, ct)
	})

	t.Run("one header, two valid values (comma space)", func(t *testing.T) {
		ct, err := getMIMETypes([]string{jpeg + ", " + webp})

		require.NoError(t, err)
		assert.Equal(t, []string{jpeg, webp}, ct)
	})

	t.Run("two headers, two valid values", func(t *testing.T) {
		s := []string{jpeg, webp}

		ct, err := getMIMETypes(s)

		require.NoError(t, err)
		assert.Equal(t, s, ct)
	})

	t.Run("two headers, three valid values", func(t *testing.T) {
		const png = "image/png"

		ct, err := getMIMETypes([]string{jpeg, webp + "," + png})

		require.NoError(t, err)
		assert.Equal(t, []string{jpeg, webp, png}, ct)
	})
}
