---
name: "benchmark"
root: "."
output: "src/scaffdog"
ignore: []
questions:
  name: "Please enter any text."
---

# `{{ inputs.name }}.tsx`

```tsx
import React from "react";

export const {{ inputs.name }}: React.FC = (children: { children: React.ReactNode }) => <div>{children}</div>;
```

