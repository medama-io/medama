# Schema

This document describes the current state of the schema for the database.

## Tables

### `users` - SQLite

Stores user data. Currently, only a single `admin` user is supported.

| Column         | Type                         | Description         |
| -------------- | ---------------------------- | ------------------- |
| `id`           | `TEXT PRIMARY KEY`           | Primary key         |
| `username`     | `TEXT NOT NULL`              | Username            |
| `password`     | `TEXT NOT NULL`              | Password            |
| `date_created` | `INTEGER NOT NULL`           | Date created (Unix) |
| `date_updated` | `INTEGER NOT NULL`           | Date updated (Unix) |
| `settings`     | `JSON NOT NULL DEFAULT '{}'` | User settings       |

#### Settings JSON Schema

```json
{
    "script_type": string, // 'default' or 'tagged-events'
    "language": string, // Only supports 'en' for now
}
```

### `views` - DuckDB

Stores page view event data.


| Column             | Type                   | Description                                                    |
| ------------------ | ---------------------- | -------------------------------------------------------------- |
| `bid`              | `TEXT PRIMARY KEY`     | Beacon ID used to link `load` and `unload` event data together |
| `hostname`         | `TEXT NOT NULL`        | Hostname                                                       |
| `pathname`         | `TEXT NOT NULL`        | Pathname                                                       |
| `is_unique_user`   | `BOOLEAN NOT NULL`     | Is unique visitor                                              |
| `is_unique_page`   | `BOOLEAN NOT NULL`     | Is unique visitor to specific page                             |
| `referrer_host`    | `TEXT`                 | Referrer hostname                                              |
| `referrer_group`   | `TEXT`                 | Referrer group name                                            |
| `country`          | `TEXT`                 | Country name                                                   |
| `language_base`    | `TEXT`                 | Base language                                                  |
| `language_dialect` | `TEXT`                 | Dialect language                                               |
| `ua_browser`       | `TEXT NOT NULL`        | Browser name                                                   |
| `ua_os`            | `TEXT NOT NULL`        | Operating system                                               |
| `ua_device_type`   | `TEXT NOT NULL`        | Device type                                                    |
| `utm_source`       | `TEXT`                 | UTM source                                                     |
| `utm_medium`       | `TEXT`                 | UTM medium                                                     |
| `utm_campaign`     | `TEXT`                 | UTM campaign                                                   |
| `duration_ms`      | `UINTEGER`             | Duration (ms)                                                  |
| `date_created`     | `TIMESTAMPTZ NOT NULL` | Date created                                                   |

### `events` - DuckDB

Stores custom properties event data.

| Column         | Type                   | Description                                                 |
| -------------- | ---------------------- | ----------------------------------------------------------- |
| `bid`          | `TEXT`                 | Beacon ID used to link to page view event                   |
| `batch_id`     | `TEXT NOT NULL`        | Batch ID used to link multiple properties of the same event |
| `group_name`   | `TEXT NOT NULL`        | Group name, typically the hostname                          |
| `name`         | `TEXT NOT NULL`        | Event key name                                              |
| `value`        | `TEXT NOT NULL`        | Event value                                                 |
| `date_created` | `TIMESTAMPTZ NOT NULL` | Date created                                                |
