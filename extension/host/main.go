package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

type FunctionType int

const (
	Unknown FunctionType = iota
	Initialize
	Apply
)

type Csv struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

type Request struct {
	FunctionType     FunctionType `json:"functionType"`
	RepositoryUrl    string       `json:"repositoryUrl"`
	TargetBranchName string       `json:"targetBranchName"`
	Username         string       `json:"username"`
	Email            string       `json:"email"`
	Csvs             []Csv        `json:"csvs"`
}

type Response struct {
	ErrorMessage string `json:"errorMessage"`
}

var filesystem billy.Filesystem
var repository *git.Repository

func main() {
	for {
		msg := readMessage()

		var req Request
		err := json.Unmarshal(msg, &req)
		if err != nil {
			sendMessage([]byte(fmt.Sprintf(`{"errorMessage":"%s"}`, err)))
			continue
		}

		switch req.FunctionType {
		case Initialize:
			err = initialize(req.RepositoryUrl)
			if err == nil {
				sendMessage([]byte("{}"))
				continue
			}
		case Apply:
			err = apply(req.TargetBranchName, req.Username, req.Email, req.Csvs)
			if err == nil {
				sendMessage([]byte("{}"))
				continue
			}
		default:
			sendMessage([]byte(`{"errorMessage":"Unknown functionType."}`))
			continue
		}

		if err != nil {
			sendMessage([]byte(fmt.Sprintf(`{"errorMessage":"%s"}`, err)))
		}
	}
}

func initialize(repositoryUrl string) error {
	filesystem = memfs.New()
	var err error
	repository, err = git.Clone(memory.NewStorage(), filesystem, &git.CloneOptions{
		URL: repositoryUrl,
	})
	return err
}

func apply(targetBranchName string, username string, email string, csvs []Csv) error {
	branchs, err := repository.Branches()
	if err != nil {
		return err
	}
	branchs.ForEach(func(branch *plumbing.Reference) error {
		if branch.Name().String() == "refs/heads/"+targetBranchName {
			err = repository.DeleteBranch(targetBranchName)
			return err
		}
		return nil
	})

	repository.Fetch(&git.FetchOptions{})

	w, err := repository.Worktree()
	if err != nil {
		return err
	}
	remoteRef, err := repository.Reference(plumbing.ReferenceName("refs/remotes/origin/"+targetBranchName), true)
	var refName plumbing.ReferenceName
	if err == nil {
		newRef := plumbing.NewHashReference(plumbing.ReferenceName(targetBranchName), remoteRef.Hash())
		repository.Storer.SetReference(newRef)
		refName = newRef.Name()
		err = w.Checkout(&git.CheckoutOptions{
			Branch: refName,
			Create: false,
			Force:  true,
		})
		if err != nil {
			return err
		}
	} else {
		refName = plumbing.ReferenceName(targetBranchName)
		err = w.Checkout(&git.CheckoutOptions{
			Branch: refName,
			Create: true,
			Force:  true,
		})
		if err != nil {
			return err
		}
	}

	var addedCsvPaths []string
	for _, csv := range csvs {
		csvPath := csv.Path + ".csv"
		file, err := filesystem.OpenFile(csvPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		file.Write([]byte(csv.Value))

		_, err = w.Add(csvPath)
		if err != nil {
			continue
		}

		addedCsvPaths = append(addedCsvPaths, csvPath)
	}

	if len(addedCsvPaths) == 0 {
		return fmt.Errorf("%s", "Not changed.")
	}

	hash, err := w.Commit("Update "+strings.Join(addedCsvPaths, ", "), &git.CommitOptions{
		Author: &object.Signature{
			Name:  username,
			Email: email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	repository.Storer.SetReference(plumbing.NewReferenceFromStrings(targetBranchName, hash.String()))

	remote, err := repository.Remote("origin")
	if err != nil {
		return err
	}
	err = remote.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec(refName + ":" + plumbing.ReferenceName("refs/heads/"+targetBranchName)),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func readMessage() []byte {
	r := bufio.NewReader(os.Stdin)
	length := make([]byte, 4)
	r.Read(length)
	msg := make([]byte, readMessageLength(length))
	io.ReadFull(r, msg)
	return msg
}

func readMessageLength(msg []byte) int {
	var length uint32
	buf := bytes.NewBuffer(msg)
	binary.Read(buf, binary.LittleEndian, &length)
	return int(length)
}

func sendMessage(msg []byte) {
	var buf bytes.Buffer
	binary.Write(os.Stdout, binary.LittleEndian, uint32(len(msg)))
	buf.Write(msg)
	buf.WriteTo(os.Stdout)
}