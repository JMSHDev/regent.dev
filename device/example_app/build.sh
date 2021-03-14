# this script will build the ExampleApp using a docker version of Golang

# create the build container image
docker build -t build-img .

# create the build container
docker create --name build-cont build-img

# copy the built binaries out of the build container
docker cp build-cont:/build_dir/ExampleApp ../bin/

# remove the build container
docker rm build-cont