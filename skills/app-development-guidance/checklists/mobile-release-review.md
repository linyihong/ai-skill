# Mobile Release Review Checklist

- Release build uses production-safe endpoint and logging configuration.
- Cleartext traffic is disabled unless a scoped exception is documented.
- Obfuscation/minification/symbol handling matches the product risk and support plan.
- Secrets and internal URLs are not present in the artifact.
- Crash reporting and analytics are reviewed for sensitive data.
- Pinning, if enabled, has a rotation and incident response plan.
- Known residual risks are documented in the project repository.
