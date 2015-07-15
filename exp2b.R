#! /usr/bin/env Rscript

library("rpart")
X <- read.table("zip.test.gz")
data <- list(image=as.matrix(X[,-1]),digit=as.factor(X[,1]))
ctl <- rpart.control(minbucket=1)
t <- rpart(digit~image, data, parms=list(split="gini"))
