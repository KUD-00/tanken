resource "aws_elasticache_cluster" "data_cache" {
  cluster_id           = var.cluster_id
  engine               = var.engine
  node_type            = var.node_type
  num_cache_nodes      = var.num_cache_nodes
  parameter_group_name = var.parameter_group_name
  engine_version       = var.engine_version
  port                 = var.port

  security_group_ids   = var.security_group_ids
  subnet_group_name    = aws_elasticache_subnet_group.data_cache_subnet_group.name
}

resource "aws_elasticache_subnet_group" "data_cache_subnet_group" {
  name       = var.subnet_group_name
  subnet_ids = var.subnet_ids

  description = "Elasticache subnet group for caching layer."
}
