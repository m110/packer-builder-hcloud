package hcloud

import "fmt"

type Artifact struct {
	imageName string
	imageID   string
}

func (Artifact) BuilderId() string {
	return BuilderID
}

func (Artifact) Files() []string {
	return nil
}

func (a Artifact) Id() string {
	return a.imageID
}

func (a Artifact) String() string {
	return fmt.Sprintf("%s (%s)", a.imageID, a.imageName)
}

func (Artifact) State(name string) interface{} {
	return nil
}

func (a Artifact) Destroy() error {
	return nil
}
