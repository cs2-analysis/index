package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/dolmen-go/jsonmap"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

func checkError(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

type File struct {
	Name     string  `json:"filename"`
	Path     string  `json:"path"`
	Manifest string  `json:"manifest"`
	Updated  string  `json:"updated"`
	Deleted  *bool   `json:"deleted,omitempty"`
	Size     *string `json:"size,omitempty"`
	Sha256   *string `json:"sha256,omitempty"`
}

func main() {
	url := flag.String("r", "https://github.com/cs2-analysis/deports.git", "repository url")
	branch := flag.String("b", "", "branch name")
	flag.Parse()

	if *url == "" {
		flag.Usage()
		checkError(fmt.Errorf("repository url is required"))
	}

	if *branch == "" {
		flag.Usage()
		checkError(fmt.Errorf("branch name is required"))
	}

	s := filesystem.NewStorage(memfs.New(), cache.NewObjectLRUDefault())
	r, err := git.Clone(s, memfs.New(), &git.CloneOptions{
		URL:           *url,
		ReferenceName: plumbing.NewBranchReferenceName(*branch),
		SingleBranch:  true,
		NoCheckout:    true,
		Progress:      os.Stderr,
	})
	checkError(err)

	b, err := r.Branch(*branch)
	checkError(err)

	ref, err := r.Reference(b.Merge, true)
	checkError(err)

	it, err := r.Log(&git.LogOptions{
		From: ref.Hash(),
		PathFilter: func(s string) bool {
			return strings.HasPrefix(s, "data/")
		},
	})
	checkError(err)

	m := make(map[string][]File)

	err = it.ForEach(func(c *object.Commit) error {
		wt, err := r.Worktree()
		if err != nil {
			return err
		}

		err = wt.Checkout(&git.CheckoutOptions{Hash: c.Hash, Force: true})
		if err != nil {
			return err
		}

		stats, err := c.Stats()
		if err != nil {
			return err
		}

		for _, s := range stats {
			if !strings.HasPrefix(s.Name, "data/") {
				continue
			}

			file, err := wt.Filesystem.Open(s.Name)
			if err != nil {
				return err
			}
			defer file.Close()

			content, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			var f File
			err = json.Unmarshal(content, &f)
			if err != nil {
				return err
			}

			m[f.Path] = append(m[f.Path], f)

		}

		return nil
	})
	checkError(err)

	mo := jsonmap.Ordered{
		Data:  make(map[string]interface{}),
		Order: []string{},
	}

	for k, v := range m {
		mo.Data[k] = v
		mo.Order = append(mo.Order, k)
	}
	sort.Strings(mo.Order)

	d, err := json.MarshalIndent(mo, "", "  ")
	checkError(err)

	fmt.Print(string(d))
}
