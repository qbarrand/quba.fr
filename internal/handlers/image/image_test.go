package image

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"testing/fstest"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/qbarrand/quba.fr/data/images"
	"github.com/qbarrand/quba.fr/internal/generated/mock_image"
	"github.com/qbarrand/quba.fr/internal/generated/mock_images"
	"github.com/qbarrand/quba.fr/internal/imgpro"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockProcessor := mock_image.NewMockProcessor(ctrl)

	t.Run("processor returns an error", func(t *testing.T) {
		randomError := errors.New("random error")

		mockProcessor.EXPECT().Init().Return(randomError)

		_, err := New(mockProcessor, nil, nil)
		assert.True(t, errors.Is(err, randomError))
	})

	t.Run("works as expected", func(t *testing.T) {
		logger, _ := test.NewNullLogger()

		mockProcessor.EXPECT().Init()
		mockMetadataFS := mock_images.NewMockMetadataFS(ctrl)

		expected := &Image{
			logger:    logger,
			mfs:       mockMetadataFS,
			processor: mockProcessor,
		}

		i, err := New(mockProcessor, mockMetadataFS, logger)

		assert.NoError(t, err)
		assert.Equal(t, expected, i)
	})
}

func TestImage_ServeHTTP(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockProcessor := mock_image.NewMockProcessor(ctrl)
	mockMetadataFS := mock_images.NewMockMetadataFS(ctrl)

	const (
		acceptHeader      = "Accept"
		contentTypeHeader = "Content-Type"
		filename          = "image_1.jpg"
		target            = "/" + filename
		mimeWebp          = "image/webp"
	)

	randomError := errors.New("random-error")

	t.Run("cannot open the image file", func(t *testing.T) {
		gomock.InOrder(
			mockProcessor.EXPECT().Init(),
			mockMetadataFS.EXPECT().OpenWithMetadata(filename).Return(nil, nil, randomError),
		)

		logger, _ := test.NewNullLogger()

		i, err := New(mockProcessor, mockMetadataFS, logger)

		require.NoError(t, err)
		assert.HTTPStatusCode(
			t,
			i.ServeHTTP,
			http.MethodGet,
			target,
			nil,
			http.StatusInternalServerError)
	})

	getFsFile := func(t *testing.T) ([]byte, fs.File) {
		t.Helper()

		b := []byte{1, 2, 3}

		// Get an actual fs.File
		fsTemp := fstest.MapFS{
			filename: &fstest.MapFile{Data: b},
		}

		fd, err := fsTemp.Open(filename)
		require.NoError(t, err)

		return b, fd
	}

	t.Run("processor cannot create a new handler", func(t *testing.T) {
		b, fd := getFsFile(t)
		defer fd.Close()

		gomock.InOrder(
			mockProcessor.EXPECT().Init(),
			mockMetadataFS.EXPECT().OpenWithMetadata(filename).Return(fd, &images.Metadata{}, nil),
			mockProcessor.EXPECT().HandlerFromBytes(b).Return(nil, randomError),
		)

		logger, _ := test.NewNullLogger()

		i, err := New(mockProcessor, mockMetadataFS, logger)

		require.NoError(t, err)
		assert.HTTPStatusCode(
			t,
			i.ServeHTTP,
			http.MethodGet,
			target,
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

		b, fd := getFsFile(t)
		defer fd.Close()

		gomock.InOrder(
			mockProcessor.EXPECT().Init(),
			mockMetadataFS.EXPECT().OpenWithMetadata(filename).Return(fd, &images.Metadata{}, nil),
			mockProcessor.EXPECT().HandlerFromBytes(b).Return(mockHandler, nil),
			mockProcessor.EXPECT().BestFormats().Return([]imgpro.Format{imgpro.Webp}),
			mockHandler.EXPECT().SetFormat(imgpro.Webp).Return(randomError),
			mockHandler.EXPECT().Destroy(),
		)

		logger, _ := test.NewNullLogger()

		i, err := New(mockProcessor, mockMetadataFS, logger)

		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, target, nil)
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

		b, fd := getFsFile(t)
		defer fd.Close()

		gomock.InOrder(
			mockProcessor.EXPECT().Init(),
			mockMetadataFS.EXPECT().OpenWithMetadata(filename).Return(fd, &images.Metadata{}, nil),
			mockProcessor.EXPECT().HandlerFromBytes(b).Return(mockHandler, nil),
			mockProcessor.EXPECT().BestFormats().Return([]imgpro.Format{imgpro.Webp}),
			mockHandler.EXPECT().SetFormat(imgpro.Webp),
			mockHandler.EXPECT().Resize(ctx, width, 0).Return(randomError),
			mockHandler.EXPECT().Destroy(),
		)

		logger, _ := test.NewNullLogger()

		i, err := New(mockProcessor, mockMetadataFS, logger)

		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, target, nil)
		r.Header.Set(acceptHeader, mimeWebp)
		r.Form = make(url.Values)
		r.Form.Set("width", strconv.Itoa(width))
		r = r.WithContext(ctx)

		assertStatusCode(t, i, r, http.StatusInternalServerError)
	})

	t.Run("no resize + cannot get bytes", func(t *testing.T) {
		ctx := getContext(t)

		mockHandler := mock_image.NewMockHandler(ctrl)

		b, fd := getFsFile(t)
		defer fd.Close()

		gomock.InOrder(
			mockProcessor.EXPECT().Init(),
			mockMetadataFS.EXPECT().OpenWithMetadata(filename).Return(fd, &images.Metadata{}, nil),
			mockProcessor.EXPECT().HandlerFromBytes(b).Return(mockHandler, nil),
			mockProcessor.EXPECT().BestFormats().Return([]imgpro.Format{imgpro.Webp}),
			mockHandler.EXPECT().SetFormat(imgpro.Webp),
			mockHandler.EXPECT().Bytes().Return(nil, randomError),
			mockHandler.EXPECT().Destroy(),
		)

		logger, _ := test.NewNullLogger()

		i, err := New(mockProcessor, mockMetadataFS, logger)

		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, target, nil)
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
			location  = "some-location"
			secs      = 1603842450
			mainColor = "#012345"
		)

		b, fd := getFsFile(t)
		defer fd.Close()

		meta := &images.Metadata{
			Date:      time.Unix(secs, 0),
			Location:  location,
			MainColor: mainColor,
		}

		gomock.InOrder(
			mockProcessor.EXPECT().Init(),
			mockMetadataFS.EXPECT().OpenWithMetadata(filename).Return(fd, meta, nil),
			mockProcessor.EXPECT().HandlerFromBytes(b).Return(mockHandler, nil),
			mockProcessor.EXPECT().BestFormats().Return([]imgpro.Format{imgpro.Webp}),
			mockHandler.EXPECT().SetFormat(imgpro.Webp),
			mockHandler.EXPECT().Resize(ctx, width, 0),
			mockHandler.EXPECT().Bytes().Return(buf, nil),
			mockHandler.EXPECT().Destroy(),
		)

		logger, _ := test.NewNullLogger()

		i, err := New(mockProcessor, mockMetadataFS, logger)

		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, target, nil)
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
		assert.Equal(t, mainColor, resHeader.Get("X-Quba-Main-Color"))
		assert.Equal(t, buf, w.Body.Bytes())
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
