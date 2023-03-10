package internal

import (
	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"os"
	"time"
)

const (
	cpuinfo = 2
	meminfo = 3
	lscpu   = 4
)

type Inode struct {
	id         fuseops.InodeID
	name       *string
	attributes fuseops.InodeAttributes
	dir        bool
	children   []fuseutil.Dirent
}

func (inode *Inode) findChildInode(name string) (fuseops.InodeID, error) {
	l := len(inode.children)
	if l == 0 {
		return 0, fuse.ENOENT
	}
	for _, child := range inode.children {
		if child.Name == name {
			return child.Inode, nil
		}
	}
	return 0, fuse.ENOENT
}

func toInode(dirent fuseutil.Dirent) *Inode {
	now := time.Now()
	inode := &Inode{
		id:   dirent.Inode,
		name: PString(dirent.Name),
		attributes: fuseops.InodeAttributes{
			Mode:  0444,
			Atime: now,
			Mtime: now,
			Ctime: now,
		},
	}
	return inode
}

func InitRootInode() *Inode {
	now := time.Now()
	root := &Inode{
		id: fuseops.RootInodeID,
		attributes: fuseops.InodeAttributes{
			Size:  4096,
			Mode:  os.ModeDir,
			Atime: now,
			Mtime: now,
			Ctime: now,
		},
		dir: true,
		children: []fuseutil.Dirent{
			{
				Offset: 1,
				Inode:  cpuinfo,
				Name:   "cpuinfo",
				Type:   fuseutil.DT_File,
			},
			{
				Offset: 2,
				Inode:  meminfo,
				Name:   "meminfo",
				Type:   fuseutil.DT_File,
			},
			{
				Offset: 3,
				Inode:  lscpu,
				Name:   "lscpu",
				Type:   fuseutil.DT_File,
			},
		},
	}
	return root
}
