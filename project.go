package main

import (
	"gopkg.in/libgit2/git2go.v24"
	"path/filepath"
)

type Project struct {
	name       string
	uri        string
	lock       chan interface{}
	repo       *git.Repository
	remoteName string
}

func (p Project) Lock() {
	p.lock <- struct{}{}
}
func (p Project) Unlock() {
	<-p.lock
}

func (p Project) Update() {
	p.Lock()
	defer p.Unlock()
	debug("updating", p.name)
	debug("updating remote", p.remoteName)
	remote, err := p.repo.Remotes.Lookup(p.remoteName)
	cep(err, true, "Error looking up remote")
	refspecs, err := remote.FetchRefspecs()
	cep(err, true, "Error retrieving refspecs for remote", p.remoteName)
	remote.Fetch(refspecs, &git.FetchOptions{}, "")
	debug(p.name, "updated")
	/*
		debug("last 10 commits to master")
		revwalk, err := p.repo.Walk()
		revwalk.PushHead()
		counter := 1
		err = revwalk.Iterate(func(c *git.Commit) bool {
			debug("Commit id", c.TreeId().String())
			debug("Msg", c.Message())

			debug("=========")
			counter = counter + 1
			if counter >= 10 {
				return false
			} else {
				return true
			}

		})

	*/

	/*
		//ref, err := p.repo.Head()
		refiter, err := p.repo.NewReferenceIterator()
		ce(err, "error extracting Reference Iterator from repo")
		//ref, err := refiter.Next()
		ref, err := p.repo.Head()
		if ce(err, "Error getting first commit") {
			return
		}
		counter := 1

		oid := ref.Target()
		commit, err := p.repo.LookupCommit(oid)
		if ce(err, "Error looking up commit") {
		}
		debug("Commit #", strconv.Itoa(counter))
		//debug(hex.EncodeToString([20]byte(commit.TreeId())))
		debug(commit.TreeId().String())
		//debug(string(commit.TreeId()[:20]))
		counter = counter + 1
		ref, err = refiter.Next()
		if ce(err, "Error getting next commit") {

		}
		/*
			for err == nil && counter < 3 {
				oid := ref.Target()
				commit, err := p.repo.LookupCommit(oid)
				if ce(err, "Error looking up commit") {
					counter = counter + 1
					continue
				}
				debug("Commit #", strconv.Itoa(counter))
				//debug(hex.EncodeToString([20]byte(commit.TreeId())))
				debug(commit.TreeId().String())
				//debug(string(commit.TreeId()[:20]))
				counter = counter + 1
				ref, err = refiter.Next()
				if ce(err, "Error getting next commit") {
					break
				}
			}
	*/

	/*
		//oid := ref.Target()
		//_, err = p.repo.LookupCommit(oid)
	*/

}
func (p Project) ListAllBranches() []string {
	return p.ListBranches(git.BranchAll)
}
func (p Project) ListLocalBranches() []string {
	return p.ListBranches(git.BranchLocal)
}

func (p Project) ListBranches(flags git.BranchType) []string {
	if p.repo == nil {
		return nil
	}

	iter, err := p.repo.NewBranchIterator(git.BranchAll)
	if err != nil {
		return nil
	}

	names := make([]string, 0)
	iter.ForEach(func(b *git.Branch, btype git.BranchType) error {
		name, err := b.Name()
		if err != nil {
			debug("error while listing branchname")
			return err
		}
		names = append(names, name)
		return nil
	})

	return names

}

func makeProject(name, uri string) *Project {
	debug("creating new project", name)
	p := Project{name, uri, make(chan interface{}, 1), nil, "origin"}
	p.Lock()
	defer p.Unlock()

	pm.projects[p.name] = &p
	path := reposPath + string(filepath.Separator) + p.name
	dirExists, err := exists(path)
	if err != nil {
		debug("dir existence check returned error", err.Error())
	}
	var repo *git.Repository
	if dirExists {
		debug("target dir for repo already exists, attempting to open local repository")
		repo, err = git.OpenRepository(path)
		cep(err, true, "")
	} else {
		debug("Cloning repository from", p.uri, "into", path)
		options := &git.CloneOptions{}
		repo, err = git.Clone(p.uri, path, options)
		cep(err, true, "")
	}
	p.repo = repo
	debug(p.name, "project created")
	return &p
}
