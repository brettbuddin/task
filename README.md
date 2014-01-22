## Task

Tasks are loosely modeled after Unix-style processes. Tasks have the following properties:

- An identitifying name.
- Accept an array of arguments (like argv)
- Have access to global environment variables
- Provide input, output and error streams
- Output and error streams can have many different consumers (tailers)
- Returning status code to indicate success or error

### Defining Tasks

Native Go:
```go
task.Define("trim", func(t *task.Task) task.Status {
    content := bytes.NewBuffer(nil)
    if _, err := t.InputStream.CopyTo(content); err != nil {
        return task.StatusErr
    }
        
    trimmed := strings.TrimString(content.String())
    t.OutputStream.Write([]byte(trimmed))
    return task.StatusOk
})

```

Using an external command (via `os/exec`):
```go
task.Define("count-words", func(t *task.Task) task.Status {
    cmd := exec.Command("wc", "-w")
    cmd.Stdin = task.InputStream
    cmd.Stdout = task.OutputStream
    cmd.Stderr = task.ErrorStream
    err := cmd.Run()
    
    if err != nil {
        return task.StatusErr
    }
    
    return task.StatusOk
})
```

### Running a Task

```go
t := task.New("trim")
    
// Set the input
t.InputStream.Set(strings.NewReader("  padded  "))

// Collect the output
output := bytes.NewBuffer(nil)
t.OutputStream.Consume(output)

// Run the task and check the status
if err := t.Run(); err != nil || t.Status == task.StatusErr {
    panic("failed to trim the string")
}

fmt.Println(output.String())
// => "padded"
```
