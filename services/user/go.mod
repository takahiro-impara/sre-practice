module github.com/lot-koichi/sre-skill-up-project/services/user

go 1.24.6

require (
	github.com/go-chi/chi/v5 v5.2.3
	github.com/go-chi/render v1.0.3
	github.com/go-playground/validator/v10 v10.27.0
	github.com/golang-migrate/migrate/v4 v4.19.0
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.11.1
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.36.0
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/lot-koichi/sre-skill-up-project/pkg => ../../pkg
