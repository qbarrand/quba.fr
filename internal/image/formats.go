package image

type Format uint

const (
	JPEG Format = iota
	Webp
)

//func AcceptHeaderToFormat(accept string) (Format, error) {
//	for _, mimeType := range strings.Split(accept, ",") {
//		mimeType = strings.Trim(mimeType, " ")
//
//		switch mimeType {
//		case "image/jpeg":
//			return JPEG, nil
//		case "image/webp":
//			return Webp, nil
//		}
//	}
//
//	return 0, fmt.Errorf("%v: unhandled Accept header", accept)
//}
