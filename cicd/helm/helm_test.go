package helm

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

const chartPath = "testdata/chart"

func (s *helmTestSuite) TestHelmVerify() {
	resources, hooks, err := s.helm.Render()
	s.Require().NoError(err)
	s.Require().NotEmpty(resources)
	s.Require().NotEmpty(hooks)
}

func (s *helmTestSuite) TestRender() {
	resources, hooks, err := s.helm.Render()
	s.Require().NoError(err)
	s.Require().NotEmpty(resources)
	s.Require().NotEmpty(hooks)
}

func (s *helmTestSuite) TestResources() {
	resources, err := s.helm.Resources()
	s.Require().NoError(err)
	s.Require().Len(resources, 4)
	s.Require().Len(resources.FilterByKind("deployment"), 1)
	s.Require().Len(resources.FilterByKind("service"), 1)
	s.Require().Len(resources.FilterByKind("ingress"), 1)
	s.Require().Len(resources.FilterByKind("job"), 1)
}

func (s *helmTestSuite) TestFilterByKind() {
	resources, err := s.helm.Resources()
	s.Require().NoError(err)
	ingresses := resources.FilterByKind("ingress")
	s.Require().Len(ingresses, 1)
	s.Require().Equal("Ingress", ingresses[0].GetString("kind"))
}

func (s *helmTestSuite) TestGetString() {
	resources, err := s.helm.Resources()
	s.Require().NoError(err)
	deployments := resources.FilterByKind("deployment")
	s.Require().Len(deployments, 1)
	s.Require().Equal("LOG_LEVEL", deployments[0].GetString("spec.template.spec.containers.0.env.0.name"))
}

func (s *helmTestSuite) TestGetNumber() {
	resources, err := s.helm.Resources()
	s.Require().NoError(err)
	services := resources.FilterByKind("service")
	s.Require().Len(services, 1)
	s.Require().Equal(5555.0, services[0].GetNumber("spec.ports.0.port"))
}

func (s *helmTestSuite) TestGetBoolean() {
	resources, err := s.helm.Resources()
	s.Require().NoError(err)
	deployments := resources.FilterByKind("deployment")
	s.Require().Len(deployments, 1)
	s.Require().Equal(true, deployments[0].GetBoolean("spec.template.spec.containers.0.securityContext.readOnlyRootFilesystem"))
}

func (s *helmTestSuite) TestIsExists() {
	resources, err := s.helm.Resources()
	s.Require().NoError(err)
	deployments := resources.FilterByKind("deployment")
	s.Require().Len(deployments, 1)
	s.Require().Equal(true, deployments[0].IsExists("spec.template.spec.containers.0"))
	s.Require().Equal(false, deployments[0].IsExists("spec.template.spec.containers.1"))
}

func (s *helmTestSuite) TestGetStruct() {
	resources, err := s.helm.Resources()
	s.Require().NoError(err)
	deployments := resources.FilterByKind("deployment")
	s.Require().Len(deployments, 1)

	type envVars struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	env := []envVars{}
	err = deployments[0].GetStruct("spec.template.spec.containers.0.env", &env)
	s.Require().NoError(err)
	s.Require().Equal([]envVars{
		{
			Name:  "LOG_LEVEL",
			Value: "trace",
		},
		{
			Name:  "ANOTHER_VAR",
			Value: "anotherValue",
		},
	}, env)
}

// ========================================================================
// Test suite setup
// ========================================================================

type helmTestSuite struct {
	suite.Suite

	helm Helm
}

func (s *helmTestSuite) SetupTest() {
	s.helm = New(chartPath, WithValuesYaml("testdata/chart/values.yaml"))
}

func TestHelmTestSuite(t *testing.T) {
	suite.Run(t, &helmTestSuite{})
}
