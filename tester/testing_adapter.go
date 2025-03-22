package tester

import (
	"context"
	"testing"
)

type testingAdapter struct {
	*testing.T
}

func Wrap(t *testing.T) T {
	return testingAdapter{t}
}

func (t testingAdapter) Run(name string, f func(t T)) bool {
	return t.T.Run(name, func(t *testing.T) {
		f(testingAdapter{t})
	})
}

func (t testingAdapter) Cleanup(f func()) {
	t.T.Cleanup(f)
}

func (t testingAdapter) Error(args ...interface{}) {
	t.T.Error(args...)
}

func (t testingAdapter) Errorf(format string, args ...interface{}) {
	t.T.Errorf(format, args...)
}

func (t testingAdapter) Fail() {
	t.T.Fail()
}

func (t testingAdapter) FailNow() {
	t.T.FailNow()
}

func (t testingAdapter) Failed() bool {
	return t.T.Failed()
}

func (t testingAdapter) Fatal(args ...interface{}) {
	t.T.Fatal(args...)
}

func (t testingAdapter) Fatalf(format string, args ...interface{}) {
	t.T.Fatalf(format, args...)
}

func (t testingAdapter) Helper() {
	t.T.Helper()
}

func (t testingAdapter) Log(args ...interface{}) {
	t.T.Log(args...)
}

func (t testingAdapter) Logf(format string, args ...interface{}) {
	t.T.Logf(format, args...)
}

func (t testingAdapter) Name() string {
	return t.T.Name()
}

func (t testingAdapter) Setenv(key, value string) {
	t.T.Setenv(key, value)
}

func (t testingAdapter) Chdir(dir string) {
	t.T.Chdir(dir)
}

func (t testingAdapter) Skip(args ...interface{}) {
	t.T.Skip(args...)
}

func (t testingAdapter) SkipNow() {
	t.T.SkipNow()
}

func (t testingAdapter) Skipf(format string, args ...interface{}) {
	t.T.Skipf(format, args...)
}

func (t testingAdapter) Skipped() bool {
	return t.T.Skipped()
}

func (t testingAdapter) TempDir() string {
	return t.T.TempDir()
}

func (t testingAdapter) Context() context.Context {
	return t.T.Context()
}
