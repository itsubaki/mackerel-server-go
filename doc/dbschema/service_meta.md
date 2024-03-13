# service_meta

## Description

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE TABLE `service_meta` (
  `org_id` varchar(16) NOT NULL,
  `service_name` varchar(16) NOT NULL,
  `namespace` varchar(128) NOT NULL,
  `meta` text,
  PRIMARY KEY (`org_id`,`service_name`,`namespace`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
```

</details>

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| org_id | varchar(16) |  | false |  |  |  |
| service_name | varchar(16) |  | false |  |  |  |
| namespace | varchar(128) |  | false |  |  |  |
| meta | text |  | true |  |  |  |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| PRIMARY | PRIMARY KEY | PRIMARY KEY (org_id, service_name, namespace) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| PRIMARY | PRIMARY KEY (org_id, service_name, namespace) USING BTREE |

## Relations

![er](service_meta.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)