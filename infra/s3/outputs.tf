output "s3_bucket_id" {
  value       = aws_s3_bucket.picture_bucket.id
  description = "The ID of the S3 bucket."
}

output "s3_access_point_id" {
  value       = aws_s3_access_point.picture_bucket_access_point.id
  description = "The ID of the S3 access point."
}
