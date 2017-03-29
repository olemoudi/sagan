package main

import (
	"gopkg.in/libgit2/git2go.v24"
	"strconv"
)

type Project struct {
	name string
	uri  string
	lock chan interface{}
	repo *git.Repository
	//repo *CodeRepo
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
	debug("updating project", p.name)
	remoteNames, err := p.repo.Remotes.List()
	if err != nil {
		info("Error retrieving remotes for project", p.name)
		return
	}

	for _, remoteName := range remoteNames {
		debug("updating", remoteName)
		remote, err := p.repo.Remotes.Lookup(remoteName)
		if err != nil {
			panic("error lookingup remote")
		}

		refspecs, err := remote.FetchRefspecs()
		if err != nil {
			panic("error retrieving refspecs for remote" + remoteName)
		}
		remote.Fetch(refspecs, &git.FetchOptions{}, "")
		debug(remoteName, "remote updated")
	}
	debug(p.name, "updated")
	debug("last commits to master")
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

	//oid := ref.Target()
	//_, err = p.repo.LookupCommit(oid)

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
	p := Project{name, uri, make(chan interface{}, 1), nil}
	return &p
}
