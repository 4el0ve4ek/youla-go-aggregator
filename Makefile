build_sensor_image:
	docker build ./cmd/sensor -t youla_dev_internship_task_go_sensor:latest

build_collector_image:
	docker build ./cmd/collector -t youla_go_collector:latest

run_sensors: build_sensor_image build_collector_image
	docker-compose up -d
