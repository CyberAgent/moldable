---
name: "quick-start"
description: "Quick Start"
prompts:
  - type: "base"
    name: "name"
    message: "Type your name"
    validate: "{{if eq (len .input) 0}}Name must have more than - characters{{end}}"
  - type: "select"
    name: "favoriteColor"
    message: "Select your favorite color"
    choices:
      - "Red"
      - "Green"
      - "Blue"
actions:
  - type: "add"
    path: "src/quick-start.json"
    template: "quick-start.json"
---

# quick-start.json

```json
{
  "name": "{{.name}}",
  "favoriteColor": "{{.favoriteColor}}"
}
```
