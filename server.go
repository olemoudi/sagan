package main

import (
	// "fmt"
	// "io"
	"gopkg.in/gin-gonic/gin.v1"
	// DO NOT UNCOMMENT GIT
	//"gopkg.in/libgit2/git2go.v24"
	//"io"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/olemoudi/sagan/proto"
	"log"
	"net/http"
)

var (
	r *gin.Engine
)

func loadRoutes() {

	r.Static("/s", "./static")
	r.GET("/projects", ListProjects)
	r.GET("/project/:name", ListProject)
	r.GET("/projectpb/:name", ListProjectProto)

}

func ListProjects(c *gin.Context) {
	for _, v := range pm.ListProjectNames() {
		c.JSON(http.StatusOK, gin.H{"project_name": v})
	}
}

func ListProject(c *gin.Context) {
	p := pm.projects[c.Param("name")]
	/*branch, err := p.repo.LookupBranch("master", git.BranchAll)
	if err != nil {
		c.String(http.StatusInternalServerError, "error with branch")
		return
	}
	name, err := branch.Name()
	if err != nil {
		c.String(http.StatusInternalServerError, "error with branch name")
		return
	}
	*/
	c.JSON(http.StatusOK, gin.H{"branches": p.ListLocalBranches()})
}

func ListProjectProto(c *gin.Context) {
	p := pm.projects[c.Param("name")]
	pb := &pb.Project{
		Id:   1,
		Name: p.name,
		Uri:  p.uri,
		Branches: []*pb.Project_Branch{
			{Name: "master"},
		},
	}
	out, err := proto.Marshal(pb)
	if err != nil {
		log.Fatalln("Failed to encode Project:", err)
	}
	c.String(http.StatusOK, string(out))
}

func webServer() {
	defer wg.Done()

	r = gin.Default()

	loadRoutes()

	//http.HandleFunc("/projects", ListProjects)
	var server http.Server
	server.Addr = ":8443"
	info(fmt.Sprintf("Starting web server %s...", server.Addr))
	server.Handler = r
	err := server.ListenAndServeTLS("tls/server.crt", "tls/server.key")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	/*
		go func() {
			select {
			case <-exiting:
				server.Close()
			}
		}()
	*/
}
