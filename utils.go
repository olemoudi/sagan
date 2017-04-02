package main

import (
	"gopkg.in/libgit2/git2go.v24"
	"os"
)

func ce(err error, msg ...string) bool {
	return cep(err, false, msg...)
}

func cep(err error, p bool, msg ...string) bool {
	if err != nil {
		if p {
			panic(err)
		} else {
			debug(msg...)
			debug(err.Error())
			return true
		}
	}
	return false
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func Pull(repo *git.Repository) error {
	branch, err := repo.Branch()
	if err != nil {
		return err
	}

	// Get the name
	name, err := branch.Name()
	if err != nil {
		return err
	}

	// Locate remote
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return err
	}

	// Fetch changes from remote
	if err := remote.Fetch([]string{}, nil, ""); err != nil {
		return err
	}

	// Get remote master
	remoteBranch, err := repo.References.Lookup("refs/remotes/origin/" + name)
	if err != nil {
		return err
	}

	remoteBranchID := remoteBranch.Target()
	// Get annotated commit
	annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)
	if err != nil {
		return err
	}

	// Do the merge analysis
	mergeHeads := make([]*git.AnnotatedCommit, 1)
	mergeHeads[0] = annotatedCommit
	analysis, _, err := repo.MergeAnalysis(mergeHeads)
	if err != nil {
		return err
	}

	// Get repo head
	head, err := repo.Head()
	if err != nil {
		return err
	}

	if analysis&git.MergeAnalysisUpToDate != 0 {
		return nil
	} else if analysis&git.MergeAnalysisNormal != 0 {
		// Just merge changes
		if err := repo.Merge([]*git.AnnotatedCommit{annotatedCommit}, nil, nil); err != nil {
			return err
		}
		// Check for conflicts
		index, err := repo.Index()
		if err != nil {
			return err
		}

		if index.HasConflicts() {
			return errors.New("Conflicts encountered. Please resolve them.")
		}

		// Make the merge commit
		sig, err := repo.DefaultSignature()
		if err != nil {
			return err
		}

		// Get Write Tree
		treeId, err := index.WriteTree()
		if err != nil {
			return err
		}

		tree, err := repo.LookupTree(treeId)
		if err != nil {
			return err
		}

		localCommit, err := repo.LookupCommit(head.Target())
		if err != nil {
			return err
		}

		remoteCommit, err := repo.LookupCommit(remoteBranchID)
		if err != nil {
			return err
		}

		repo.CreateCommit("HEAD", sig, sig, "", tree, localCommit, remoteCommit)

		// Clean up
		repo.StateCleanup()
	} else if analysis&git.MergeAnalysisFastForward != 0 {
		// Fast-forward changes
		// Get remote tree
		remoteTree, err := repo.LookupTree(remoteBranchID)
		if err != nil {
			return err
		}

		// Checkout
		if err := repo.CheckoutTree(remoteTree, nil); err != nil {
			return err
		}

		branchRef, err := repo.References.Lookup("refs/heads/" + name)
		if err != nil {
			return err
		}

		// Point branch to the object
		branchRef.SetTarget(remoteBranchID, "")
		if _, err := head.SetTarget(remoteBranchID, ""); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Unexpected merge analysis result %d", analysis)
	}

	return nil
}
