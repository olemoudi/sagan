package main

import (
	"gopkg.in/libgit2/git2go.v24"
	"path/filepath"
	"time"
)

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
			go pm.CreateProject(newproject)
		case <-time.After(time.Second * 30):
			go pm.UpdateProjects()
		}

	}
}

func (pm ProjectManager) ListProjectNames() []string {
	projects := make([]string, 0)
	for k, _ := range pm.projects {
		projects = append(projects, k)
	}
	return projects
}

func getProject(name string) *Project {
	return pm.projects[name]
}

//TODO: update on bursts
func (pm ProjectManager) UpdateProjects() {
	debug("updating projects...")
	for _, p := range pm.projects {
		go p.Update()
	}
	debug("projects update complete")

}

func (pm ProjectManager) CreateProject(p *Project) {
	debug("adding new project", p.name)
	p.Lock()
	defer p.Unlock()

	pm.projects[p.name] = p

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
