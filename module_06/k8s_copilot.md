## K8s Copilot

```bash
# interactive
$ k8scopilot ask [llm]

# deploy a res: llm → YAML → (Function Calling) → client-go → k8s
> "Create a deployment for me, image is nginx"
> "Create a pod for me, image is nginx"

# get/list res(s): llm → (Function Calling) → client-go → k8s
> "list pod/svc/deploy under default ns"

# del a res: llm → (Function Calling) → client-go → k8s
> "Delete nginx deployment under default ns"
```

```bash
# interactive
$ k8scopilot analyze event

# features
# - analyse cluster event
# - get unhealthy pod log
# - ask gpt for suggestions
```

### Note

```bash
$ cobra-cli add ask
$ cobra-cli add chatgpt -p "askCmd"

$ cobra-cli add analyze
$ cobra-cli add event -p 'analyzeCmd'
```

```bash
# wsl env to win
$ export API_KEY="sk-NUkr3DrnPJtiWZS32d667dAaCd304f20B43bAbF0D6872b21"
$ export BASE_URL="https://vip.apiyi.com/v1"
$ export WSLENV=API_KEY/w:BASE_URL/w
```
