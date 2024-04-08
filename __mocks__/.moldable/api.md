---
name: "api"
description: "API Generator"
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
    message: "Select Express HTTP methods you need"
    choices:
      - "get"
      - "post"
      - "put"
      - "delete"
      - "patch"
      - "head"
      - "options"
actions:
  - type: "add"
    path: "src/apis{{if .parent}}/{{.parent}}{{end}}/{{.name}}.ts"
    template: "api.ts"
---

# api.ts

```ts
import express from "express";

const app = express();
app.use(express.json());
{{range .methods}}
app.{{.}}("/", async (req, res) => {
  res.json({ message: "Hi There!" });
});
{{end}}
```
