#! /usr/bin/env Rscript

library("rpart")

X <- read.table("zip.test.gz")
data <- list(image=as.matrix(X[,-1]), digit=as.factor(X[,1]))
rm(X)
varNames <- sapply(1:NCOL(data$image), function(col) paste("x",col-1,sep=""))
colnames(data$image) <- varNames
rm(varNames)

ctl <- rpart.control(
	minsplit=0,
	minbucket=0,
	cp=0,
	maxcompete=0,
	maxsurrogate=0,
	xval=5,
	maxdepth=30
)
prm <- list(
	split="information" # can be "gini" or "information"
)
t <- rpart(digit~image, data, parms=prm, control=ctl)
