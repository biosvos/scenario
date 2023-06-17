package scenario

import (
	"encoding/json"
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
		printHeader(writer, step, idx)
		err := s.artifacts.Fill(step.Input())
		if err != nil {
			return errors.WithStack(err)
		}
		printInput(writer, step.Input())
		start := time.Now()
		outputs, err := step.Run()
		if err != nil {
			return errors.WithStack(err)
		}
		if IsStepDone(outputs, step.DefinitionOfDone()) {
			return errors.New("step is not done.")
		}
		printOutput(writer, outputs)
		printElapsedTime(writer, start)
		err = s.artifacts.Merge(outputs)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func printElapsedTime(writer io.Writer, start time.Time) {
	printf(writer, "elapsed: %v\n", time.Since(start))
}

func printHeader(writer io.Writer, step Step, idx int) {
	printf(writer, "step %v. %v\n", idx+1, step.Title())
}

func printInput(writer io.Writer, input any) {
	marshal, _ := json.MarshalIndent(input, "", "\t")
	printf(writer, "input: %v\n", string(marshal))
}

func printOutput(writer io.Writer, output *Artifacts) {
	marshal, _ := json.MarshalIndent(output, "", "\t")
	printf(writer, "output: %v\n", string(marshal))
}

func IsStepDone(outputs *Artifacts, steps []string) bool {
	return !outputs.IsExists(steps...)
}
