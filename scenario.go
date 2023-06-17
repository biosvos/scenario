package scenario

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"time"
)

type Scenario struct {
	artifacts *Artifacts
	steps     []Step
}

func NewScenario(steps []Step, opts ...Option) (*Scenario, error) {
	artifacts, _ := NewArtifacts()
	options := Options{
		artifacts: artifacts,
	}
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &Scenario{
		artifacts: artifacts,
		steps:     steps,
	}, nil
}

func printf(writer io.Writer, format string, a ...any) {
	if isEmptyWriter(writer) {
		return
	}
	_, _ = io.WriteString(writer, fmt.Sprintf(format, a...))
}

func isEmptyWriter(writer io.Writer) bool {
	return writer == nil
}

func (s *Scenario) Run(writer io.Writer) error {
	for idx, step := range s.steps {
		printf(writer, "step %v. %v\n", idx+1, step.Title())
		err := s.artifacts.Fill(step.Input())
		if err != nil {
			return errors.WithStack(err)
		}
		start := time.Now()
		outputs, err := step.Run()
		printf(writer, "%v elapsed\n", time.Since(start))
		if err != nil {
			return errors.WithStack(err)
		}
		if !outputs.IsExists(step.DefinitionOfDone()...) {
			return errors.New("step is not done.")
		}
		for key, value := range outputs.artifacts {
			err := s.artifacts.Add(key, value)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
}
