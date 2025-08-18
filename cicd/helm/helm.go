package helm

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
)

var (
	ErrKindNotSupported = errors.New("kind is not support")
	ErrNoKindDefined    = errors.New("no kind field is defined")
)

type Helm interface {
	Render() (resources []byte, hooks []byte, err error)
	MustRender() (resources []byte, hooks []byte)
	Resources() (Resources, error)
	MustResources() Resources
}

type helm struct {
	chartPath  string
	valueFiles []string
	values     []string
}

type Option func(*helm)

func WithValuesYaml(path string) Option {
	return func(h *helm) {
		h.valueFiles = append(h.valueFiles, path)
	}
}

func WithSet(key, value string) Option {
	return func(h *helm) {
		h.values = append(h.values, key+"="+value)
	}
}

func New(chart string, opts ...Option) Helm {
	h := &helm{
		chartPath: chart,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (h *helm) Render() ([]byte, []byte, error) {
	settings := cli.New()
	actionConfig := &action.Configuration{}

	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debugLog); err != nil {
		return nil, nil, errors.Wrap(err, "error initializing helm")
	}

	chart, err := loader.Load(h.chartPath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error loading chart")
	}

	client := action.NewInstall(actionConfig)
	client.DryRun = true
	client.ReleaseName = "my-release"
	client.Namespace = settings.Namespace()
	client.ClientOnly = true
	client.Verify = true
	client.DisableHooks = false

	valueOpts := &values.Options{
		Values:     h.values,
		ValueFiles: h.valueFiles,
	}

	vals, err := valueOpts.MergeValues(nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error merging values")
	}

	rel, err := client.Run(chart, vals)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error rendering")
	}

	var hooks []string
	for _, hook := range rel.Hooks {
		hooks = append(hooks, hook.Manifest)
	}

	return []byte(rel.Manifest), []byte(strings.Join(hooks, "---")), nil
}

func (h *helm) MustRender() (resources []byte, hooks []byte) {
	resources, hooks, err := h.Render()
	if err != nil {
		panic(err)
	}

	return resources, hooks
}

func (h *helm) Resources() (Resources, error) {
	resources, hooks, err := h.Render()
	if err != nil {
		return nil, errors.Wrap(err, "error rendering Helm chart")
	}

	data := append(resources, hooks...)

	documents := Resources{}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	for {
		document := Resource{}
		if err := decoder.Decode(&document); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, errors.Wrap(err, "error decoding YAML document")
		}
		if document.IsEmpty() {
			continue
		}
		documents = append(documents, document)
	}

	return documents, nil
}

func (h *helm) MustResources() Resources {
	rs, err := h.Resources()
	if err != nil {
		panic(err)
	}
	return rs
}

func debugLog(msg string, args ...interface{}) {
	log.Debugf(msg, args...)
}
