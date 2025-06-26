# 🔐 Security Considerations

This document outlines security best practices when using the `wnc`.

## 🔒 TLS Certificate Verification

By default, the CLI enforces strict TLS certificate verification:

```bash
# Secure connection (default)
wnc show ap --controllers "wnc.example.com:token"

# Skip verification only for development/testing
wnc show ap --controllers "wnc-dev.local:token" --insecure
```

> [!WARNING]
> The `--insecure` flag disables TLS certificate verification and should only be used in trusted development environments.

## 🔑 Authentication Token Security

### ✅ Best Practices

1. **Environment Variables**: Store tokens in environment variables, never in scripts:

   ```bash
   export WNC_CONTROLLERS="wnc.example.com:$(wnc generate token -u admin -p $PASSWORD)"
   ```

2. **Token Rotation**: Regenerate tokens regularly:

   ```bash
   # Automated token refresh
   NEW_TOKEN=$(wnc generate token --username admin --password "$PASSWORD")
   export WNC_CONTROLLERS="wnc.example.com:$NEW_TOKEN"
   ```

3. **Secure Storage**: Use secure credential management:

   ```bash
   # Example with macOS Keychain
   PASSWORD=$(security find-generic-password -a admin -s wnc-password -w)
   TOKEN=$(wnc generate token --username admin --password "$PASSWORD")
   ```

### ❌ What to Avoid

- Never log authentication tokens
- Don't store tokens in source code
- Avoid hardcoding credentials in scripts
- Don't share tokens between environments

## 🌐 Network Security

### 🔥 Firewall Considerations

- HTTPS traffic on port 443 (default)
- Outbound connections to controller management interfaces
- Consider VPN access for production environments

### 🚪 Access Control

- Use dedicated service accounts with minimal privileges
- Implement read-only access where possible
- Monitor API access logs on controllers
- Regularly audit user permissions

## 🏭 Production Deployment

### 🏗️ Environment Isolation

```bash
# Development
export WNC_CONTROLLERS="wnc-dev.local:$DEV_TOKEN"

# Staging
export WNC_CONTROLLERS="wnc-staging.company.com:$STAGING_TOKEN"

# Production
export WNC_CONTROLLERS="wnc-prod.company.com:$PROD_TOKEN"
```

### 📊 Monitoring and Logging

- Monitor CLI usage patterns
- Log command execution for audit trails
- Set up alerts for authentication failures
- Track API call volumes

---

**Back to:** [CLI Reference](CLI_REFERENCE.md)
