LOCAL_PATH := github.com/antongolenev23/voltake-services

.PHONY: fmt
fmt:
	@echo "Formatting all Go files..."
	@goimports -local $(LOCAL_PATH) -w . > /dev/null 2>&1

.PHONY: fmt-check
fmt-check:
	@echo "Checking all Go files format..."
	@if goimports -local $(LOCAL_PATH) -l . 2>&1 | grep . > /dev/null; then \
		echo "Not all Go files are formatted. Run 'make fmt'"; \
		exit 1; \
	else \
		echo "All Go files formatted correctly"; \
	fi

.PHONY: pre-commit
pre-commit: fmt-check
	@echo "Pre-commit checks passed"