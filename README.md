# Run GitHub Action

A Bitrise Step for running GitHub Javascript Actions

## Usage

Please note that not all GitHub Action features are supported, if the Action
relies on an unsupported feature (e.g masking sensitive values), the step might
fail or produce an unexpected output.

Add it to your bitrise.yml using git source. Inputs:
- *uses*: the identifier of the GitHub Action to run, use [the official syntax](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idstepsuses) (required)
- *with*: inputs for the GitHub Action, use [the official
  syntax](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idstepswith), basically key-value pairs in YML syntax. For sensitive data, consider using secrets and use an environment variable here. Alternatively, if there's an enviroment variable already set with the same input name (e.g. if you set it up as an app/worklfow env or a secret), the step can pick that up and use it, even if it's not explicitly defined in the *with* input.

```yaml
workflows:
  your-workflow:
    steps:
      - git::https://github.com/gaborszakacs/steps-run-github-action.git@master:
          inputs:
            - uses: gaborszakacs/kitchen-sink-action@main
            - with: |
                name: Bello
                a_sensitive_input: $SECRET_ENV
```

## Supported and not yet supported features

- [x] Basic logging (not yet Bitrise flavoured though)
- [ ] Adding context to logs (file, line, column)
- [ ] Grouping log lines
- [ ] Switching debug level on or off
- [x] Using input variables
- [x] Setting output variables
- [ ] Masking sensitive values in logs
- [ ] Adding to the system path

## How does it work?

The step
- clones the GitHub Action, and parses its config (inputs, entrypoint)
- prepares the inputs: the core toolkit expects the inputs to be present as an
  env var prefixed with `INPUT_`
- runs the action
- processes the output, executes the commands in a Bitrise specific way (e.g
  sets output variables) or prints the output. (GitHub Action [uses special
keywords](https://docs.github.com/en/actions/reference/workflow-commands-for-github-actions) in the output as workflow commands be executed by the runner.)

An alternative approach would be to patch the core toolkit library before
running the step, that seems to be a bit more fragile as Action publishers don't
have to use the toolkit.


## Contributing

Run an E2E test with `go test e2e`. This runs the `test` bitrise workflow which
runs this step with a [prepared kitchen sink GitHub Action](https://github.com/gaborszakacs/kitchen-sink-action) and verifies the step's behaviour by matcing the build log output with the expected patterns.
