.PHONY: build

GIT_BRANCH = $$(git symbolic-ref --short HEAD)
DOCKER_UP = "docker-compose up"
DOCKER_DOWN = "docker-compose down"
DOCKER_RUN = "docker-compose run"

.EXPORT_ALL_VARIABLES:

.DEFAULT: help
help:
	@echo "\n \
	------------------------------ \n \
	++ Python Related ++ \n \
	yq-lint: Runs a linter against a YAML file. Pass in the file with the variable YAML \n \
	  Ex. make yq-lint YAML="extract/postgres_pipeline/manifests/gitlab_com_db_manifest.yaml" \n \
	lint: Runs a linter (Black) over the whole repo. \n \
	mypy: Runs a type-checker in the extract dir. \n \
	pylint: Runs the pylint checker over the whole repo. \n \
	pytest: Unit tests with the pytest python module. \n \
	radon: Runs a cyclomatic complexity checker and shows anything with less than an A rating. \n \
	xenon: Runs a cyclomatic complexity checker that will throw a non-zero exit code if the criteria aren't met. \n \
	\n \
	++ Terraform Related ++ \n \
	format-tf: Format terraform files with 'terraform fmt'. \n \
	\n \
	++ Utilities ++ \n \
	cleanup: WARNING: DELETES DB VOLUME, frees up space and gets rid of old containers/images. \n \
	------------------------------ \n"
build-images:
	@echo "Building images with docker-compose..."
	@docker-compose build

submissions-scraper:
	@echo "Attaching to reddit-scraper..."
	@"$(DOCKER_RUN)" -e MODE=submissions scraper

comments-scraper:
	@echo "Attaching to reddit-scraper..."
	@"$(DOCKER_RUN)" -e MODE=comments scraper

cleanup:
	@echo "Cleaning things up..."
	@"$(DOCKER_DOWN)" -v
	@docker system prune -f

lint:
	@echo "Linting the repo..."
	@black .

yq-lint:
ifdef YAML
	@echo "Linting the YAML file... "
	@echo "Running: yq eval 'sortKeys(..)' $(YAML)"
	@"$(DOCKER_RUN)" data_image bash -c "yq eval 'sortKeys(..)' $(YAML)"
else
	@echo "No file. Exiting."
endif

mypy:
	@echo "Running mypy..."
	@mypy src/ --ignore-missing-imports

pylint:
	@echo "Running pylint..."
	@cd src && pylint *

pytest:
	@echo "Running pytest..."
	@cd src && python -m pytest . --disable-pytest-warnings

radon:
	@echo "Run Radon to compute complexity..."
	@radon cc src --total-average -nb

xenon:
	@echo "Running Xenon..."
	@xenon --max-absolute B --max-modules A --max-average A src -i transform,shared_modules

format-tf:
	@echo "Formatting terraform files..."
	@terraform fmt -recursive tf
