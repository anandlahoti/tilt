// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cli

import (
	context "context"
	build "github.com/windmilleng/tilt/internal/build"
	docker "github.com/windmilleng/tilt/internal/docker"
	engine "github.com/windmilleng/tilt/internal/engine"
	k8s "github.com/windmilleng/tilt/internal/k8s"
	model "github.com/windmilleng/tilt/internal/model"
)

// Injectors from wire.go:

func wireManifestCreator(ctx context.Context, browser engine.BrowserMode) (model.ManifestCreator, error) {
	env, err := k8s.DetectEnv()
	if err != nil {
		return nil, err
	}
	config, err := k8s.ProvideRESTConfig()
	if err != nil {
		return nil, err
	}
	coreV1Interface, err := k8s.ProvideRESTClient(config)
	if err != nil {
		return nil, err
	}
	portForwarder := k8s.ProvidePortForwarder()
	k8sClient := k8s.NewK8sClient(ctx, env, coreV1Interface, config, portForwarder)
	sidecarSyncletManager := engine.NewSidecarSyncletManager(k8sClient)
	syncletBuildAndDeployer := engine.NewSyncletBuildAndDeployer(k8sClient, sidecarSyncletManager)
	dockerCli, err := docker.DefaultDockerClient(ctx, env)
	if err != nil {
		return nil, err
	}
	containerUpdater := build.NewContainerUpdater(dockerCli)
	containerResolver := build.NewContainerResolver(dockerCli)
	analytics, err := provideAnalytics()
	if err != nil {
		return nil, err
	}
	localContainerBuildAndDeployer := engine.NewLocalContainerBuildAndDeployer(containerUpdater, containerResolver, env, k8sClient, analytics)
	console := build.DefaultConsole()
	writer := build.DefaultOut()
	labels := _wireLabelsValue
	dockerImageBuilder := build.NewDockerImageBuilder(dockerCli, console, writer, labels)
	imageBuilder := build.DefaultImageBuilder(dockerImageBuilder)
	imageBuildAndDeployer := engine.NewImageBuildAndDeployer(imageBuilder, k8sClient, env, analytics)
	buildOrder := engine.DefaultBuildOrder(syncletBuildAndDeployer, localContainerBuildAndDeployer, imageBuildAndDeployer, env)
	fallbackTester := engine.DefaultShouldFallBack()
	compositeBuildAndDeployer := engine.NewCompositeBuildAndDeployer(buildOrder, fallbackTester)
	imageReaper := build.NewImageReaper(dockerCli)
	upper := engine.NewUpper(ctx, compositeBuildAndDeployer, k8sClient, browser, imageReaper)
	return upper, nil
}

var (
	_wireLabelsValue = build.Labels{}
)
