# Security Policy

## Reporting a Vulnerability
Please use GitHub private security advisories (preferred). Do not open a public issue for sensitive reports.

## Notes
portik may shell out to system tools (e.g., `ss`, `lsof`, `ps`, `docker`) and parses their output.
Treat external command output as untrusted.