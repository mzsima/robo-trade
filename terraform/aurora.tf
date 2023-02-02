resource "aws_rds_cluster" "robo_trade1_rds_cluster" {
  cluster_identifier = "robotrade1"
  engine             = "aurora-mysql"
  engine_mode        = "provisioned"

  master_username = "admin"
  master_password = "password"
  engine_version     = "8.0.mysql_aurora.3.02.0"

  vpc_security_group_ids = [aws_security_group.robo_trade1_rds_security_group.id]

  serverlessv2_scaling_configuration {
    min_capacity = 1
    max_capacity = 1
  }

  skip_final_snapshot = true
}

resource "aws_rds_cluster_instance" "example" {
  cluster_identifier = aws_rds_cluster.robo_trade1_rds_cluster.id
  identifier         = "robo-trade1-serverless-instance"

  engine                  = aws_rds_cluster.robo_trade1_rds_cluster.engine
  engine_version          = aws_rds_cluster.robo_trade1_rds_cluster.engine_version
  instance_class          = "db.serverless"

  publicly_accessible = true
}

resource "aws_security_group" "robo_trade1_rds_security_group" {
  name = "robo_trade1_rds_security_group"

  ingress {
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}