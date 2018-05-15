rm -f artifacts/main
mkdir -p artifacts

cd main
go build -o ../artifacts/echo
cd ..
