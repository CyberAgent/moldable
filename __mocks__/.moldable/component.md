---
name: "component"
description: "Component Generator"
prompts:
  - type: "select"
    name: "type"
    message: "Select a component type"
    choices:
      - "container"
      - "functional"
      - "layout"
      - "page"
  - type: "select"
    name: "pageType"
    message: "Select a page component type"
    choices:
      - "page"
      - "contents"
    when: '{{if ne .type "page"}}Skipped to select a page component type{{end}}'
  - type: "base"
    name: "name"
    message: "Type a component name"
    suffix: "* Use PascalCase"
    validate: "{{if eq (len .input) 0}}Name must have more than - characters{{end}}"
  - type: "base"
    name: "parent"
    message: "Type a path of the parent folder if needed"
    suffix: "* [Optional] Skip to enter"
    defaultValues: ""
  - type: "confirm"
    name: "isSC"
    message: "Is it React Server Component?"
    defaultValues: true
actions:
  - type: "add"
    path: "src/components/container/{{.type}}{{if .parent}}/{{.parent}}{{end}}/{{.name | pascal}}.tsx"
    template: "Container.tsx"
    skip: '{{if ne .type "container"}}Skipped to generate a container component{{end}}'
  - type: "add"
    path: "src/components/functional/{{.type}}{{if .parent}}/{{.parent}}{{end}}/{{.name | pascal}}.tsx"
    template: "Functional.tsx"
    skip: '{{if ne .type "functional"}}Skipped to generate a functional component{{end}}'
  - type: "add"
    path: "src/components/layout/{{.type}}{{if .parent}}/{{.parent}}{{end}}/{{.name | pascal}}Layout.tsx"
    template: "Layout.tsx"
    skip: '{{if ne .type "layout"}}Skipped to generate a layout component{{end}}'
  - type: "add"
    path: "src/components/page/{{.type}}{{if .parent}}/{{.parent}}{{end}}/{{.name | pascal}}Page.tsx"
    template: "Page.tsx"
    skip: '{{ if or (ne .type "page") (ne .pageType "page") }}Skipped to generate a page component{{end}}'
  - type: "add"
    path: "src/components/page/{{.type}}{{if .parent}}/{{.parent}}{{end}}/{{.name | pascal}}PageContents.tsx"
    template: "PageContents.tsx"
    skip: '{{ if or (ne .type "page") (ne .pageType "contents") }}Skipped to generate a page contents component{{end}}'
---

# Container.tsx

```tsx
{{if not .isSC}}
"use client";

{{end}}
export const {{.name | pascal}} = () => {
  return <div>{{.name | pascal}} Component</div>;
}
```

# Functional.tsx

```tsx
{{if not .isSC}}
"use client";

{{end}}
export const {{.name | pascal}} = () => {
  return <div>{{.name | pascal}} Component</div>;
}
```

# Layout.tsx

```tsx
{{if not .isSC}}
"use client";

{{end}}
export const {{.name | pascal}}Layout = ({ children }: { children: React.ReactNode }) => {
  return <div>{children}</div>;
}
```

# Page.tsx

```tsx
{{if not .isSC}}
"use client";

{{end}}
export const {{.name | pascal}}Page = () => {
  return <div>{{.name | pascal}}Page Component</div>;
}
```

# PageContents.tsx

```tsx
{{if not .isSC}}
"use client";

{{end}}
export const {{.name | pascal}}PageContents = () => {
  return <div>{{.name | pascal}}PageContents Component</div>;
}
```
