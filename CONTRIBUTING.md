# 🚀 Contributing Manual

This project is open to contributions from all levels and any help is appreciated! 

> [!Tip]
>
> **For new contributors:** Take a look at [https://github.com/firstcontributions/first-contributions](https://github.com/firstcontributions/first-contributions) for a simple quick start guide on GitHub contributions.

## Development

### Prerequisites

Install the following tools for the project:

- [Go](https://go.dev/dl/) - API server (core)
- [gcc](https://gcc.gnu.org/install/) - Required for CGO in Go
- [Bun](https://bun.sh/docs/installation) - JavaScript runtime (dashboard, tracker)
- [Taskfile](https://taskfile.dev/installation/) - Global task runnner

We recommend using [mise-en-place](https://mise.jdx.dev/) as a convenient version manager tool for all the languages and tools required for this project.

### Setup

We use [Taskfile](https://taskfile.dev/) to run scripts and tasks to simplify the development process. Please refer to each sub-project's Taskfile for understanding the available tasks.

A full development setup can be done by running the following commands:


```bash [Terminal 1]
cd ./core
task dev -- start
```

```bash [Terminal 2]
cd ./dashboard
task dev
```

More details on the development setup can be found in the respective sub-project READMEs.

- [Core](./core/README.md)
- [Dashboard](./dashboard/README.md)
- [Tracker](./tracker/README.md)
