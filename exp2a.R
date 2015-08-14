#! /usr/bin/env Rscript

library("tree")

X <- read.table("zip.train.gz")
data <- list(image=as.matrix(X[,-1]), digit=as.factor(X[,1]))
rm(X)
varNames <- sapply(1:NCOL(data$image), function(col) paste("x",col-1,sep=""))
colnames(data$image) <- varNames
rm(varNames)

ctl <- tree.control(NROW(data$image),
	mincut=1,
	minsize=2,
	mindev=0
)
splitMethod <- "deviance" # can be "gini" or "deviance"
t <- tree(digit~image, data, control=ctl, split=splitMethod)
