//go:generate go run github.com/golang/mock/mockgen -destination mock_images/images.go github.com/qbarrand/quba.fr/data/images MetadataFS
//go:generate go run github.com/golang/mock/mockgen -destination mock_imgpro/processor.go github.com/qbarrand/quba.fr/internal/imgpro Handler,Processor

package generated
