resource "aws_dynamodb_table" "reddit_submissions" {
    name = var.submissions_table_name
    billing_mode = "PAY_PER_REQUEST"
    hash_key = "submission_id"

    attribute {
        name = "submission_id"
        type = "S"
    }

    ttl {
        enabled = true
        attribute_name = "expiration_date"
    }
}