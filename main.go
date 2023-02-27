package main

import (
	"context"
	"github.com/jacobsa/fuse"
	. "github.com/lnzx/wawu/internal"
	"log"
)

func main() {
	fs := NewWawu()
	mfs, err := fuse.Mount(DEFAULT_DIR, fs, &fuse.MountConfig{
		FSName:                    DEFAULT_FSNAME,
		ReadOnly:                  true,
		ErrorLogger:               nil,
		DebugLogger:               nil,
		DisableWritebackCaching:   true,
		EnableVnodeCaching:        false,
		EnableSymlinkCaching:      false,
		EnableNoOpenSupport:       true,
		EnableNoOpendirSupport:    true,
		DisableDefaultPermissions: false,
		UseVectoredRead:           false,
		Options:                   nil,
		Subtype:                   "wawu",
		EnableAsyncReads:          false,
	})
	if err != nil {
		log.Fatalf("Mount: %v", err)
	}
	log.Println("Successfully Mounted", DEFAULT_DIR)
	// Wait for it to be unmounted.
	if err = mfs.Join(context.Background()); err != nil {
		log.Fatalf("Join: %v", err)
	}
	log.Println("Successfully exiting.")
}
