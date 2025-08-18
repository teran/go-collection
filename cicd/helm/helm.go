package helm

import (
	"bytes"
	"io"
	"os"

	"github.com/pkg/errors"
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
	Render() ([]byte, error)
	MustRender() []byte
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

func (h *helm) Render() ([]byte, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), nil); err != nil {
		return nil, errors.Wrap(err, "error initializing helm")
	}

	chart, err := loader.Load(h.chartPath)
	if err != nil {
		return nil, errors.Wrap(err, "error loading chart")
	}

	client := action.NewInstall(actionConfig)
	client.DryRun = true
	client.ReleaseName = "my-release"
	client.Namespace = settings.Namespace()
	client.ClientOnly = true

	valueOpts := &values.Options{
		Values:     h.values,
		ValueFiles: h.valueFiles,
	}

	vals, err := valueOpts.MergeValues(nil)
	if err != nil {
		return nil, errors.Wrap(err, "error merging values")
	}

	rel, err := client.Run(chart, vals)
	if err != nil {
		return nil, errors.Wrap(err, "error rendering")
	}

	return []byte(rel.Manifest), nil
}

func (h *helm) MustRender() []byte {
	out, err := h.Render()
	if err != nil {
		panic(err)
	}

	return out
}

func (h *helm) Resources() (Resources, error) {
	data, err := h.Render()
	if err != nil {
		return nil, errors.Wrap(err, "error rendering Helm chart")
	}
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
