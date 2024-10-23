## [Cobra](https://github.com/spf13/cobra)

A library for creating powerful modern **CLI** applications.

:thumbsup:

- Easy (Nested) sub commands.
- Automatic help flag recognition of `-h`, `--help`, etc.
- Intelligent suggestions (`app srver`... did you mean `app server`?).
- Automatically generated shell autocomplete.
- Global, local and cascading flags, for example kubectl `--kubeconfig` is a global flag.

### Concepts

**Commands** represent **actions**.

**Args** are **things**/objects.

**Flags** are **modifiers** for those actions.

```bash
# kubectl - AppName
# edit    - COMMAND
# service - SUB COMMAND
# nginx   - ARG
# -n      - FLAG
$ kubectl edit service nginx -n kube-system
```

### Flags

> A way to modify the behavior of a command.

Optional

```go
# force
rootCmd.MarkFlagRequired("source")
```

#### Data Type

String

```bash
StringVar(&variable, "flag", "default", "description")                     # single flag
StringVarP(&variable, "flag", "shorthand", "default", "description")       # single flag with abbr
StringSliceVar(&variable, "flag", []string{}, "description")               # multi flags
StringSliceVarP(&variable, "flag", "shorthand", []string{}, "description") # multi flags with abbr
```

Bool

```bash
BoolVar(&variable, "flag", false, "description")
BoolVarP(&variable, "flag", "shorthand", false, "description")
```

Int

```bash
IntVar(&variable, "flag", 0, "description")
IntVarP(&variable, "flag", "shorthand", 0, "description")
```

Float

```bash
Float64Var(&variable, "flag", 0.0, "description")
Float64VarP(&variable, "flag", "shorthand", 0.0, "description")
```

#### Type

**Persistent**: can be used in current command & subcommands

```go
rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "k", defaultKubeconfig, "kubeconfig file")
rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "namespace")
```

**Local**: can be used in current command only

```go
worldCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
```

### QuickStart

```bash
$ go install github.com/spf13/cobra-cli@latest
$ cobra-cli init --author "aiops" --license mit
$ go build -o k8scopilot
$ ./k8scopilot
```

```bash
# ++command
$ cobra-cli add hello
$ go build -o k8scopilot
$ ./k8scopilot
```

:construction_worker: Parent only prints usage rather than biz logic.

```bash
# ++subcommand, specify parent
$ cobra-cli add world -p 'helloCmd'
$ go build -o k8scopilot
$ ./k8scopilot
```

```go
// ++ flag in root.go
// global flags = persistent flags under root
var kubeconfig string
var namespace string

func init() {
	// default to read ~/.kube/config
	homeDir, _ := os.UserHomeDir()
	defaultKubeconfig := filepath.Join(homeDir, ".kube", "config")
	rootCmd.PersistentFlags().
    	StringVarP(&kubeconfig, "kubeconfig", "conf", defaultKubeconfig, "Path to the kubeconfig file")
	rootCmd.PersistentFlags().
    	StringVarP(&namespace, "namespace", "n", defaultKubeconfig, "Namespace")
}
```

```go
// read from world.go
var worldCmd = &cobra.Command{
	// ...
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kubeconfig: ", kubeconfig)
		fmt.Println("namespace: ", namespace)
	},
}
```

```bash
$ go build -o k8scopilot
$ ./k8scopilot
$ ./k8scopilot hello world --kubeconfig config -n test
```

```go
// ++ local flag in world.go
var content string

func init() {
	helloCmd.AddCommand(worldCmd)
	worldCmd.Flags().StringVarP(&content, "content", "c", "world", "The content of msg")
}
```

```bash
$ go build -o k8scopilot
$ ./k8scopilot
$ ./k8scopilot hello world -h
```

```go
// version in root.go
&cobra.Command{
    Version: "v0.0.1",
}

// deprecation in world.go
&cobra.Command{
    Deprecated: "This cmd is deprecated, please use xxx instead",
}
```

