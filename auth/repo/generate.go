package repo

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i AuthRepo -o ./mocks/ -s "_minimock.go"
