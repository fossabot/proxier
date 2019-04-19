package e2e

import (
	"github.com/draveness/proxier/pkg/controller/proxier"
	"github.com/draveness/proxier/test/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProxierCreateSuite struct {
	suite.Suite
}

func (suite *ProxierCreateSuite) TestProxierCreateBackends() {
	ctx := f.NewTestCtx(suite.T())
	defer ctx.Cleanup(suite.T())

	namespace := ctx.CreateNamespace(suite.T(), f.KubeClient)
	ctx.SetupProxierRBAC(suite.T(), namespace, f.KubeClient)

	suite.T().Parallel()

	exampleProxier := framework.MakeBasicProxier(namespace, "test", []string{"v1", "v2"}, []int32{100, 10})

	_, err := f.CreateProxierAndWaitUntilReady(namespace, exampleProxier)

	assert.Nil(suite.T(), err, "create proxier error")

	svcList, err := f.KubeClient.CoreV1().Services(namespace).List(metav1.ListOptions{
		LabelSelector: "maegus.com/proxier-name=" + exampleProxier.Name,
	})

	assert.Nil(suite.T(), err, "list service error")
	assert.Equal(suite.T(), 2, len(svcList.Items), "proxier should create backend services")
}

func (suite *ProxierCreateSuite) TestProxierCreateNginxDeployment() {
	ctx := f.NewTestCtx(suite.T())
	defer ctx.Cleanup(suite.T())

	namespace := ctx.CreateNamespace(suite.T(), f.KubeClient)
	ctx.SetupProxierRBAC(suite.T(), namespace, f.KubeClient)

	suite.T().Parallel()

	exampleProxier := framework.MakeBasicProxier(namespace, "test", []string{"v1", "v2"}, []int32{100, 10})

	_, err := f.CreateProxierAndWaitUntilReady(namespace, exampleProxier)
	assert.Nil(suite.T(), err, "create proxier error")

	deploymentName := proxier.NewDeploymentName(exampleProxier)
	deployment, err := f.KubeClient.AppsV1().Deployments(namespace).Get(deploymentName, metav1.GetOptions{})

	assert.Nil(suite.T(), err, "get deployment error")
	assert.Equal(suite.T(), "nginx", deployment.Spec.Template.Spec.Containers[0].Name, "invalid nginx name")
	assert.Equal(suite.T(), "nginx:1.15.9", deployment.Spec.Template.Spec.Containers[0].Image, "invalid nginx image")
}

func (suite *ProxierCreateSuite) TestProxierCreateService() {
	ctx := f.NewTestCtx(suite.T())
	defer ctx.Cleanup(suite.T())

	namespace := ctx.CreateNamespace(suite.T(), f.KubeClient)
	ctx.SetupProxierRBAC(suite.T(), namespace, f.KubeClient)

	suite.T().Parallel()

	exampleProxier := framework.MakeBasicProxier(namespace, "echo", []string{"v1", "v2"}, []int32{100, 10})

	_, err := f.CreateProxierAndWaitUntilReady(namespace, exampleProxier)
	assert.Nil(suite.T(), err, "create proxier error")

	proxierService, err := f.KubeClient.CoreV1().Services(namespace).Get(exampleProxier.Name, metav1.GetOptions{})
	assert.Nil(suite.T(), err, "get service error")
	assert.Equal(suite.T(), exampleProxier.Name, proxierService.Name, "proxier and service shoud have the same name")
}
