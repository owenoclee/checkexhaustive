# checkexhaustive

Checkexhaustive is a simple analyzer for Go that ensures any struct literal
labelled with `//check:exhaustive` (on the line preceding it) has its fields
filled exhaustively.

```go
package example

type Contact struct {
	Name  string
	Email string
	Phone string
}

func NewContact(name string) Contact {
	//check:exhaustive
	return Contact{
		Name: name,
	}
}
```

```
$ checkexhaustive ./example
example/example.go:11:9: Contact is missing fields: Email, Phone
```

## Installation

```
go install github.com/owenoclee/checkexhaustive/cmd/checkexhaustive
```

## Usage

```
checkexhaustive [packages]
```

or, as part of [go vet](https://pkg.go.dev/cmd/vet):

```
go vet -vettool=$(which checkexhaustive) [packages]
```

or, integrated into Visual Studio Code with the following `settings.json`:

```jsonc
{
	"go.vetOnSave": "workspace",
	"go.vetFlags": [
	    "-vettool=/path/to/checkexhaustive"
	],
	"go.languageServerExperimentalFeatures": {
	    "diagnostics": false, // https://github.com/golang/vscode-go/issues/1109
	},
}
```
