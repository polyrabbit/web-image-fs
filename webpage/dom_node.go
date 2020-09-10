package webpage

import (
	"hash/fnv"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
)

type DomNode interface {
	GetSize() uint64
	FileName() string
	GetLink() string
	IsImage() bool
	// Get self stat as a fs-mode
	FSMode() uint32
	// Hash file path into inode number, so we can ensure the same file always gets the same inode number
	InodeHash() uint64
}

// LinkNode represents a a-link node
type LinkNode struct {
	Name     string
	SelfLink string
}

func (n *LinkNode) FileName() string {
	if n.Name == "" {
		n.Name = filepath.Base(n.SelfLink)
	}
	return strings.TrimSpace(n.Name)
}

func (n *LinkNode) GetSize() uint64 {
	return 0
}

func (n *LinkNode) IsImage() bool {
	return false
}

func (n *LinkNode) GetLink() string {
	return n.SelfLink
}

func (n *LinkNode) FSMode() uint32 {
	return 0755 | uint32(syscall.S_IFDIR)
}

func (n *LinkNode) InodeHash() uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(n.SelfLink))
	return h.Sum64()
}

// ImageNode represents a image-node
type ImageNode struct {
	LinkNode
	Size        uint64
	contentType string

	rwMu    sync.RWMutex // Protect file content
	content []byte       // Internal buffer to hold the current file content
}

func (n *ImageNode) FileName() string {
	fname := n.LinkNode.FileName()
	if filepath.Ext(fname) == "" {
		linkExt := filepath.Ext(n.SelfLink)
		if len(linkExt) != 0 {
			fname += linkExt
		} else {
			//TODO: deduce by content-type
		}
	}
	return fname
}

func (n *ImageNode) GetSize() uint64 {
	return n.Size
}

func (n *ImageNode) IsImage() bool {
	return true
}

func (n *ImageNode) FSMode() uint32 {
	return 0644 | uint32(syscall.S_IFREG)
}

var (
	_ DomNode = &LinkNode{}
	_ DomNode = &ImageNode{}
)
