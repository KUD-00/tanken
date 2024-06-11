variable "bucket_name" {
  type        = string
  description = "The name of the S3 bucket."
}

variable "access_point_name" {
  type        = string
  description = "The name of the S3 access point."
}

variable "vpc_id" {
  type        = string
  description = "The VPC ID to associate with the S3 access point."
}

variable "bucket_tags" {
  type        = map(string)
  description = "A mapping of tags to assign to the S3 bucket."
}
