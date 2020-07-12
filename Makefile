  
.SILENT: build_env test_all build cleanup
all: build_env test_all build clean

build_env:
	echo "Creating an integrated build system"
	docker-compose -f docker-compose.yml up --build -d

clean:
	echo "Bringing down the build system"
	docker-compose down -v

test_all: test_ui test_backend test_integration

test_ui:
	echo "Running frontend unit tests and coverage tool"
	docker exec something-frontend-build npm test -- --coverage --watchAll=false

test_backend:
	echo "Running backend unit tests"
	docker exec something-backend-build go test ./...

test_integration:
	echo "Running backend integration tests"
	docker exec something-backend-build sh "./db-test.sh"

build:
	echo "Building the app"
	docker build -t rjmarques/something-of-the-day .
	echo "Runnable docker image: rjmarques/something-of-the-day"

run_with_localdb:
	echo "Running the app in a container"
	docker run -e POSTGRES_URL -e CLIENT_SECRET -e CLIENT_ID --name something-of-the-day --rm --publish=80:80 --network=something-of-the-day_integration-tests rjmarques/something-of-the-day

run:
	echo "Running the app in a container"
	docker run -e POSTGRES_URL -e CLIENT_SECRET -e CLIENT_ID --name something-of-the-day --rm --publish=80:80 rjmarques/something-of-the-day	

deploy:
	echo "Deploying the app to ECS"
	aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${ECR}
	docker tag rjmarques/something-of-the-day:latest ${ECR_SOTD_REPO}:latest
	docker push ${ECR_SOTD_REPO}:latest
	aws ecs update-service --cluster sotd-cluster --service sotd-ecs-service --region ${AWS_REGION} --force-new-deployment | cat