package tester

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"sync"
	"testing"
)

// T is an interface that mimics *testing.T.
type T interface {
	Run(name string, f func(t T)) bool
	Cleanup(f func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Setenv(key, value string)
	Chdir(dir string)
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	TempDir() string
	Context() context.Context
}

// TestResult represents the reason for a test termination via panic.
type TestResult string

const (
	ResultFail  TestResult = "FAIL"
	ResultFatal TestResult = "FATAL"
	ResultSkip  TestResult = "SKIP"
)

// TestContext mimics *testing.T for standalone use.
type TestContext struct {
	ctx      context.Context
	cf       MCPClientFactory
	name     string
	tempRoot string
	logger   *log.Logger
	failed   bool
	skipped  bool
	cleanup  []func()
	tempDirs []string
	mu       sync.Mutex // Protects shared state
}

// NewTestContext creates a new TestContext with the given name and context.
func NewTestContext(ctx context.Context, name string, tempDir string, cf MCPClientFactory) *TestContext {
	prefix := ""
	if name != "" {
		prefix = fmt.Sprintf("[%s] ", name)
	}
	return &TestContext{
		ctx:      ctx,
		cf:       cf,
		name:     name,
		tempRoot: tempDir,
		logger:   log.New(os.Stdout, prefix, 0),
	}
}

func (t *TestContext) Run(name string, f func(t T)) bool {
	if t.name != "" {
		name = t.name + "_" + name
	}
	prefix := fmt.Sprintf("[%s] ", name)
	tt := &TestContext{
		ctx:      t.ctx,
		logger:   log.New(os.Stdout, prefix, 0),
		name:     name,
		tempRoot: t.tempRoot,
	}
	// Run the test in a panic-safe block
	defer func() {
		if r := recover(); r != nil {
			switch r {
			case ResultFail:
				tt.Log("Test failed (via FailNow)")
			case ResultFatal:
				tt.Log("Test failed (via Fatal)")
			case ResultSkip:
				tt.Log("Test skipped (via SkipNow)")
				tt.skipped = true
			default:
				// Unexpected panic, rethrow
				panic(r)
			}
		}
		tt.RunCleanup()
		if tt.Skipped() {
			tt.Log(color.YellowString("SKIPPED"))
		} else if tt.Failed() {
			tt.Log(color.RedString("FAILED"))
			t.Fail()
		} else {
			tt.Log(color.GreenString("PASSED"))
		}
	}()
	f(tt)
	return true
}

// Cleanup registers a function to be called when the test completes.
func (t *TestContext) Cleanup(f func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.cleanup = append(t.cleanup, f)
}

// RunCleanup executes all registered cleanup functions in reverse order.
func (t *TestContext) RunCleanup() {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i := len(t.cleanup) - 1; i >= 0; i-- {
		t.cleanup[i]()
	}
	t.cleanup = nil
}

// Error logs an error and marks the test as failed.
func (t *TestContext) Error(args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.logger.Println("ERROR:", fmt.Sprint(args...))
	t.failed = true
}

// Errorf logs a formatted error and marks the test as failed.
func (t *TestContext) Errorf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.logger.Printf("ERROR: "+format, args...)
	t.failed = true
}

// Fail marks the test as failed without logging.
func (t *TestContext) Fail() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.failed = true
}

// FailNow marks the test as failed and panics to stop execution.
func (t *TestContext) FailNow() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.failed = true
	panic(ResultFail)
}

// Failed reports whether the test has failed.
func (t *TestContext) Failed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.failed
}

// Fatal logs an error and panics to stop execution.
func (t *TestContext) Fatal(args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.logger.Println("FATAL:", fmt.Sprint(args...))
	t.failed = true
	panic(ResultFatal)
}

// Fatalf logs a formatted error and panics to stop execution.
func (t *TestContext) Fatalf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.logger.Printf("FATAL: "+format, args...)
	t.failed = true
	panic(ResultFatal)
}

// Helper is a no-op; irrelevant outside the testing framework.
func (t *TestContext) Helper() {
	// No-op
}

// Log outputs a message.
func (t *TestContext) Log(args ...interface{}) {
	t.logger.Println(fmt.Sprint(args...))
}

// Logf outputs a formatted message.
func (t *TestContext) Logf(format string, args ...interface{}) {
	t.logger.Printf(format, args...)
}

// Name returns the test name.
func (t *TestContext) Name() string {
	return t.name
}

// Setenv sets an environment variable for the duration of the test.
func (t *TestContext) Setenv(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		t.Errorf("Setenv failed: %v", err)
		return
	}
	t.cleanup = append(t.cleanup, func() {
		if err := os.Unsetenv(key); err != nil {
			t.Logf("Failed to unset env %s: %v", key, err)
		}
	})
}

// Chdir changes the current working directory, reverting on cleanup.
func (t *TestContext) Chdir(dir string) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get current directory: %v", err)
		return
	}
	if err := os.Chdir(dir); err != nil {
		t.Errorf("Failed to change directory to %s: %v", dir, err)
		return
	}
	t.cleanup = append(t.cleanup, func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Logf("Failed to restore directory to %s: %v", originalDir, err)
		}
	})
}

// Skip logs a skip message and marks the test as skipped.
func (t *TestContext) Skip(args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.logger.Println("SKIP:", fmt.Sprint(args...))
	t.skipped = true
}

// SkipNow marks the test as skipped and panics to stop execution.
func (t *TestContext) SkipNow() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.skipped = true
	panic(ResultSkip)
}

// Skipf logs a formatted skip message and marks the test as skipped.
func (t *TestContext) Skipf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.logger.Printf("SKIP: "+format, args...)
	t.skipped = true
}

// Skipped reports whether the test was skipped.
func (t *TestContext) Skipped() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.skipped
}

// TempDir creates and returns a temporary directory, cleaned up later.
func (t *TestContext) TempDir() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	dir, err := os.MkdirTemp(t.tempRoot, "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.tempDirs = append(t.tempDirs, dir)
	t.cleanup = append(t.cleanup, func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("Failed to remove temp dir %s: %v", dir, err)
		}
	})
	return dir
}

// Context returns the test’s context.
func (t *TestContext) Context() context.Context {
	return t.ctx
}

func Run(t interface{}, name string, f func(t T)) bool {
	if tt, ok := t.(T); ok {
		return tt.Run(name, f)
	}
	if tt, ok := t.(*testing.T); ok {
		return tt.Run(name, func(t *testing.T) {
			f(testingAdapter{t})
		})
	}
	return false
}
