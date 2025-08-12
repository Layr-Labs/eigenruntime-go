# RuntimeSpec Template Guide

## Structure

RuntimeSpec defines components as a flexible map with user-defined keys:

```yaml
apiVersion: eigenruntime.io/v1alpha1
kind: Runtime
name: my-runtime
version: v1.0.0
spec:
  <component-name>:
    registry: <container-registry-url>
    digest: <sha256-digest>
    command: [optional-command-array]
    env: [optional-env-declarations]
    resources: {optional-resource-config}
```

## Field Reference

### Top-Level Fields

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `apiVersion` | ✓ | string | API version |
| `kind` | ✓ | string | Resource type (`Runtime`, `Hourglass`) |
| `name` | ✓ | string | Runtime instance name |
| `version` | ✓ | string | Runtime version |
| `spec` | ✓ | map[string]Component | Component definitions (min 1) |

### Component Fields

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `registry` | ✓ | string | Container registry URL |
| `digest` | ✓ | string | SHA256 digest |
| `command` | ✗ | []string | Override container command |
| `env` | ✗ | []EnvVar | Environment variable declarations |
| `resources` | ✗ | Resources | Resource configuration |

### Environment Variables

**Declarative only** - defines expected variables, not values:

| Field | Required | Type | Description                    |
|-------|----------|------|--------------------------------|
| `name` | ✓ | string | Variable name                  |
| `type` | ✗ | string | `secret`,`, or `runtime`       |
| `required` | ✗ | bool | Must be provided at deployment |

### Resources

| Field | Type | Description |
|-------|------|-------------|
| `teeEnabled` | bool | Enable TEE |

## Example

```yaml
apiVersion: eigenruntime.io/v1alpha1
kind: Runtime
name: distributed-runtime
version: v1.0.0
spec:
  executor:
    registry: ghcr.io/example/executor
    digest: sha256:bbb222
    
  performer:
    registry: ghcr.io/example/performer
    digest: sha256:ccc333
    env:
      - name: DATABASE_URL
        type: secret
        required: true
```

## Component Naming

User-defined keys in the `spec` map. Common patterns:
- **Role-based**: `aggregator`, `executor`, `performer`
- **Service-based**: `api-server`, `worker`, `scheduler`
- **Layer-based**: `frontend`, `backend`, `database`

## Validation Rules

**Required fields:**
- Top-level: `apiVersion`, `kind`, `name`, `version` (non-empty)
- `spec` must contain at least one component
- Each component: `registry` and `digest` required
- Environment variables: `name` required (non-empty)

**Not validated:**
- Component names (any valid YAML key)
- Environment variable types (any string)
- Number of components

## Common Errors

| Error | Fix |
|-------|-----|
| `apiVersion is required` | Add non-empty `apiVersion` field |
| `spec must contain at least one component` | Add component under `spec` |
| `registry is required` | Add `registry` field to component |
| `digest is required` | Add SHA256 digest (not tags) |
| `environment variable name cannot be empty` | Add `name` to env var |

## Best Practices

1. **Use digests** not tags (`sha256:abc123` not `:latest`)
2. **Mark critical variables** with `required: true`
3. **Use descriptive component names** reflecting their role
4. **Remember `env` is declarative** - no actual values in spec
5. **Use camelCase** for resources (`teeEnabled` not `tee_enabled`)

## Notes

- Component names become deployment identifiers
- Environment values injected at deployment based on `type`
- `command` overrides container's default command
- Resource constraints beyond `teeEnabled` handled by deployment platform