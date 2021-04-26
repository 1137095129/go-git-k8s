package public

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/wang1137095129/go-git-k8s/config"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type PublicGitHandler struct {
	lastPullTime        *time.Time
	localRepositoryPath string
}

var once = &sync.Once{}

var lock = &sync.WaitGroup{}

func (p *PublicGitHandler) OpenRepository(c *config.Config) (*git.Repository, error) {
	lock.Wait()
	lock.Add(1)
	defer lock.Done()

	once.Do(func() {
		p.lastPullTime = new(time.Time)
		p.localRepositoryPath = filepath.Join(c.Git.Local, c.Git.Repository)
	})

	if _, err := os.Stat(p.localRepositoryPath); err != nil && os.IsNotExist(err) {
		repository, err := git.PlainClone(
			c.Git.Local,
			false,
			&git.CloneOptions{
				URL:           c.Git.Url,
				RemoteName:    c.Git.Remote,
				ReferenceName: plumbing.NewBranchReferenceName(c.Git.Branch),
			},
		)

		if err != nil {
			return nil, err
		}

		head, err := repository.Head()
		if err != nil {
			return nil, err
		}

		commitIter, err := repository.Log(&git.LogOptions{From: head.Hash()})
		if err != nil {
			return nil, err
		}

		t := time.Time{}
		err = commitIter.ForEach(func(commit *object.Commit) error {
			if t.Before(commit.Committer.When) {
				t = commit.Committer.When
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		*p.lastPullTime = t

		return repository, nil
	} else if err != nil {
		return nil, err
	}

	repository, err := git.PlainOpen(c.Git.Local)
	if err != nil {
		return nil, err
	}

	err = repository.Fetch(&git.FetchOptions{
		RemoteName: c.Git.Remote,
	})
	if err != nil {
		return nil, err
	}

	revision, err := repository.ResolveRevision(plumbing.Revision(fmt.Sprintf("%s/%s", c.Git.Remote, c.Git.Branch)))
	if err != nil {
		return nil, err
	}

	commit, err := repository.CommitObject(*revision)
	if err != nil {
		return nil, err
	}

	if p.lastPullTime.Before(commit.Committer.When) {
		worktree, err := repository.Worktree()
		if err != nil {
			return nil, err
		}
		err = worktree.Pull(&git.PullOptions{
			RemoteName:    c.Git.Remote,
			ReferenceName: plumbing.NewBranchReferenceName(c.Git.Branch),
		})
		if err != nil {
			return nil, err
		}
		*p.lastPullTime = commit.Committer.When
	}

	return repository, nil
}
