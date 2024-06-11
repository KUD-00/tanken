resource "aws_db_instance" "postgres" {
  identifier              = var.identifier
  instance_class          = var.instance_class
  engine                  = var.engine
  engine_version          = var.engine_version
  username                = var.username
  password                = var.password
  db_subnet_group_name    = aws_db_subnet_group.rds_subnet_group.name
  vpc_security_group_ids  = var.vpc_security_group_ids
  
  allocated_storage       = var.allocated_storage
  storage_type            = var.storage_type

  tags                    = var.tags
}

resource "aws_db_subnet_group" "rds_subnet_group" {
  name       = var.subnet_group_name
  subnet_ids = var.subnet_ids

  tags = {
    Name = var.subnet_group_name
  }
}

