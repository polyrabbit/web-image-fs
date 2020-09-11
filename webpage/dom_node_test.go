package webpage

import "testing"

func TestImageNode_FileName(t *testing.T) {
	type fields struct {
		LinkNode    LinkNode
		Size        uint64
		Content     []byte
		ContentType string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{{
		name: "Get name from name field",
		fields: fields{LinkNode: LinkNode{
			Name: "parent/AAA",
		}},
		want: "AAA",
	}, {
		name: "Get name from link field",
		fields: fields{LinkNode: LinkNode{
			SelfLink: "https://l35h2znmhf1scosj14ztuxt1-wpengine.netdna-ssl.com/wp-content/themes/unherdv3/src/img/share-twitter.png",
		}},
		want: "share-twitter.png",
	}, {
		name: "Get ext from link field",
		fields: fields{LinkNode: LinkNode{
			Name:     "AAA",
			SelfLink: "https://l35h2znmhf1scosj14ztuxt1-wpengine.netdna-ssl.com/wp-content/themes/unherdv3/src/img/share-twitter.png",
		}},
		want: "AAA.png",
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &ImageNode{
				LinkNode:    tt.fields.LinkNode,
				Size:        tt.fields.Size,
				Content:     tt.fields.Content,
				ContentType: tt.fields.ContentType,
			}
			if got := n.FileName(); got != tt.want {
				t.Errorf("FileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
