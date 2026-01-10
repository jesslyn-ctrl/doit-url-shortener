terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# ====================
# Provider
# ====================
provider "aws" {
  region = "ap-southeast-1"
}

# ====================
# IAM Role (Least Privilege)
# ====================
resource "aws_iam_role" "lambda_role" {
  name = "doit-url-shortener-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_policy" "lambda_policy" {
  name = "doit-url-shortener-lambda-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      # CloudWatch Logs
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "*"
      },

      # DynamoDB access (scoped to one table)
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem"
        ]
        Resource = aws_dynamodb_table.short_urls.arn
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_attach" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.lambda_policy.arn
}

# ====================
# Managed Storage (DynamoDB)
# ====================
resource "aws_dynamodb_table" "short_urls" {
  name         = "doit-url-shortener"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "code"

  attribute {
    name = "code"
    type = "S"
  }
}

# ====================
# Serverless Compute (Lambda)
# ====================
resource "aws_lambda_function" "url_shortener" {
  function_name = "doit-url-shortener"
  role          = aws_iam_role.lambda_role.arn

  runtime = "provided.al2"
  handler = "bootstrap"

  # Placeholder artifact (not required to actually exist for validation)
  filename         = "lambda.zip"

  timeout     = 10
  memory_size = 256

  environment {
    variables = {
      APP_NAME        = "doit-url-shortener"
      APP_PORT        = "8081"
      URL_TTL_SECONDS = "86400"
      TABLE_NAME      = aws_dynamodb_table.short_urls.name
    }
  }
}
