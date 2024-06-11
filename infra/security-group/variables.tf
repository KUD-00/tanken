variable "vpc_id" {
  type        = string
  description = "The VPC ID where the security groups will be created."
}

variable "cache_ingress_cidr_blocks" {
  type        = list(string)
  description = "List of CIDR blocks allowed to access the Elasticache."
  default     = []
}

variable "rds_ingress_cidr_blocks" {
  type        = list(string)
  description = "List of CIDR blocks allowed to access the RDS instance."
  default     = []
}
