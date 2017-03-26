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
			debug("go create")
			go pm.CreateProject(newproject)
		case <-time.After(time.Second * 30):
			go pm.UpdateProjects()
		}

	}
}

func (pm ProjectManager) UpdateProjects() {
	debug("updating projects...")
	debug("projects update complete")
}

func (pm ProjectManager) CreateProject(p *Project) {
	p.lock <- struct{}{}

	path := reposPath + string(filepath.Separator) + p.name
	e, err := exists(path)
	if err != nil {
		debug(err.Error())
	}
	var repo *git.Repository
	if e {
		debug("target dir for repo already exists, attempting to open local repository")
		repo, err = git.OpenRepository(path)
		if err != nil {
			panic(err)
		}
	} else {
		debug("Cloning repository from", p.uri, "into", path)
		options := &git.CloneOptions{}
		repo, err = git.Clone(p.uri, reposPath+string(filepath.Separator)+p.name, options)
		if err != nil {
			panic(err)
		}
	}
	p.repo = repo
	debug("project created")
	<-p.lock
}

//func (p Project) listBranches()
