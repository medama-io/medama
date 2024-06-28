<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./.github/images/banner-dark.svg">
    <source media="(prefers-color-scheme: light)" srcset="./.github/images/banner-light.svg">
    <img alt="Medama: Cookie-free privacy-focused website analytics." src="./.github/images/banner-light.svg">
  </picture>
  <br>
  <a href="https://oss.medama.io/introduction">Website</a> |
  <a href="https://oss.medama.io/deployment/installation">Installation</a> |
  <a href="https://demo.medama.io">Demo</a>
</p>

## Overview

Medama Analytics is an open-source project dedicated to providing self-hostable, cookie-free website analytics. With a lightweight tracker of less than 1KB, it aims to offer useful analytics while prioritising user privacy.

<p align="center">
    <a href="https://demo.medama.io" target="_blank">
        <img src="./.github/images/demo.png" alt="Demo Screenshot" width="70%" height="70%">
    </a>
</p>

### Features

- 📊 **Real-Time Analytics:** Gain instant insights with real-time analytics, allowing you to monitor website performance and user interactions as they happen.
- 🔒 **Privacy-Focused:** Using a lightweight tracker of less than 1KB that operates without relying on cookies, IP addresses, or additional identifiers, this project ensures compliance with GDPR, PECR, and other privacy regulations.
- 💼 **Self-Hostable:** Using embedded databases such as SQLite and DuckDB, this approach has a single-binary straightforward setup with no external dependencies. It's a lightweight solution that can efficiently run on a VM with 256MB of memory for most small websites.

## License

The `/core` and `/dashboard` directory is licensed under the Apache License 2.0. See the core [LICENSE](./core/LICENSE) and dashboard [LICENSE](./dashboard/LICENSE) for more information.

The `/tracker` directory is licensed under the MIT License. See [LICENSE](./tracker/LICENSE) for more information.
