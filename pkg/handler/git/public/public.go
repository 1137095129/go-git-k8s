package public

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/wang1137095129/go-git-k8s/config"
	"github.com/wang1137095129/go-git-k8s/utils"
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
		utils.CheckError(err)

		commitIter, err := repository.Log(&git.LogOptions{From: head.Hash()})
		utils.CheckError(err)

		t := time.Time{}
		err = commitIter.ForEach(func(commit *object.Commit) error {
			if t.Before(commit.Committer.When) {
				t = commit.Committer.When
			}
			return nil
		})
		utils.CheckError(err)
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
	utils.CheckError(err)

	revision, err := repository.ResolveRevision(plumbing.Revision(fmt.Sprintf("%s/%s", c.Git.Remote, c.Git.Branch)))
	utils.CheckError(err)

	commit, err := repository.CommitObject(*revision)
	utils.CheckError(err)

	if p.lastPullTime.Before(commit.Committer.When) {
		worktree, err := repository.Worktree()
		utils.CheckError(err)
		err = worktree.Pull(&git.PullOptions{
			RemoteName:    c.Git.Remote,
			ReferenceName: plumbing.NewBranchReferenceName(c.Git.Branch),
		})
		utils.CheckError(err)
		*p.lastPullTime = commit.Committer.When
	}

	return repository, nil
}
