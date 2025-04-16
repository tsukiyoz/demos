package errors

import (
	"errors"
	"fmt"
	"testing"

	githuberrors "github.com/pkg/errors"
	k8sutilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func TestK8sUtilErrorsNew(t *testing.T) {
	var errs []error

	err1 := errors.New("error1 message")
	errs = append(errs, err1)

	err2 := errors.New("error2 message")
	errs = append(errs, err2)

	agg := k8sutilerrors.NewAggregate(errs)

	t.Logf("errs: %s\n", errs)
	t.Logf("aggregate: %s\n", agg)
	t.Logf("errs len: %d, aggregate len: %d\n", len(errs), len(agg.Errors()))
	t.Logf("errors: %s\n", agg.Errors())
	t.Logf("err1: %s, err2: %s\n", err1, err2)
	t.Logf("errors is [agg, err1]: %t\n", errors.Is(agg, err1))
	t.Logf("errors is [agg, err2]: %t\n", errors.Is(agg, err2))
	t.Logf("errors is [agg, new error]: %t\n", errors.Is(agg, errors.New("new error")))

	filteredErr := k8sutilerrors.FilterOut(agg, func(err error) bool {
		return errors.Is(err, err1)
	})

	t.Logf("filtered errors: %s\n", filteredErr)
}

func TestK8sUtilErrors_Flatten(t *testing.T) {
	processPod := func(podID int) error {
		if podID%2 == 0 {
			return fmt.Errorf("pod %d processing failed", podID)
		}
		return nil
	}

	batchProcessPods := func(podIDs []int) error {
		var errs []error
		for _, podID := range podIDs {
			if err := processPod(podID); err != nil {
				errs = append(errs, err)
			}
		}
		return k8sutilerrors.NewAggregate(errs)
	}

	handleAllBatches := func() error {
		var errs []error

		batch1Err := batchProcessPods([]int{1, 2, 3})
		if batch1Err != nil {
			errs = append(errs, batch1Err)
		}

		batch2Err := batchProcessPods([]int{4, 5, 6})
		if batch2Err != nil {
			errs = append(errs, batch2Err)
		}

		batch3Err := batchProcessPods([]int{7, 8, 9})
		if batch3Err != nil {
			errs = append(errs, batch3Err)
		}

		return k8sutilerrors.Flatten(k8sutilerrors.NewAggregate(errs))
	}

	err := handleAllBatches()
	if err != nil {
		if agg, ok := err.(k8sutilerrors.Aggregate); ok {
			t.Logf("total errors: %d\n", len(agg.Errors()))
			for i, e := range agg.Errors() {
				t.Logf("error %d: %s\n", i, e)
			}
		} else {
			t.Logf("not aggregate error: %s\n", err)
		}
	} else {
		t.Logf("all batches processed successfully\n")
	}
}

// 测试代码片段
func BenchmarkErrorsNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = errors.New("error message")
	}
}

func BenchmarkFmtErrorf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Errorf("error %s", "message")
	}
}

func BenchmarkGithubPkgErrorsNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = githuberrors.New("error message")
	}
}

func BenchmarkK8sUtilErrorsNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = k8sutilerrors.NewAggregate([]error{errors.New("error message")})
	}
}

// ❯ go test -bench .
// goos: darwin
// goarch: arm64
// pkg: github.com/tsukiyoz/demos/usecases/errors
// cpu: Apple M2
// BenchmarkErrorsNew-8                    1000000000               0.2932 ns/op
// BenchmarkFmtErrorf-8                    24209167                49.17 ns/op
// BenchmarkGithubPkgErrorsNew-8            4761103               248.4 ns/op
// BenchmarkK8sUtilErrorsNew-8             41058860                28.19 ns/op
// PASS
// ok      github.com/tsukiyoz/demos/usecases/errors       4.593s
