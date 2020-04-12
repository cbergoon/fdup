package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

func shortenString(s string, l int) string {
	if len(s) > l {
		return fmt.Sprintf("%s...%s", s[:20], s[len(s)-70:])
	}
	return s
}

func hashDirectory(path string, count int64) string {
	hashes := [][]byte{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("error: %+v", err)
		} else {
			fmt.Printf("\033[2K\r\tds %d -> %s", count, shortenString(path, 90))

			if !info.IsDir() && (info.Mode()&os.ModeSymlink) != os.ModeSymlink {
				b, err := ioutil.ReadFile(path)
				if err != nil {
					log.Printf("error: %+v", err)
				}
				digest := sha1.New()
				digest.Write(b)
				hashes = append(hashes, digest.Sum(nil))
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("error: %+v", err)
	}
	if len(hashes) > 0 {
		digest := sha1.New()
		for _, h := range hashes {
			digest.Write(h)
		}
		return hex.EncodeToString(digest.Sum(nil))
	} else {
		return ""
	}
}

func wrapMainWalk(fmap map[string][]string, dmap map[string][]string, dirComp bool) func(path string, info os.FileInfo, err error) error {
	var count int64 = 0
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("error: %+v", err)
		} else {
			count++
			fmt.Printf("\033[2K\r\t%d -> %s", count, shortenString(path, 90))

			if dirComp && info.IsDir() && (info.Mode()&os.ModeSymlink) != os.ModeSymlink {
				dh := hashDirectory(path, count)
				dmap[dh] = append(dmap[dh], path)
			} else if !info.IsDir() && (info.Mode()&os.ModeSymlink) != os.ModeSymlink {
				b, err := ioutil.ReadFile(path)
				if err != nil {
					log.Printf("error: %+v", err)
				}
				digest := sha1.New()
				digest.Write(b)
				h := digest.Sum(nil)
				fmap[hex.EncodeToString(h)] = append(fmap[hex.EncodeToString(h)], path)
			}
		}
		return nil
	}
}

func displayDuplicated(m map[string][]string, workingDirectory string) {
	for k, v := range m {
		if len(v) > 1 {
			fmt.Printf("\t%s: \n", fmt.Sprintf("%s...%s", k[:3], k[35:]))
			for _, f := range v {
				v := strings.TrimPrefix(f, workingDirectory)
				fmt.Printf("\t\t -> %s\n", v)
			}
		}
	}
}

func displayStats(fmap map[string][]string, dmap map[string][]string, runtime time.Duration, dirComp bool) {
	duplicateFiles := 0
	filesScanned := 0
	duplicatedSize := int64(0)
	for _, v := range fmap {
		if len(v) > 1 {
			duplicateFiles += len(v) - 1
		}
		filesScanned += len(v)

		if len(v) > 1 {
			f, err := os.Open(v[0])
			if err != nil {
				log.Printf("error: %+v", err)
			} else {
				info, err := f.Stat()
				if err != nil {
					log.Printf("error: %+v", err)
				} else {
					x := info.Size()
					duplicatedSize += (x * int64(len(v)-1))
				}
			}
			f.Close()
		}
	}
	duplicateDirectories := 0
	directoriesScanned := 0
	if dirComp {
		for _, v := range dmap {
			if len(v) > 1 {
				duplicateDirectories += len(v) - 1
			}
			directoriesScanned += len(v)
		}
	}

	fmt.Printf("\n")
	fmt.Printf("Runtime: \t%s\n", runtime)
	if dirComp {
		fmt.Printf("Scanned:\t%d directories\t%d files\n", directoriesScanned, filesScanned)
		fmt.Printf("Duplicates:\t%d directories\t%d files\n", duplicateDirectories, duplicateFiles)
	} else {
		fmt.Printf("Scanned:\t%d files\n", filesScanned)
		fmt.Printf("Duplicates:\t%d files\n", duplicateFiles)
	}
	fmt.Printf("Duplicate Size: %s\n\n", humanize.IBytes(uint64(duplicatedSize)))

}

func main() {
	dirCompPtr := flag.Bool("dir-comparison", false, "performs directory level comparison (not reccomended for large directories)")

	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error: %+v", err)
	}

	fmt.Printf("Starting scan on: %s\n", wd)

	fmap := make(map[string][]string)
	dmap := make(map[string][]string)

	start := time.Now()
	filepath.Walk(wd, wrapMainWalk(fmap, dmap, *dirCompPtr))
	fmt.Printf("\033[2K\r")
	runtime := time.Since(start)

	if *dirCompPtr {
		fmt.Printf("\nDuplicated Directories in %s\n", wd)
		displayDuplicated(dmap, wd)
	}
	fmt.Printf("\nDuplicated Files in %s\n", wd)
	displayDuplicated(fmap, wd)

	displayStats(fmap, dmap, runtime, *dirCompPtr)
}
