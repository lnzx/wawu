package internal

import (
	"context"
	"fmt"
	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strings"
)

const (
	DEFAULT_DIR = "/var/lib/wawu"
)

var suffixes = [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n"}
var DEFAULT_FSNAME = "/dev/sd"

func switchFsname() {
	for _, suffix := range suffixes {
		cur := DEFAULT_FSNAME + suffix
		if _, err := os.Stat(cur); err != nil && os.IsNotExist(err) {
			DEFAULT_FSNAME = cur
			log.Println("DEFAULT_FSNAME:", DEFAULT_FSNAME)
			return
		}
	}
	log.Fatalln("ERROR: Failed to generate DEFAULT_FSNAME")
}

func init() {
	_, err := os.Stat(DEFAULT_DIR)
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(DEFAULT_DIR, os.ModeDir)
		if err != nil {
			log.Fatalln(err)
			return
		}
		log.Println("init dir:", DEFAULT_DIR)
	}

	switchFsname()
}

type Wawu struct {
	fuseutil.NotImplementedFileSystem
	inodes map[fuseops.InodeID]*Inode
}

func NewWawu() fuse.Server {
	fs := &Wawu{
		inodes: make(map[fuseops.InodeID]*Inode),
	}

	root := InitRootInode()
	fs.inodes[fuseops.RootInodeID] = root
	for _, dirent := range root.children {
		fs.inodes[dirent.Inode] = toInode(dirent)
	}
	return fuseutil.NewFileSystemServer(fs)
}

func (fs *Wawu) StatFS(ctx context.Context, op *fuseops.StatFSOp) (err error) {
	used := GetDiskUsed()
	const BLOCK_SIZE = 4096
	const TOTAL_SPACE = 1 * 1024 * 1024 * 1024 * 1024 // 1TB
	const TOTAL_BLOCKS = TOTAL_SPACE / BLOCK_SIZE
	const INODES = 1 * 1000 * 1000 // 10 million

	usedBlocks := used / BLOCK_SIZE

	op.BlockSize = BLOCK_SIZE
	op.Blocks = TOTAL_BLOCKS
	op.BlocksFree = TOTAL_BLOCKS - usedBlocks
	op.BlocksAvailable = uint64(float64(TOTAL_BLOCKS)*0.9526627218934911 - float64(usedBlocks))
	op.IoSize = 1 * 1024 * 1024 // 1MB
	op.Inodes = INODES
	op.InodesFree = uint64(INODES - len(fs.inodes))
	return
}

func (fs *Wawu) getInodeOrDie(id fuseops.InodeID) *Inode {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("stacktrace from panic: %v \n"+string(debug.Stack()), err)
			err = fuse.EIO
		}
	}()
	inode := fs.inodes[id]
	if inode == nil {
		panic(fmt.Sprintf("Unknown inode: %v", id))
	}
	return inode
}

func (fs *Wawu) GetInodeAttributes(ctx context.Context, op *fuseops.GetInodeAttributesOp) error {
	inode := fs.getInodeOrDie(op.Inode)

	op.Attributes = inode.attributes
	return nil
}

func (fs *Wawu) LookUpInode(ctx context.Context, op *fuseops.LookUpInodeOp) error {
	parent := fs.getInodeOrDie(op.Parent)
	child, err := parent.findChildInode(op.Name)
	if err != nil {
		return err
	}

	// Copy over information.
	op.Entry.Child = child
	op.Entry.Attributes = fs.inodes[child].attributes
	return nil
}

func (fs *Wawu) ReadDir(ctx context.Context, op *fuseops.ReadDirOp) error {
	inode := fs.getInodeOrDie(op.Inode)
	if !inode.dir {
		return fuse.EIO
	}
	entries := inode.children
	// Grab the range of interest.
	if op.Offset > fuseops.DirOffset(len(entries)) {
		return fuse.EIO
	}
	entries = entries[op.Offset:]
	// Resume at the specified offset into the array.
	for _, e := range entries {
		n := fuseutil.WriteDirent(op.Dst[op.BytesRead:], e)
		if n == 0 {
			break
		}
		op.BytesRead += n
	}
	return nil
}

func (fs *Wawu) OpenFile(ctx context.Context, op *fuseops.OpenFileOp) error {
	// Allow opening any file.
	op.KeepPageCache = false
	op.UseDirectIO = true
	return nil
}

func (fs *Wawu) ReadFile(ctx context.Context, op *fuseops.ReadFileOp) error {
	var info string
	switch op.Inode {
	case meminfo:
		info = GetMeminfo()
	default:
		log.Println("not supported")
		return fuse.ENOENT
	}
	reader := strings.NewReader(info)
	var err error
	op.BytesRead, err = reader.ReadAt(op.Dst, op.Offset)
	// Special case: FUSE doesn't expect us to return io.EOF.
	if err == io.EOF {
		return nil
	}
	return nil
}
