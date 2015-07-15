#! /usr/bin/env Rscript

library("tree")
X <- read.table("zip.test.gz")
data <- list(image=as.matrix(X[,-1]),digit=as.factor(X[,1]))
t <- tree(digit~image, data, split="deviance")
