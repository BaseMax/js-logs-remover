# js-logs-remover in Go

`js-logs-remover` is a command-line utility written in Go that removes specified `console` logging methods (e.g., `console.log()`, `console.error()`) from JavaScript/TypeScript files within a directory. This tool is useful for cleaning up log statements before deploying or publishing code.

## Features

- Removes specified `console` methods from JavaScript, TypeScript, JSX, and TSX files.
- Processes files recursively within directories, excluding certain common directories such as `node_modules`, `.git`, `dist`, and `build`.
- Supports both specified log methods and a global "all" option to remove all defined `console` methods.

## Installation

To install `js-logs-remover`, you must have Go installed on your machine.

1. Clone the repository:

   ```bash
   git clone https://github.com/BaseMax/js-logs-remover-go.git
   ```

2. Navigate to the directory:

   ```bash
   cd js-logs-remover-go
   ```

3. Build the Go application:

   ```bash
   go build
   ```

After building, you will have a js-logs-remover binary in the project directory.

## Usage

Run the tool from the command line:

```bash
$ js-logs-remover [path] [log-methods]
```

`[path]`: The directory or file path where you want to remove the logs from. If no path is provided, it defaults to the current directory (.).

`[log-methods]`: A comma-separated list of console methods you want to remove (e.g., log,error,warn). If you specify all, it will remove all known console methods.

## Example Usage

Remove specific console methods: To remove console.log and console.warn from all .js, .ts, .jsx, and .tsx files within the src directory:

```bash
$ js-logs-remover src log,warn
```

Remove all console methods: To remove all console methods (like console.log, console.error, etc.) from all files in the src directory:

```bash
$ js-logs-remover src all
```

Process the current directory: To process the current directory and remove console.log:

```bash
$ js-logs-remover . log
```

### Excluded Directories

The following directories are automatically excluded from processing:

- node_modules
- .git
- dist
- build

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Copyright

© 2025 Max Base. All rights reserved.
