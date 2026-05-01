# Kotlin And Java Notes

Use this for Android Kotlin/Java implementation details.

- Review exported components, intent extras, and deeplink handlers for trust boundary mistakes.
- Keep debug interceptors, logging interceptors, and staging endpoints out of release builds.
- Prefer platform-backed storage for secrets and avoid custom crypto wrappers unless justified.
- Treat OkHttp/HTTP client configuration as auditable release behavior.

See [`../platforms/mobile/android.md`](../platforms/mobile/android.md) for Android platform guidance.
