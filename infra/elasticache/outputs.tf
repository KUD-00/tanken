output "elasticache_cluster_id" {
  value       = aws_elasticache_cluster.data_cache.id
  description = "The ID of the Elasticache cluster."
}

output "elasticache_subnet_group_name" {
  value       = aws_elasticache_subnet_group.data_cache_subnet_group.name
  description = "The name of the subnet group used by the Elasticache cluster."
}
