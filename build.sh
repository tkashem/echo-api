rm -f artifacts/main
mkdir -p artifacts

cd main
GOOS=linux go build -o ../artifacts/main
