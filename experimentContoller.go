package choices

import (
	"context"
	"fmt"

	"github.com/Nordstrom/choices/util"
	"github.com/pkg/errors"
)

// CreateExperiment will create a new experiment and namespace based on the
// input that it receives.
func CreateExperiment(ctx context.Context, s experimentController, exp *Experiment, ns *Namespace, nsSegments, expSegments int) error {
	var err error
	if ns == nil {
		ns = newNamespace(ns.Name, nsSegments)
	} else {
		ns, err = s.Namespace(ns.Name)
		if err != nil {
			switch err {
			case ErrNotFound:
				ns = newNamespace(ns.Name, nsSegments)
			default:
				return errors.Wrap(err, "could not get namespace")
			}
		}
	}
	// sample the namespaces segments
	seg := ns.sampleSegments(expSegments)

	if exp == nil {
		return errors.New("experiment is nil")
	} else if len(exp.Labels) == 0 {
		return errors.New("experiment labels are empty")
	}
	exp.Segments = seg
	if exp.Name == "" {
		exp.Name = util.BasicNameGenerator.GenerateName("")
	}
	if exp.ID == "" {
		exp.ID = util.BasicNameGenerator.GenerateName(fmt.Sprintf("exp-%s-", exp.Name))
	}

	if err := s.SetNamespace(ns); err != nil {
		return errors.Wrap(err, "could not save namespace")
	}
	if err := s.SetExperiment(exp); err != nil {
		return errors.Wrap(err, "could not save experiment")
	}
	return nil
}

var (
	ErrNotFound = errors.New("not found")
)

type experimentController interface {
	SetNamespace(namespace *Namespace) error
	Namespace(name string) (*Namespace, error)
	SetExperiment(*Experiment) error
	Experiment(id string) (*Experiment, error)
}
