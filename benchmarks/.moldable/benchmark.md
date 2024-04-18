---
name: "benchmark"
description: "This is a benchmark generator"
prompts:
  - type: "base"
    name: "name"
    message: "Please enter any text."
actions:
  - type: "add"
    path: "src/moldable/{{.name}}.tsx"
    template: "benchmark.tsx"
---

# benchmark.tsx

```tsx
import React from "react";

export const {{.name}}: React.FC = (children: { children: React.ReactNode }) => <div>{children}</div>;`,
```
