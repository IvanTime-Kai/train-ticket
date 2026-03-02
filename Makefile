.PHONY: user-service-run user-service-build user-service-sqlgen

# User service
user-service-run:
	$(MAKE) -C user-service run

user-service-build:
	$(MAKE) -C user-service build

user-service-sqlgen:
	$(MAKE) -C user-service sqlgen
