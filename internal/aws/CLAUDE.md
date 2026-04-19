# aws

AWS SDK v2 client wrappers. One file per AWS service. All clients take `aws.Config` from session.go.

Pattern: each client struct wraps the SDK client; methods return domain types (not SDK types). Errors are wrapped with `fmt.Errorf("action: %w", err)`.

Security: SSM WithDecryption must always be false. SecureString values are masked with `"****"`.
