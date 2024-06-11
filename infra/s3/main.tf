resource "aws_s3_bucket" "picture_bucket" {
  bucket = var.bucket_name

  tags = var.bucket_tags
}

resource "aws_s3_access_point" "picture_bucket_access_point" {
  bucket = aws_s3_bucket.picture_bucket.id
  name   = var.access_point_name

  vpc_configuration {
    vpc_id = var.vpc_id
  }
}

resource "aws_s3_bucket_policy" "picture_bucket_policy" {
  bucket = aws_s3_bucket.picture_bucket.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "arn:aws:s3:::${var.bucket_name}/*"
      },
    ]
  })
}
