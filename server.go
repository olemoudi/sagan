package main

import (
	// "fmt"
	// "io"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/libgit2/git2go.v24"
	//"io"
	"github.com/golang/protobuf/proto"
	"log"
	"net/http"
)

var (
	r *gin.Engine
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

/*
func ListProjects(w http.ResponseWriter, req *http.Request) {
	for _, v := range pm.ListProjectNames() {
		io.WriteString(w, v+"<br>")
	}

}
*/

func ListProjects(c *gin.Context) {
	for _, v := range pm.ListProjectNames() {
		c.JSON(http.StatusOK, gin.H{"project_name": v})
	}
}

func ListProject(c *gin.Context) {
	p := pm.projects[c.Param("name")]
	branch, err := p.repo.LookupBranch("master", git.BranchAll)
	if err != nil {
		c.String(http.StatusInternalServerError, "error with branch")
		return
	}
	name, err := branch.Name()
	if err != nil {
		c.String(http.StatusInternalServerError, "error with branch name")
		return
	}
	c.JSON(http.StatusOK, gin.H{"branch_name": name})
}

func ListProjectProto(c *gin.Context) {
	p := pb.PbProject{}

}

func loadRoutes() {

	r.Static("/s", "./static")
	r.GET("/projects", ListProjects)
	r.GET("/project/:name", ListProject)

}

func webServer() {
	defer wg.Done()

	r = gin.Default()

	loadRoutes()

	//http.HandleFunc("/projects", ListProjects)
	info("Starting web server...")
	var server http.Server
	server.Addr = ":8443"
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
