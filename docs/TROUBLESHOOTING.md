# 🚨 Troubleshooting Guide

This guide helps you resolve common issues when using this CLI.

## 🔐 Authentication Failures

**Symptom:** `Authentication failed` or `401 Unauthorized`

**Solutions:**

1. **Verify Credentials:**

   ```bash
   # Test token generation
   wnc generate token --username admin --password your_password
   ```

2. **Check Token Format:**

   ```bash
   # Ensure proper encoding
   echo "YWRtaW46cGFzc3dvcmQ=" | base64 -d
   # Should output: admin:password
   ```

3. **Validate Controller URL:**

   ```bash
   # Test basic connectivity
   curl -k https://your-controller.example.com/restconf/
   ```

## ⏱️ Connection Timeouts

**Symptom:** `Connection timeout` or `Network unreachable`

**Solutions:**

1. **Increase Timeout:**

   ```bash
   wnc show ap --controllers "wnc.example.com:token" --timeout 120
   ```

2. **Check Network Connectivity:**

   ```bash
   # Test basic network access
   ping your-controller.example.com
   telnet your-controller.example.com 443
   ```

3. **Verify DNS Resolution:**

   ```bash
   nslookup your-controller.example.com
   ```

## 🔒 Certificate Errors

**Symptom:** `Certificate verification failed` or `SSL handshake failed`

**Solutions:**

1. **Development Environment:**

   > [!CAUTION]
   > The `--insecure` flag disables certificate verification and should only be used in development environments with self-signed certificates.

   ```bash
   # Temporary workaround for self-signed certificates
   wnc show ap --controllers "wnc-dev.local:token" --insecure
   ```

2. **Production Environment:**

   ```bash
   # Install proper CA certificates
   sudo curl -o /usr/local/share/ca-certificates/company-ca.crt \
     https://ca.company.com/company-ca.crt
   sudo update-ca-certificates
   ```

## 📭 Empty Results

**Symptom:** No data returned or empty tables

**Solutions:**

1. **Verify API Permissions:**

   ```bash
   # Test with basic overview command
   wnc show overview --controllers "wnc.example.com:token"
   ```

2. **Check Filter Parameters:**

   ```bash
   # Remove filters to see all data
   wnc show client --controllers "wnc.example.com:token"
   # vs filtered
   wnc show client --controllers "wnc.example.com:token" --ssid "NonExistentSSID"
   ```

3. **Validate Controller Status:**

   ```bash
   # Check if controller is responding
   curl -k -H "Authorization: Basic $TOKEN" \
     https://your-controller.example.com/restconf/data/Cisco-IOS-XE-wireless-access-point-oper:access-point-oper-data
   ```
