---
name: "pages"
description: "Page Generator"
prompts:
  - type: "base"
    name: "name"
    message: "Type a page name"
    suffix: "* Use kebab-case"
    validate: "{{if eq (len .input) 0}}Name must have more than 0 characters.{{end}}"
  - type: "base"
    name: "parent"
    message: "Type a path of the parent folder if needed"
    suffix: "* [Optional] Skip to enter"
    defaultValues: ""
  - type: "confirm"
    name: "requiredLayout"
    message: "Do you need a layout?"
    defaultValues: true
  - type: "confirm"
    name: "requiredTemplate"
    message: "Do you need a template?"
    defaultValues: false
actions:
  - type: "add"
    path: "src/app{{if .parent}}/{{.parent}}{{end}}/{{.name}}/page.tsx"
    template: "{{.name | kebab}}/page.tsx"
  - type: "add"
    path: "src/app{{if .parent}}/{{.parent}}{{end}}/{{.name}}/layout.tsx"
    template: "{{.name | kebab}}/layout.tsx"
    skip: "{{if not .requiredLayout}}Skipped to generate a layout file.{{end}}"
  - type: "add"
    path: "src/app{{if .parent}}/{{.parent}}{{end}}/{{.name}}/template.tsx"
    template: "{{.name | kebab}}/template.tsx"
    skip: "{{if not .requiredTemplate}}Skipped to generate a template file.{{end}}"
---

# `{{.name | kebab}}/page.tsx`

```tsx
export default function {{.name | pascal}}Index() {
    return (
        <div>{{.name | pascal}}</div>
    );
};
```

# `{{.name | kebab}}/layout.tsx`

```tsx
export default function {{.name | pascal}}Layout({
  children
}: {
  children: React.ReactNode,
}) {
  return (
    <section>
      {children}
    </section>
  );
};
```

# `{{.name | kebab}}/template.tsx`

```tsx
export default function {{.name | pascal}}Template({
  children
}: {
  children: React.ReactNode
}) {
  return (
    <div>
      {children}
    </div>
  );
};
```
