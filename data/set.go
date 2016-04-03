package data

// Set implements an abstract interface to represent a test data set,
// consisting of training data for setting up the method and test data
// for assessment.
type Set interface {
	GetName() string
	TrainingData() (data *Data, err error)
	TestData() (data *Data, err error)
}
