package main

import (
	"errors"
	"os"
	"time"

	"github.com/fsouza/go-dockerclient"
	"archive/tar"
	"io"
	"strings"
	"bytes"
	"io/ioutil"
	"fmt"
	"github.com/offers/dope/out"
)

func newDockerClient() (*docker.Client, error) {
	endpoint := "unix:///var/run/docker.sock"
	return docker.NewClient(endpoint)
}

func dockerPull(repo string, tag string) error {
	client, err := newDockerClient()
	if err != nil {
		return err
	}

	opts := docker.PullImageOptions{
		Repository:        repo,
		Tag:               tag,
		InactivityTimeout: time.Duration(30) * time.Second,
	}

	out.Infof("Pulling image %s:%s...\n", repo, tag)
	return client.PullImage(opts, docker.AuthConfiguration{})
}

func dockerRmi(image string) error {
	client, err := newDockerClient()
	if err != nil {
		return err
	}

	return client.RemoveImage(image)
}

func dockerGetDopeFile(repo string, tag string) ([]byte, error) {
	exportFile, err := ioutil.TempFile("", "dope-docker-export")
	if err != nil {
		return []byte{}, err
	}
	defer exportFile.Close()
	defer os.Remove(exportFile.Name())


	client, err := newDockerClient()
	if err != nil {
		return []byte{}, err
	}

	image := fmt.Sprintf("%s:%s", repo, tag)
	imageData, err := client.InspectImage(image)
	if err != nil {
		return []byte{}, err
	}

	exportOpts := docker.ExportImageOptions {
		Name: imageData.ID,
		OutputStream: exportFile,
	}
	out.Info("Searching docker image for .dope.json...")
	err = client.ExportImage(exportOpts)
	if err != nil {
		return []byte{}, err
	}

	// Rewind to beginning of .tar for reading
	exportFile.Seek(0, 0)
	if err != nil {
		return []byte{}, err
	}

	tr := tar.NewReader(exportFile)
	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasSuffix(hdr.Name, "layer.tar")  {
			buf := make([]byte, hdr.Size)
			_, err := tr.Read(buf)
			if err != nil {
				return []byte{}, err
			}

			data, err := dockerGetFileFromLayer(".dope.json", buf)
			if nil == err {
				return data, nil
			}
		}
	}
	return []byte{}, errors.New(".dope.json not found in tar archive")
}

func dockerGetFileFromLayer(filename string, layer []byte) ([]byte, error){
	tr := tar.NewReader(bytes.NewReader(layer))
	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasSuffix(hdr.Name, filename)  {
			buf := make([]byte, hdr.Size)
			_, err := tr.Read(buf)
			if err != nil {
				return []byte{}, err
			}

			return buf, nil
		}
	}
	return []byte{}, errors.New(filename + " not found in layer")
}
