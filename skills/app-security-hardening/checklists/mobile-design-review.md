# Mobile Design Review Checklist

- What data, token, or action is high risk in this feature?
- Which controls are server-owned, client-owned, build-owned, and monitoring-owned?
- Can the server reject modified clients or replayed requests?
- Is local/offline data necessary, scoped, and expired?
- Are logs, analytics, and crash reports safe by design?
- Are platform-specific risks documented for Android, iOS, Flutter, or React Native?
- What test or review step proves the control exists?
