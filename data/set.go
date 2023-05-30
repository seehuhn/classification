package data

// Set implements an abstract interface to represent a test data set,
// consisting of training data for setting up the method and test data
// for assessment.
type Set interface {
	GetName() string
	TrainingData() (data *Data, err error)
	TestData() (data *Data, err error)
}

func MakeSet(name string, train, test *Data) Set {
	return &set{name, train, test}
}

type set struct {
	name  string
	train *Data
	test  *Data
}

func (s *set) GetName() string {
	return s.name
}

func (s *set) TrainingData() (data *Data, err error) {
	return s.train, nil
}

func (s *set) TestData() (data *Data, err error) {
	return s.test, nil
}
