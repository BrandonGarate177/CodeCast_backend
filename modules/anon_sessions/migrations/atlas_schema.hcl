schema "public" {
  comment = "standard public schema"
}

table "anon_sessions" {
  schema = schema.public
  column "id" {
    type = text
    null = false
  }
  column "session_code" {
    type = text
    null = false
  }
  column "creator_token" {
    type = text
    null = false
  }
  column "display_name" {
    type = text
    null = false
  }
  column "is_active" {
    type = boolean
    null = false
  }
  column "created_at" {
    type = timestamp
    null = false
  }
  primary_key {
    columns = [column.id]
  }
}

table "anon_participants" {
  schema = schema.public
  column "id" {
    type = text
    null = false
  }
  column "session_id" {
    type = text
    null = false
  }
  column "display_name" {
    type = text
    null = false
  }
  column "joined_at" {
    type = timestamp
    null = false
  }
  primary_key {
    columns = [column.id]
  }
}

table "anon_snippets" {
  schema = schema.public
  column "id" {
    type = text
    null = false
  }
  column "session_id" {
    type = text
    null = false
  }
  column "file_name" {
    type = text
    null = false
  }
  column "content" {
    type = text
    null = false
  }
  column "pushed_at" {
    type = timestamp
    null = false
  }
  primary_key {
    columns = [column.id]
  }
}