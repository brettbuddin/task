package task

import (
	"bytes"
	"math/rand"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestDefinition(t *testing.T) {
	container := NewContainer()
	container.Define("trim", func(task *Task) Status {
		return StatusOk
	})

	if handler := container.Handler("trim"); handler == nil {
		t.Fail()
	}
}

func TestTaskExecution(t *testing.T) {
	container := NewContainer()
	container.Define("trim", func(task *Task) Status {
		body := bytes.NewBuffer(nil)
		if _, err := task.InputStream.CopyTo(body); err != nil {
			return StatusErr
		}

		// Jitter the time it takes this task to run
		// so we can check the waiting on output close.
		random := rand.New(rand.NewSource(0))
		time.Sleep(jitter(500*time.Millisecond, random))

		trimmed := strings.TrimSpace(body.String())
		task.OutputStream.Write([]byte(trimmed))

		return StatusOk
	})

	output := bytes.NewBuffer(nil)
	output2 := bytes.NewBuffer(nil)
	output3 := bytes.NewBuffer(nil)

	task := container.NewTask("trim")
	task.InputStream.Set(strings.NewReader("  some string  "))
	task.OutputStream.Consume(output)
	task.OutputStream.Consume(output2)
	task.OutputStream.Consume(output3)
	err := task.Run()

	if err != nil || task.Status == StatusErr {
		t.Fail()
	}

	for _, o := range []*bytes.Buffer{output, output2, output3} {
		if o.String() != "some string" {
			t.Fail()
		}
	}
}

func TestExecCmd(t *testing.T) {
	container := NewContainer()
	container.Define("num-chars", func(task *Task) Status {
		cmd := exec.Command("wc", "-c")
		cmd.Stdin = task.InputStream
		cmd.Stdout = task.OutputStream
		cmd.Stderr = task.ErrorStream
		cmd.Run()

		return StatusOk
	})

	output := bytes.NewBuffer(nil)

	task := container.NewTask("num-chars")
	task.InputStream.Set(strings.NewReader("youareadummy"))
	task.OutputStream.Consume(output)
	task.Run()

	if strings.TrimSpace(output.String()) != "12" {
		t.Fail()
	}
}

func jitter(d time.Duration, r *rand.Rand) time.Duration {
	if d == 0 {
		return 0
	}
	return d + time.Duration(r.Int63n(2*int64(d)))
}
