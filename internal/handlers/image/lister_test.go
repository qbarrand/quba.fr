package image

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/golang/mock/gomock"
	"github.com/qbarrand/quba.fr/internal/generated/mock_images"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newImageLister(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockMetadataFS := mock_images.NewMockMetadataFS(ctrl)

	t.Run("error opening the filesystem's root", func(t *testing.T) {
		mockMetadataFS.EXPECT().Open(".").Return(nil, errors.New("random-error"))

		_, err := Lister(mockMetadataFS, nil)
		assert.Error(t, err)
	})

	t.Run("should work as expected", func(t *testing.T) {
		fsys := fstest.MapFS{
			"fileA": &fstest.MapFile{},
			"fileB": &fstest.MapFile{},
		}

		handler, err := Lister(fsys, nil)

		require.NoError(t, err)
		require.NotNil(t, handler)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		handler.ServeHTTP(w, req)

		body, err := io.ReadAll(w.Result().Body)
		require.NoError(t, err)

		assert.Equal(t, "application/json", w.Result().Header.Get("Content-Type"))
		assert.JSONEq(t, `["fileA", "fileB"]`, string(body))
	})
}
