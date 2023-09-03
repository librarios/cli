package yes24

import (
	"reflect"
	"testing"
)

func TestExtractPreviewImageFilenames(t *testing.T) {
	html := `
	'pages': { 'PAGE' : [{"GoodsNo":121961003,"OrderNo":1,"OriginalName":null,"SmallImage":{"Name":"3pvvjqvluxau6ygp01.jpg","Width":400,"Height":514,"FileSize":0},"MiddleImage":{"Name":"3pvvjqvluxau6ygp01.jpg","Width":544,"Height":700,"FileSize":0},"LargeImage":{"Name":"3pvvjqvluxau6ygp01.jpg","Width":1089,"Height":1400,"FileSize":0},"Bookmark":""},{"GoodsNo":121961003,"OrderNo":2,"OriginalName":null,"SmallImage":{"Name":"kt2vqir9uqvbc9ra02.jpg","Width":400,"Height":514,"FileSize":0},"MiddleImage":{"Name":"kt2vqir9uqvbc9ra02.jpg","Width":544,"Height":700,"FileSize":0},"LargeImage":{"Name":"kt2vqir9uqvbc9ra02.jpg","Width":1089,"Height":1400,"FileSize":0},"Bookmark":""}] },
	`

	expected := []string{"3pvvjqvluxau6ygp01.jpg", "kt2vqir9uqvbc9ra02.jpg"}
	filenames, err := extractPreviewImageFilenames(html)

	if err != nil {
		t.Errorf("failed to extract filenames: %+v", err)
	}

	if !reflect.DeepEqual(expected, filenames) {
		t.Errorf("mismatch. expected: %+v, actual: %+v", expected, filenames)
	}
}
