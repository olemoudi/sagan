package main

import (
	"gopkg.in/libgit2/git2go.v24"
	"path/filepath"
	"time"
)

type Project struct {
	name string
	uri  string
	lock chan interface{}
	repo *git.Repository
}

func (p Project) Lock() {
	p.lock <- struct{}{}
}
func (p Project) Unlock() {
	<-p.lock
}

func (p Project) listBranches() []string {
	if p.repo == nil {
		return nil
	}

	iter, err := p.repo.NewBranchIterator(git.BranchAll)
	if err != nil {
		return nil
	}

	names := make([]string, 1)
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

func makeProject(name, uri string) Project {
	p := Project{name, uri, make(chan interface{}, 1), nil}
	return p
}

type ProjectManager struct {
	add      chan *Project
	projects map[string]*Project
}

func (pm ProjectManager) Run() {
	defer wg.Done()
	debug("running pm")

loop:
	for {
		select {
		case <-exiting:
			break loop
		case newproject := <-pm.add:
			debug("received new project")
			pm.projects[newproject.name] = newproject
			go pm.CreateProject(newproject)
		case <-time.After(time.Second * 30):
			go pm.UpdateProjects()
		}

	}
}

func (pm ProjectManager) ListProjectNames() []string {
	projects := make([]string, 1)
	for k, _ := range pm.projects {
		projects = append(projects, k)
	}
	return projects
}

func getProject(name string) *Project {
	return pm.projects[name]
}

func (pm ProjectManager) UpdateProjects() {
	debug("updating projects...")
	debug("projects update complete")
}

func (pm ProjectManager) CreateProject(p *Project) {
	p.Lock()
	defer p.Unlock()

	path := reposPath + string(filepath.Separator) + p.name
	dirExists, err := exists(path)
	if err != nil {
		debug(err.Error())
	}
	var repo *git.Repository
	if dirExists {
		debug("target dir for repo already exists, attempting to open local repository")
		repo, err = git.OpenRepository(path)
		if err != nil {
			panic(err)
		}
	} else {
		debug("Cloning repository from", p.uri, "into", path)
		options := &git.CloneOptions{}
		repo, err = git.Clone(p.uri, path, options)
		if err != nil {
			panic(err)
		}
	}
	p.repo = repo
	debug("project created")
}
