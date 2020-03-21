  
.SILENT: build_env test_all build cleanup
all: build_env test_all build cleanup

build_env:
	echo "Creating an integrated build system"
	docker-compose -f docker-compose.yml up --build -d

cleanup:
	echo "Bringing down the build system"
	docker-compose down -v
	# docker rm -f temp-something-of-the-day

test_all: test_ui test_backend test_integration

test_ui:
	echo "Running frontend unit tests and coverage tool"
	docker exec something-frontend-build npm test -- --coverage --watchAll=false

test_backend:
	echo "Running backend unit tests"
	docker exec something-backend-build go test ./...

test_integration:
	echo "Running backend integration tests"
	docker exec something-backend-build sh -c "./db-test.sh"

build:
	echo "Building the app"
	docker build . -t rjmarques/something-of-the-day
	docker create -ti --name temp-something-of-the-day rjmarques/something-of-the-day bash
	docker cp temp-something-of-the-day:/home/something-of-the-day - > ./something-of-the-day.tar
	echo "Artifacts stored in: something-of-the-day.tar"
	echo "Runnable docker image: rjmarques/something-of-the-day"

# run:
# 	echo "Running the app in a container"
# 	docker run --name something-of-the-day --rm --publish=80:80 rjmarques/something-of-the-day