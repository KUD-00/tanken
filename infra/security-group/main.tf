resource "aws_security_group" "data_cache_security_group" {
  name   = "data_cache_security_group"
  vpc_id = var.vpc_id

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = var.cache_ingress_cidr_blocks
  }
}

resource "aws_security_group" "rds_security_group" {
  name   = "rds-security-group"
  vpc_id = var.vpc_id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = var.rds_ingress_cidr_blocks
  }
}