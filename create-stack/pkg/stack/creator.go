package stack

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"time"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . PackageFinder
type PackageFinder interface {
	GetPackagesList(image string) ([]string, error)
	GetPackageMetadata(imageName string) (string, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . MixinsGenerator
type MixinsGenerator interface {
	GetMixins(buildPackages, runPackages []string) ([]string, []string)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ImageClient
type ImageClient interface {
	Build(tag, dockerfilePath string, buildArgs ...string) error
	Push(tag string) (string, error)
	Pull(tag string, keychain authn.Keychain) (v1.Image, error)
	SetLabel(tag, key, value string) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Stack
type Stack interface {
	GetName() string
	GetSources() string
	GetBuildPackages() string
	GetRunPackages() string
	GetBaseBuildDockerfilePath() string
	GetBaseRunDockerfilePath() string
	GetCNBBuildDockerfilePath() string
	GetCNBRunDockerfilePath() string
	GetBuildDescription() string
	GetRunDescription() string
}

type Creator struct {
	PackageFinder   PackageFinder
	MixinsGenerator MixinsGenerator
	ImageClient     ImageClient
}

func (c Creator) CreateBionicStack(stack Stack, buildBaseTag, runBaseTag string, publish bool) error {
	_, err := c.ImageClient.Pull("ubuntu:bionic", authn.DefaultKeychain)
	if err != nil {
		return fmt.Errorf("error pulling bionic image: %w", err)
	}

	buildBaseRef, runBaseRef, err := c.buildBaseStackImages(stack, buildBaseTag, runBaseTag, publish)
	if err != nil {
		return fmt.Errorf("error building base stack images: %w", err)
	}

	err = c.buildCnbStackImages(stack, buildBaseTag, runBaseTag, buildBaseRef, runBaseRef, publish)
	if err != nil {
		return fmt.Errorf("error building CNB stack images: %w", err)
	}

	return nil
}

func (c Creator) buildBaseStackImages(stack Stack, buildBaseTag, runBaseTag string, publish bool) (string, string, error) {
	buildBaseRef, err := c.buildBaseImage(buildBaseTag, stack.GetSources(), stack.GetBuildPackages(),
		stack.GetBaseBuildDockerfilePath(), publish)
	if err != nil {
		return "", "", fmt.Errorf("error building base build image: %w", err)
	}

	runBaseRef, err := c.buildBaseImage(runBaseTag, stack.GetSources(), stack.GetRunPackages(),
		stack.GetBaseRunDockerfilePath(), publish)
	if err != nil {
		return "", "", fmt.Errorf("error building base run image: %w", err)
	}

	return buildBaseRef, runBaseRef, nil
}

func (c Creator) buildCnbStackImages(stack Stack, buildBaseTag, runBaseTag, buildBaseRef, runBaseRef string, publish bool) error {

	buildMixins, runMixins, err := c.getMixins(buildBaseTag, runBaseTag)
	if err != nil {
		return fmt.Errorf("error getting mixins: %w", err)
	}

	releaseDate := time.Now()

	err = c.buildCnbImage(stack.GetCNBBuildDockerfilePath(), buildBaseTag, buildBaseRef, stack.GetBuildDescription(), releaseDate, buildMixins, publish)
	if err != nil {
		return fmt.Errorf("error building cnb build image: %w", err)
	}

	err = c.buildCnbImage(stack.GetCNBRunDockerfilePath(), runBaseTag, runBaseRef, stack.GetRunDescription(), releaseDate, runMixins, publish)
	if err != nil {
		return fmt.Errorf("error building cnb run image: %w", err)
	}

	return nil
}

func (c Creator) buildBaseImage(tag, sources, packages, dockerfilePath string, publish bool) (string, error) {
	err := c.ImageClient.Build(tag, dockerfilePath,
		fmt.Sprintf("sources=%s", sources),
		fmt.Sprintf("packages=%s", packages))
	if err != nil {
		return "", fmt.Errorf("failed to build base image: %w", err)
	}

	if publish {
		imageRef, err := c.ImageClient.Push(tag)
		if err != nil {
			return "", fmt.Errorf("failed to push tag %s: %w", tag, err)
		}
		return imageRef, nil
	}

	return "", nil
}

func (c Creator) buildCnbImage(dockerfilePath, baseTag, baseRef, description string, releaseDate time.Time, mixinsList []string, publish bool) error {

	mixinsJson, err := json.Marshal(mixinsList)
	if err != nil {
		return fmt.Errorf("failed to marshal mixin array: %w", err)
	}

	cnbTag := fmt.Sprintf("%s-cnb", baseTag)

	type Metadata struct {
		BaseImage string `json:"base-image,omitempty"`
	}

	var metadata Metadata
	if publish {
		metadata = Metadata{BaseImage: baseRef}
	}

	metadataMap, err := json.Marshal(metadata)

	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = c.ImageClient.Build(cnbTag, dockerfilePath,
		fmt.Sprintf(`base_image=%s`, baseTag),
		fmt.Sprintf("description=%s", description),
		fmt.Sprintf("mixins=%s", string(mixinsJson)),
		fmt.Sprintf("released=%s", releaseDate.Format(time.RFC3339)),
		fmt.Sprintf(`metadata=%s`, string(metadataMap)))

	if err != nil {
		return fmt.Errorf("failed to build cnb image: %w", err)
	}

	packageMetadata, err := c.PackageFinder.GetPackageMetadata(baseTag)
	if err != nil {
		return fmt.Errorf("failed to get package metadata: %w", err)
	}

	err = c.ImageClient.SetLabel(cnbTag, "io.paketo.stack.packages", packageMetadata)
	if err != nil {
		return fmt.Errorf("failed to set label: %w", err)
	}

	if publish {
		_, err := c.ImageClient.Push(cnbTag)
		if err != nil {
			return fmt.Errorf("failed to push tag %s: %w", cnbTag, err)
		}
		return nil
	}

	return nil
}

func (c Creator) getMixins(buildBaseTag, runBaseTag string) ([]string, []string, error) {
	buildBasePackageList, err := c.PackageFinder.GetPackagesList(buildBaseTag)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting packages from base build image: %w", err)
	}

	runBasePackageList, err := c.PackageFinder.GetPackagesList(runBaseTag)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting packages from base run image: %w", err)
	}

	buildMixins, runMixins := c.MixinsGenerator.GetMixins(buildBasePackageList, runBasePackageList)

	return buildMixins, runMixins, nil

}
