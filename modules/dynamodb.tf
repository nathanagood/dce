# Redbox Account table
# Tracks the status of AWS Accounts in our pool
resource "aws_dynamodb_table" "redbox_account" {
  name             = "RedboxAccount${title(var.namespace)}"
  read_capacity    = 5
  write_capacity   = 5
  hash_key         = "Id"
  stream_enabled   = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  global_secondary_index {
    name            = "AccountStatus"
    hash_key        = "AccountStatus"
    projection_type = "ALL"
    read_capacity   = 5
    write_capacity  = 5
  }

  global_secondary_index {
    name            = "GroupId"
    hash_key        = "GroupId"
    projection_type = "ALL"
    read_capacity   = 5
    write_capacity  = 5
  }

  server_side_encryption {
    enabled = true
  }

  # AWS Account ID
  attribute {
    name = "Id"
    type = "S"
  }

  # Status of the Account
  # May be one of:
  #   - ASSIGNED
  #   - READY
  #   - NOT_READY
  attribute {
    name = "AccountStatus"
    type = "S"
  }

  # Azure AD Group associated with the account
  attribute {
    name = "GroupId"
    type = "S"
  }

  tags = var.global_tags
  /*
  Other attributes:
  - LastModifiedOn (Integer, epoch timestamps)
  - CreatedOn (Integer, epoch timestamps)
  */
}

resource "aws_dynamodb_table" "redbox_account_assignment" {
  name             = "RedboxAccountAssignment${title(var.namespace)}"
  read_capacity    = 5
  write_capacity   = 5
  hash_key         = "AccountId"
  range_key        = "UserId"
  stream_enabled   = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  server_side_encryption {
    enabled = true
  }

  global_secondary_index {
    name            = "UserId"
    hash_key        = "UserId"
    projection_type = "ALL"
    read_capacity   = 5
    write_capacity  = 5
  }

  global_secondary_index {
    name            = "AssignmentStatus"
    hash_key        = "AssignmentStatus"
    projection_type = "ALL"
    read_capacity   = 5
    write_capacity  = 5
  }

  # AWS Account ID
  attribute {
    name = "AccountId"
    type = "S"
  }

  # Assignment status.
  # May be one of:
  # - ACTIVE
  # - FINANCE_LOCK
  # - RESET_LOCK
  # - DECOMMISSIONED
  attribute {
    name = "AssignmentStatus"
    type = "S"
  }

  # User ID
  attribute {
    name = "UserId"
    type = "S"
  }

  tags = var.global_tags
  /*
  Other attributes:
    - CreatedOn (Integer, epoch timestamps)
    - LastModifiedOn (Integer, epoch timestamps)
  */
}

