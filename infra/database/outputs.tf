output "postgres_instance_id" {
  value       = aws_db_instance.postgres.id
  description = "The ID of the RDS instance."
}

output "db_subnet_group_name" {
  value       = aws_db_subnet_group.rds_subnet_group.name
  description = "The name of the subnet group used by the RDS instance."
}
