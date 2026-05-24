# 🚀 Contributing Manual

This project is open to contributions from all levels and any help is appreciated! 

> [!Tip]
>
> **For new contributors:** Take a look at [https://github.com/firstcontributions/first-contributions](https://github.com/firstcontributions/first-contributions) for a simple quick start guide on GitHub contributions.

## Development

### Prerequisites

Install the following tools for the project:

- [Go](https://go.dev/dl/) - API server (core)
- [gcc](https://gcc.gnu.org/install/) - Required for CGO in Go (core)
- [Bun](https://bun.sh/docs/installation) - JavaScript runtime (dashboard, tracker)
- [mise-en-place](https://mise.jdx.dev/) - Tool and workflow manager

We recommend using [mise-en-place](https://mise.jdx.dev/) as a convenient version manager tool for all the languages and tools required for this project.

### Setup

We use mise tasks to run project workflows. Run `mise install` from the repository root first, then use the package tasks for normal development.

A full development setup can be done by running the following commands in separate terminals:


```bash [Terminal 1]
mise run core:dev # Start API server
```

```bash [Terminal 2]
mise run dashboard:dev # Start dashboard dev server with HMR
```

More details on the development setup can be found in the respective sub-project READMEs.

- [Core](./core/README.md)
- [Dashboard](./dashboard/README.md)
- [Tracker](./tracker/README.md)
