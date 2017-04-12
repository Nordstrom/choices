package choices

import (
	"context"
	"fmt"

	"github.com/Nordstrom/choices/util"
	"github.com/foolusion/elwinprotos/storage"
	"github.com/pkg/errors"
)

// CreateExperiment will create a new experiment and namespace based on the
// input that it receives.
func CreateExperiment(
	ctx context.Context,
	expcontroller experimentController,
	sExperiment *storage.Experiment,
	sNamespace *storage.Namespace,
	nsNumSegments int,
	expNumSegments int,
) error {
	if sExperiment == nil {
		return errors.New("experiment is nil")
	} else if len(sExperiment.Labels) == 0 {
		return errors.New("experiment labels are empty")
	}
	exp := FromExperiment(sExperiment)
	var ns *Namespace
	if sNamespace == nil {
		ns = newNamespace(exp.Namespace, nsNumSegments)
	} else {
		ns = FromNamespace(sNamespace)
	}
	// sample the namespaces segments
	seg := ns.Segments.sample(expNumSegments)
	exp.Segments = &segments{b: seg, len: ns.Segments.len}
	if exp.Name == "" {
		exp.Name = util.BasicNameGenerator.GenerateName("")
	}
	if exp.ID == "" {
		exp.ID = util.BasicNameGenerator.GenerateName(fmt.Sprintf("exp-%s-", exp.Name))
	}

	if err := expcontroller.SetNamespace(ctx, ns.ToNamespace()); err != nil {
		return errors.Wrap(err, "could not save namespace")
	}
	if err := expcontroller.SetExperiment(ctx, exp.ToExperiment()); err != nil {
		return errors.Wrap(err, "could not save experiment")
	}
	return nil
}

// BadSegments implements an Error interface. It is used when then
// experiments claimed segments do not match the namespaces claimed
// segments.
type BadSegments struct {
	NamespaceSegments *segments
	*Experiment
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
	nsSet := make(map[string]*segments, len(namespaces))
	for _, ns := range namespaces {
		nsSet[ns.Name] = FromSegments(ns.Segments)
	}
	experiments, err := e.AllExperiments(ctx)
	if err != nil {
		return errors.Wrap(err, "could not get all experiments from storage")
	}
	expSet := make(map[string]*segments, len(namespaces))
	for _, exp := range experiments {

		if s, ok := expSet[exp.Namespace]; !ok {
			expSet[exp.Namespace] = FromSegments(exp.Segments)
		} else {
			// check for overlapping experiments
			out, err := s.Claim(FromSegments(exp.Segments))
			if err != nil {
				return &BadSegments{
					NamespaceSegments: s,
					Experiment:        FromExperiment(exp),
					Err:               err,
				}
			}
			s.b = out
		}
		// check all namespace segments are claimed
		if s, ok := nsSet[exp.Namespace]; !ok {
			return NamespaceDoesNotExist{FromExperiment(exp)}
		} else if !s.contains(FromSegments(exp.Segments)) {
			return &BadSegments{NamespaceSegments: s, Experiment: FromExperiment(exp)}
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
			if err := e.SetNamespace(ctx, &storage.Namespace{
				Name:        err.Namespace,
				NumSegments: int64(err.Segments.len),
				Segments:    err.Segments.ToSegments(),
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
	SetNamespace(context.Context, *storage.Namespace) error
	Namespace(context.Context, string) (*storage.Namespace, error)
	AllNamespaces(context.Context) ([]*storage.Namespace, error)
	SetExperiment(context.Context, *storage.Experiment) error
	Experiment(context.Context, string) (*storage.Experiment, error)
	AllExperiments(context.Context) ([]*storage.Experiment, error)
}
