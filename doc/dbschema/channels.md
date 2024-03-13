# channels

## Description

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE TABLE `channels` (
  `org_id` varchar(16) NOT NULL,
  `id` varchar(16) NOT NULL,
  `name` varchar(16) NOT NULL,
  `type` enum('email','slack','webhook') NOT NULL,
  `url` text,
  `enabled_graph_image` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
```

</details>

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| org_id | varchar(16) |  | false |  |  |  |
| id | varchar(16) |  | false |  |  |  |
| name | varchar(16) |  | false |  |  |  |
| type | enum('email','slack','webhook') |  | false |  |  |  |
| url | text |  | true |  |  |  |
| enabled_graph_image | tinyint(1) | 1 | false |  |  |  |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| PRIMARY | PRIMARY KEY | PRIMARY KEY (id) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| PRIMARY | PRIMARY KEY (id) USING BTREE |

## Relations

![er](channels.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)