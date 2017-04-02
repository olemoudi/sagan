package main

import (
	"fmt"
	"gopkg.in/libgit2/git2go.v24"
	"os"
	"strconv"
)

type Project struct {
	name string
	uri  string
	lock chan interface{}
	repo *git.Repository
	//repo *CodeRepo
}

func makeProject(name, uri string) *Project {
	p := Project{name, uri, make(chan interface{}, 1), nil}
	return &p
}

func ce(err error) {
	if err != nil {
		panic(err)
	}
}

/*
func (p Project) Update() {
	p.Lock()
	defer p.Unlock()
	debug("updating", p.name)
	//remoteNames, err := p.repo.Remotes.List()
	debug("updating remote", p.remoteName)
	remote, err := p.repo.Remotes.Lookup(p.remoteName)
	cep(err, true, "Error looking up remote")
	refspecs, err := remote.FetchRefspecs()
	cep(err, true, "Error retrieving refspecs for remote", p.remoteName)
	remote.Fetch(refspecs, &git.FetchOptions{}, "")
	debug(p.name, "updated")
}
*/

var counter = 1

var c1, c2 *git.Commit

func main() {

	//p := makeProject("sqlmap", "git://github.com/sqlmapproject/sqlmap")
	p := makeProject("freeCodeCamp", "git://github.com/freeCodeCamp/freeCodeCamp")

	// create repo from uri

	repo, err := git.Clone(p.uri, p.name, &git.CloneOptions{})
	//ce(err)
	fmt.Println(repo.Remotes.List())
	os.Exit(1)

	// open repo from dir

	//repo, err := git.OpenRepository(p.name)
	ce(err)

	p.repo = repo

	// get head

	/*headref, err := p.repo.Head()
	ce(err)
	headoid := headref.Target()
	headcommit, err := p.repo.LookupCommit(headoid)
	ce(err)

	fmt.Println("head tree ID", headcommit.TreeId().String())
	fmt.Println("head commit message", headcommit.Message())
	*/

	revwalk, err := p.repo.Walk()
	ce(err)
	revwalk.Sorting(git.SortTime | git.SortTopological)
	err = revwalk.PushHead()
	ce(err)
	err = revwalk.HideGlob("tags/*")
	ce(err)
	err = revwalk.Iterate(func(commit *git.Commit) bool {
		c1 = c2
		c2 = commit
		if c1 != nil {
			fmt.Println("========")
			fmt.Println("Commit " + c1.TreeId().String())
			fmt.Println("Parent Count = " + strconv.Itoa(int(c1.ParentCount())))
			for i := 0; i < int(c1.ParentCount()); i++ {
				fmt.Println("Parent #" + strconv.Itoa(i) + " - " + c1.Parent(uint(i)).TreeId().String())

				fmt.Println("MSG " + c1.Message())
				tree1, err := c1.Tree()
				ce(err)
				parent := c1.Parent(uint(i))
				if parent == nil {
					continue
				}
				tree2, err := parent.Tree()
				ce(err)
				opts, err := git.DefaultDiffOptions()
				ce(err)
				diff, err := p.repo.DiffTreeToTree(tree1, tree2, &opts)

				numDiffs := 0
				numAdded := 0
				numDeleted := 0
				err = diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
					numDiffs++
					fmt.Println("Diff #" + strconv.Itoa(numDiffs))
					fmt.Println("DiffFile OldFile", file.OldFile)
					fmt.Println("DiffFile NewFile", file.NewFile)
					fmt.Println()
					fmt.Println("/* Start Snippet */")
					fmt.Println()

					switch file.Status {
					case git.DeltaAdded:
						numAdded++
						fmt.Println("delta added")
					case git.DeltaDeleted:
						numDeleted++
						fmt.Println("delta deleted")
					}
					return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
						return func(line git.DiffLine) error {
							switch line.Origin {
							case git.DiffLineContext:
								fmt.Println(line.Content)
							case git.DiffLineAddition:
								fmt.Println("+", line.Content)
							case git.DiffLineDeletion:
								fmt.Println("-", line.Content)
							case git.DiffLineContextEOFNL:
								fmt.Println("ContextEOFLN", line.Content)
							case git.DiffLineAddEOFNL:
								fmt.Println("AddEOFLN", line.Content)
							case git.DiffLineDelEOFNL:
								fmt.Println("DelEOFLN", line.Content)
							case git.DiffLineFileHdr:
								fmt.Println("FileHdr", line.Content)
							case git.DiffLineHunkHdr:
								fmt.Println("HunkHdr", line.Content)
							}
							//fmt.Println("Diffline content", line.Content)

							return nil
						}, nil
					}, nil
				}, git.DiffDetailLines)
				/*

					numdeltas, err := diff.NumDeltas()
					ce(err)
					err = diff.ForEach(func(delta DiffDelta, f float64) (git.DiffForEachHunkCallback, error) {
						return nil, nil
					}, git.DiffDetailFiles)
					/*
					/*
						for i := 0; i < numdeltas; i++ {
							delta, err := diff.GetDelta(i)

							ce(err)
						}
				*/
				fmt.Println()
				fmt.Println("/* End Snippet */")
				fmt.Println()
			}
		}
		counter = counter + 1
		if counter > 4000 {
			return false
		} else {
			return true
		}

	})
	/*

		// get ref iterator

		refiter, err := p.repo.NewReferenceIterator()
		ce(err)

		ref, err := refiter.Next()
		oid := ref.Target()
		refcommit, err := p.repo.LookupCommit(oid)
		ce(err)
		fmt.Println("first iter commit", refcommit.TreeId().String())
		fmt.Println("first iter commit msg", refcommit.Message())

		// second ref on iter
		ref, err = refiter.Next()
		oid = ref.Target()
		refcommit, err = p.repo.LookupCommit(oid)
		ce(err)
		fmt.Println("second iter commit", refcommit.TreeId().String())
		fmt.Println("second iter commit msg", refcommit.Message())
	*/

}
