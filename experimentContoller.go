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
		ns, err = s.Namespace(ctx, ns.Name)
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

	if err := s.SetNamespace(ctx, ns); err != nil {
		return errors.Wrap(err, "could not save namespace")
	}
	if err := s.SetExperiment(ctx, exp); err != nil {
		return errors.Wrap(err, "could not save experiment")
	}
	return nil
}

// BadSegments implements an Error interface. It is used when then
// experiments claimed segments do not match the namespaces claimed
// segments.
type BadSegments struct {
	NamespaceSegments segments
	Experiment
	Err error
}

func (bs *BadSegments) Error() string {
	if bs.Err == nil {
		return fmt.Sprintf("namespace %s segments %x don't match experiment %s segments %x", bs.Namespace, bs.NamespaceSegments, bs.ID, bs.Segments)
	}
	return fmt.Sprintf("namespace %s segments %x don't match experiment %s segments %x: %v", bs.Namespace, bs.NamespaceSegments, bs.ID, bs.Segments, bs.Err)
}

// NamespaceDoesNotExist is an error thrown when an experiment has
// a namespace listed that is not in storage.
type NamespaceDoesNotExist struct {
	*Experiment
}

func (n NamespaceDoesNotExist) Error() string {
	return fmt.Sprintf("namespace %s does not exist", n.Namespace)
}

// ValidateNamespaces checks whether all exeriments have namespaces
// and if the segments they claimed have also been claimed from the
// namespace. If everything is OK it ValidateNamespaces will return
// nil otherwise it will return an error. If the error is
// ErrNamespaceDoesNotExists you can fix this by creating a namespace
// to match the experiment.
func ValidateNamespaces(ctx context.Context, e experimentController) error {
	namespaces, err := e.AllNamespaces(ctx)
	nsSet := make(map[string]segments, len(namespaces))
	for _, ns := range namespaces {
		nsSet[ns.Name] = ns.Segments
	}
	experiments, err := e.AllExperiments(ctx)
	if err != nil {
		return errors.Wrap(err, "could not get all experiments from storage")
	}
	expSet := make(map[string]segments, len(namespaces))
	for _, exp := range experiments {
		if s, ok := expSet[exp.Namespace]; !ok {
			expSet[exp.Namespace] = exp.Segments
		} else {
			// check for overlapping experiments
			out, err := s.Claim(exp.Segments)
			if err != nil {
				return &BadSegments{
					NamespaceSegments: s,
					Experiment:        *exp,
					Err:               err,
				}
			}
			expSet[exp.Namespace] = out
		}
		// check all namespace segments are claimed
		if s, ok := nsSet[exp.Namespace]; !ok {
			return NamespaceDoesNotExist{exp}
		} else if !s.contains(exp.Segments) {
			return &BadSegments{NamespaceSegments: s, Experiment: *exp}
		}
	}
	return nil
}

// AutoFix will attempt to add namespaces for experiments
// that are missing a namespace. In the future we could potentially
// add more autofixes here.
func AutoFix(ctx context.Context, e experimentController) error {
	err := ValidateNamespaces(ctx, e)
	if err != nil {
		switch err := err.(type) {
		case *BadSegments:
			return errors.Wrap(err, "could not fix bad segments")
		case *NamespaceDoesNotExist:
			if err := e.SetNamespace(ctx, &Namespace{
				Name:        err.Namespace,
				NumSegments: len(err.Segments) * 8,
				Segments:    err.Segments,
			}); err != nil {
				return errors.Wrap(err, "could not add namespace")
			}
			return AutoFix(ctx, e)
		}
	}
	return nil
}

var (
	// ErrNotFound is the error that should be returnned when calls to
	// Namespace and Experiment fail because they don't exist in
	// storage.
	ErrNotFound = errors.New("not found")
)

type experimentController interface {
	SetNamespace(context.Context, *Namespace) error
	Namespace(context.Context, string) (*Namespace, error)
	AllNamespaces(context.Context) ([]*Namespace, error)
	SetExperiment(context.Context, *Experiment) error
	Experiment(context.Context, string) (*Experiment, error)
	AllExperiments(context.Context) ([]*Experiment, error)
}
