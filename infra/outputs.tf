output "cache_node_addresses" {
  value = aws_elasticache_cluster.data-cache.cache_nodes.*.address
}

output "cache_port" {
  value = aws_elasticache_cluster.data-cache.cache_nodes.*.port
}