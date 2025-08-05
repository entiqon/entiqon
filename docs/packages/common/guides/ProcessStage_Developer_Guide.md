# ProcessStage Developer Guide

## Overview

`ProcessStage` is a simple but powerful string alias used to identify specific stages or steps within a multi-phase process or workflow. Incorporating `ProcessStage` into your error handling helps provide rich contextual information about *where* an error occurred, facilitating better debugging and user feedback.

Typical predefined stages include:

- `Init`
- `Parse`
- `Validate`
- `Execute`
- `Finalize`

## Using `ProcessStage`

Use the predefined constants from the `common` package whenever possible to maintain consistency:

```go
import "github.com/entiqon/entiqon/common"
import "github.com/entiqon/entiqon/common/errors"

err := errors.NewProcessStageError(common.StageParse, "Parser", "Invalid syntax", underlyingErr)
```

When handling errors, you can extract the stage information to customize behavior:

```go
var pse *errors.ProcessStageError
if errors.As(err, &pse) {
    switch pse.Stage() {
    case common.StageValidate:
        // Special handling for validation errors
    default:
        // Generic error handling
    }
}
```

## Extending `ProcessStage`

Since `ProcessStage` is a string alias, you can define custom stages tailored to your domain:

```go
package mypackage

import "github.com/entiqon/entiqon/common"

const (
    StageCustomProcessing common.ProcessStage = "CustomProcessing"
    StagePostProcess      common.ProcessStage = "PostProcess"
)
```

Use your custom stages seamlessly wherever `ProcessStage` is expected:

```go
err := errors.NewProcessStageError(StageCustomProcessing, "CustomModule", "Unexpected condition", nil)
```

### Naming Guidelines

- Choose descriptive and concise names.
- Use PascalCase or CamelCase consistently.
- Document custom stages for team clarity.

## Best Practices

- Prefer existing predefined stages to avoid fragmentation.
- Add custom stages only when meaningful distinctions are needed.
- Keep stage granularity balanced — too fine-grained may cause noise.
- Always attach stage info to errors to enhance observability.

---

If you want, I can also help format this as a standalone markdown file or prepare other documentation parts. Just say the word!
