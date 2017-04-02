package main

import (
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
		case _ = <-pm.add:
			debug("received new project")

		case <-time.After(time.Second * 10):
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
		p.Update()
	}
	debug("projects update complete")

}
