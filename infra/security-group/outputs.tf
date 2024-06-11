output "data_cache_security_group_id" {
  value       = aws_security_group.data_cache_security_group.id
  description = "The ID of the security group for Elasticache."
}

output "rds_security_group_id" {
  value       = aws_security_group.rds_security_group.id
  description = "The ID of the security group for RDS."
}
