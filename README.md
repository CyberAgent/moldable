# `moldable`

A CLI tool for code generation under the influence of [plop](https://plopjs.com) and [scaffdog](https://scaff.dog). It has the best of each tool, yet it is faster and more flexible.
See the [benchmarks](#benchmarks) for more information.

## Installation

```bash
npm install -g moldable
```

## Quick Start

Create a `.moldable` folder in the root directory of your application, and create Markdown file within it shown below.

```md

```

Execute `moldable` in the root directory of your application, and select the Generator you want to use.

```bash

```

Answer the questions that appear, and the code will be generated.

```bash

```

## Advanced Usage

You can describe a generator for generating any code.
Create a `/.moldable` folder in the root directory of your application,
and create Markdown files within it.

```
.
└──  .moldable
    ├── foo.md
    └── bar.md
```

The name of the Markdown file corresponds to the Generator name displayed when you launch moldable.
When you run moldable, an interactive question starts according to the selected Generator, and when all are answered, code is generated.
This Generator is divided into metadata (in YAML format) and content (in Markdown format).
The structure of the metadata is as follows. Some properties can use the syntax of Go's [text/template](https://pkg.go.dev/text/template).

| Name        | Description                                                                   | Required           | Type             | Example                                                                                                      |
| ----------- | ----------------------------------------------------------------------------- | ------------------ | ---------------- | ------------------------------------------------------------------------------------------------------------ |
| name        | Template name (must match the Markdown file name)                             | :white_check_mark: | String           | "pages"                                                                                                      |
| description | Template description                                                          | :white_check_mark: | String           | "Page Generator"                                                                                             |
| prompts     | Questions to be asked sequentially after selecting a template                 | :white_check_mark: | Array of mapping | <pre>type: "base"<br>name: "name"<br>message: "Type a page name"<br></pre>                                   |
| actions     | Settings for generating code using the values of the answers to the questions | :white_check_mark: | Array of mapping | <pre>type: "add"<br>path: "src/app/{{.name}}/page.tsx"<br>template: "{{.name \| pascal}}/page.tsx"<br></pre> |

The structure of each element of the prompts value is as follows.

| Name          | Description                                                                                                                                            | Required                                  | Type                                             |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------- | ------------------------------------------------ |
| type          | Format of the question<br>・base: Text input<br>・select: Single selection<br>・multiselect: Multiple selection<br>・confirm: Y/N binary choice        | :white_check_mark:                        | "base" \| "select" \| "multiselect" \| "confirm" |
| name          | Variable name to store the result of the question. In places where text/template is valid, it can be referenced with a . start.                        | :white_check_mark:                        | String                                           |
| message       | Content of the question                                                                                                                                | :white_check_mark:                        | String                                           |
| prefix        | Supplement to the question, displayed before the question                                                                                              |                                           | String                                           |
| suffix        | Supplement to the question, displayed after the question                                                                                               |                                           | String                                           |
| when          | Place to skip the question. The results of the questions up to the current question can be referenced.                                                 |                                           | String(text/template)                            |
| validate      | Place to validate the question. The value of the question being entered can be referenced with .input.                                                 |                                           | String(text/template)                            |
| defaultValues | Default value of the result of the question<br>Data type<br>・base: String<br>・select: String<br>・multiselect: Array of String<br>・confirm: Boolean |                                           | Any                                              |
| choises       | Choices for select and multiselect                                                                                                                     | Required if type is select or multiselect | Array of String                                  |

The structure of each element of the actions value is as follows.

| Name     | Description                                                               | Required                   | Type                  |
| -------- | ------------------------------------------------------------------------- | -------------------------- | --------------------- |
| type     | Type of code generation                                                   | :white_check_mark:         | "add" \| "append"     |
| path     | Destination of the code output                                            | :white_check_mark:         | String(text/template) |
| template | Template name for code generation (must match the heading in the content) | :white_check_mark:         | String(text/template) |
| skip     | Place to skip code generation                                             |                            | String(text/template) |
| pattern  | Destination to add code in append                                         | Required if type is append | String                |

In addition, for variables referenced within text/template, you can convert to the following cases using the format .var | case.

| Name     | Case          |
| -------- | ------------- |
| camel    | camelCase     |
| snake    | snake_case    |
| kebab    | kebab-case    |
| dash     | dash-case     |
| dot      | dot.case      |
| proper   | ProperCase    |
| pascal   | PascalCase    |
| lower    | lower case    |
| sentence | Sentence case |
| constant | CONSTANT_CASE |
| title    | Title Case    |

> [!NOTE]
>
> By attaching OnlyAlphanumeric to the end of each case modifier (e.g. pascalOnlyAlphanumeric), if the first character of the string is **not alphanumeric**, it is excluded and the case is converted.
> For example, [projectId] is converted to ProjectId with pascalOnlyAlphanumeric.

The content consists of headings and code blocks as follows.

````md
---
actions:
  - type: "add"
    template: "page.tsx"
---

# page.tsx

```tsx
export const Page = (children: { children: React.ReactNode }) => <div>{children}</div>;
```
````

As mentioned earlier, the template of any element contained in actions needs to match the heading in the content.
The syntax of text/template can also be used for headings.

## Benchmarks

### Environment

- OS: macOS Monterey 12.6.3
- CPU: Apple M1 Max
- Memory: 64GB
- Node.js: v20.11.1
- pnpm: 8.1.1

We measured the time it took to generate a simple React component with each package.

```tsx
import React from "react";

export const Page = (children: { children: React.ReactNode }) => <div>{children}</div>;
```

### Result

| Package         | Version | Time (ms) |
| --------------- | ------- | --------- |
| plop(turbo gen) |         |           |
| scaffdog        |         |           |
| hygen           |         |           |
| **moldable**    |         |           |

It is the fastest. This is because it is written in Go and executed by node after compiling it into a binary.

## Third Party License

Parts of the code from the following projects have been borrowed and modified for use.

- [Plop.js(MIT License)](https://github.com/plopjs/plop/blob/b797aafa8ae90adb59462114412505f9ec0c5bb1/LICENSE.md)
- [Scaffdog(MIT License)](https://github.com/scaffdog/scaffdog/blob/9221585670365f67f6c439aea33a09ca79ae499c/LICENSE)

The list of licenses is available at [THIRD_PARTY_LICENSE](https://github.com/cyberagent-oss/moldable/blob/main/THIRD_PARTY_LICENSE.md).

## License

[MIT License](./LICENSE)

## Developers

- Shuta Hirai - [GitHub](https://github.com/shuta13)
- Satoshi Kotake - [GitHub](https://github.com/cotapon)
- Takuma Shibuya - [GitHub](https://github.com/sivchari)

## Copyright

CyberAgent, Inc. All rights reserved.
