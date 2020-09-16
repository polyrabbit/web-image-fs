package fs

import (
	"context"
	"fmt"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/polyrabbit/imagefs/webpage"
	"github.com/sirupsen/logrus"
)

// Set file owners to the current user,
// otherwise in OSX, we will fail to start.
var uid, gid uint32

func init() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	uid32, _ := strconv.ParseUint(u.Uid, 10, 32)
	gid32, _ := strconv.ParseUint(u.Gid, 10, 32)
	uid = uint32(uid32)
	gid = uint32(gid32)
}

// A tree node in filesystem, it acts as both a directory and file
type Node struct {
	fs.Inode
	client   *webpage.HTTPClient
	domNode  webpage.DomNode // An embeded dom node
	children map[string]webpage.DomNode
}

// NewRoot returns a file node - acting as a root, with inode sets to 1 and leaf sets to false
func NewRoot(client *webpage.HTTPClient, rootDom webpage.DomNode) *Node {
	return &Node{
		client:  client,
		domNode: rootDom,
	}
}

// List keys under a certain prefix from etcd, and output the next hierarchy level
func (n *Node) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	logrus.WithField("url", n.domNode.GetLink()).Debug("Node Readdir")
	domNodes, err := n.GetChildren(ctx)
	if err != nil {
		logrus.WithError(err).WithField("url", n.domNode.GetLink()).Errorf("Failed to get node's children")
		return nil, syscall.EIO
	}

	entries := make([]fuse.DirEntry, 0, len(domNodes))
	for _, dom := range domNodes {
		entries = append(entries, fuse.DirEntry{
			Mode: dom.FSMode(),
			Name: dom.FileName(),
			Ino:  dom.InodeHash(),
		})
	}
	return fs.NewListDirStream(entries), fs.OK
}

// Lookup finds a file under the current node(directory)
func (n *Node) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	logrus.WithField("name", name).WithField("url", n.domNode.GetLink()).Debug("Node Lookup")
	children, err := n.GetChildren(ctx)
	if err != nil {
		logrus.WithError(err).WithField("url", n.domNode.GetLink()).Errorf("Failed to get node's children")
		return nil, syscall.EIO
	}
	childNode := children[name]
	if childNode == nil {
		return nil, syscall.ENOENT
	}
	child := Node{
		domNode: childNode,
		client:  n.client,
	}
	return n.NewInode(ctx, &child, fs.StableAttr{Mode: childNode.FSMode(), Ino: childNode.InodeHash()}), fs.OK
}

func (n *Node) GetChildren(ctx context.Context) (map[string]webpage.DomNode, error) {
	if len(n.children) == 0 {
		domNodes, err := n.client.Parse(ctx, n.domNode.GetLink())
		if err != nil {
			return nil, fmt.Errorf("parser web page: %w", err)
		}
		groupedDoms := make(map[string]webpage.DomNode, len(domNodes))
		for _, dom := range domNodes {
			groupedDoms[dom.FileName()] = dom
		}
		n.children = groupedDoms
	}
	return n.children, nil
}

// Getattr outputs file attributes
// TODO: how to invalidate them?
func (n *Node) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = n.domNode.FSMode()
	out.Size = n.domNode.GetSize()
	out.Ino = n.domNode.InodeHash()
	now := time.Now()
	out.SetTimes(&now, &now, &now)
	out.Uid = uid
	out.Gid = gid
	return fs.OK
}

// Open gets value via http, and saves it in "content" for later read
func (n *Node) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	logrus.WithField("url", n.domNode.GetLink()).Debug("Node Open")
	if imgNode, ok := n.domNode.(*webpage.ImageNode); ok {
		if len(imgNode.Content) == 0 {
			content, err := n.client.Download(ctx, imgNode.GetLink())
			if err != nil {
				logrus.WithError(err).Errorf("Failed to download image content")
				return nil, 0, syscall.EIO
			}
			imgNode.Content = content
		}
	}
	return n, fuse.FOPEN_KEEP_CACHE, fs.OK
}

//
// Read returns bytes from "content", which should be filled by a prior Open operation
func (n *Node) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	logrus.WithField("url", n.domNode.GetLink()).Debug("Node Read")
	if imgNode, ok := n.domNode.(*webpage.ImageNode); !ok {
		return nil, syscall.EIO
	} else {
		end := int(off) + len(dest)
		if end > len(imgNode.Content) {
			end = len(imgNode.Content)
		}
		// We could copy to the `dest` buffer, but since we have a
		// []byte already, return that.
		return fuse.ReadResultData(imgNode.Content[off:end]), fs.OK
	}
}

var (
	_ fs.NodeGetattrer = &Node{}
	_ fs.NodeReaddirer = &Node{}
	_ fs.NodeLookuper  = &Node{}
	_ fs.NodeOpener    = &Node{}
	_ fs.FileReader    = &Node{}
)
