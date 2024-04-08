---
name: "route-handlers"
description: "Route Handlers Generator"
prompts:
  - type: "base"
    name: "name"
    message: "Type a route handler name"
    validate: "{{if eq (len .input) 0}}Name must have more than - characters{{end}}"
  - type: "base"
    name: "parent"
    message: "Type a path of the parent folder if needed"
    suffix: "* [Optional] Skip to enter"
    defaultValues: ""
  - type: "multiselect"
    name: "methods"
    message: "Select HTTP methods you need"
    choices:
      - "GET"
      - "POST"
      - "PUT"
      - "DELETE"
      - "PATCH"
      - "HEAD"
      - "OPTIONS"
actions:
  - type: "add"
    path: "src/app/api{{if .parent}}/{{.parent}}{{end}}/{{.name}}/route.ts"
    template: "handler.ts"
---

# handler.ts

```ts
import { NextRequest, NextResponse } from "next/server";
{{range .methods}}
export async function {{.}}(req: NextRequest, res: NextResponse) {
  return res;
}
{{end}}
```
